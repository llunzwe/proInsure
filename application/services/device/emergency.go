package device

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// InitiateEmergencyLock remotely locks the device
func (d *Device) InitiateEmergencyLock() error {
	// Check if remote lock is enabled
	if !d.RemoteLockEnabled {
		return fmt.Errorf("remote lock not enabled for this device")
	}

	// Check if device is already locked
	if d.IsLocked {
		return fmt.Errorf("device is already locked")
	}

	// Set device as locked
	d.IsLocked = true
	d.Status = "locked"

	// Log the lock event
	lockEvent := map[string]interface{}{
		"device_id":     d.ID,
		"action":        "emergency_lock",
		"timestamp":     time.Now(),
		"imei":          d.IMEI,
		"last_location": d.LastKnownLocation,
		"reason":        "Emergency lock initiated",
	}

	// In production, this would trigger actual device lock via MDM/API
	fmt.Printf("Emergency Lock Initiated: %+v\n", lockEvent)

	// Update security status
	if !d.IsStolen {
		d.TriggerSecurityAlert("Device locked for security - potential theft or loss")
	}

	return nil
}

// TriggerDataWipe initiates remote data wipe
func (d *Device) TriggerDataWipe() error {
	// Critical security check - device must be reported stolen or compromised
	if !d.IsStolen && d.Status != "compromised" {
		return fmt.Errorf("data wipe requires device to be reported stolen or compromised")
	}

	// Check if device is verified (to prevent accidental wipes)
	if !d.IsVerified {
		return fmt.Errorf("device must be verified before remote wipe")
	}

	// Generate wipe confirmation code
	wipeCode := d.GenerateRecoveryCode()

	// Log the wipe request
	wipeEvent := map[string]interface{}{
		"device_id":         d.ID,
		"action":            "remote_wipe",
		"timestamp":         time.Now(),
		"imei":              d.IMEI,
		"confirmation_code": wipeCode,
		"last_location":     d.LastKnownLocation,
		"owner_id":          d.OwnerID,
	}

	// In production, this would trigger actual data wipe via MDM/API
	fmt.Printf("CRITICAL: Data Wipe Initiated: %+v\n", wipeEvent)

	// Update device status
	d.Status = "wiped"
	now := time.Now()
	d.LastLocationUpdate = &now

	// Mark device as end of life after wipe
	d.Category = "wiped_device"

	return nil
}

// GetLastBackupDate returns last backup timestamp
func (d *Device) GetLastBackupDate() *time.Time {
	// In production, this would query backup service
	// For now, simulate based on device activity

	// If device has recent inspection, assume backup was done
	if d.LastInspection != nil && time.Since(*d.LastInspection) < 30*24*time.Hour {
		return d.LastInspection
	}

	// If device was recently verified, assume backup
	if d.LastVerifiedAt != nil && time.Since(*d.LastVerifiedAt) < 60*24*time.Hour {
		return d.LastVerifiedAt
	}

	// Check for recent claim activity (usually involves backup)
	for _, claim := range d.Claims {
		if claim.Status == "processing" || claim.Status == "approved" {
			if time.Since(claim.CreatedAt) < 90*24*time.Hour {
				return &claim.CreatedAt
			}
		}
	}

	// No recent backup found
	return nil
}

// IsRecoverable checks if device can be recovered
func (d *Device) IsRecoverable() bool {
	// Cannot recover if wiped
	if d.Status == "wiped" {
		return false
	}

	// Cannot recover if blacklisted
	if d.BlacklistStatus == "blocked" {
		return false
	}

	// Check if device is stolen but location is known
	if d.IsStolen {
		// Can potentially recover if we have recent location
		if d.LastLocationUpdate != nil && time.Since(*d.LastLocationUpdate) < 24*time.Hour {
			return true
		}
		// Can recover if Find My Device is enabled
		if d.FindMyDeviceEnabled {
			return true
		}
		return false
	}

	// Check physical condition for recovery
	if d.Condition == "poor" && d.WaterDamageIndicator != "white" {
		return false // Severely damaged
	}

	// Check if device is locked but can be unlocked
	if d.IsLocked && d.RemoteLockEnabled {
		return true // Can be remotely unlocked
	}

	// Device is recoverable in most other cases
	return true
}

// GenerateRecoveryCode creates recovery access code
func (d *Device) GenerateRecoveryCode() string {
	// Generate secure recovery code
	bytes := make([]byte, 8)
	rand.Read(bytes)
	code := hex.EncodeToString(bytes)

	// Format as readable code
	formattedCode := fmt.Sprintf("%s-%s-%s-%s",
		code[0:4], code[4:8], code[8:12], code[12:16])

	return formattedCode
}

// InitiateRecoveryMode puts device in recovery mode
func (d *Device) InitiateRecoveryMode() error {
	// Check if device is eligible for recovery
	if !d.IsRecoverable() {
		return fmt.Errorf("device is not recoverable in current state")
	}

	// Generate recovery code
	recoveryCode := d.GenerateRecoveryCode()

	// Create recovery session
	recoverySession := map[string]interface{}{
		"device_id":        d.ID,
		"recovery_code":    recoveryCode,
		"initiated_at":     time.Now(),
		"expires_at":       time.Now().Add(24 * time.Hour),
		"device_status":    d.Status,
		"is_stolen":        d.IsStolen,
		"is_locked":        d.IsLocked,
		"last_location":    d.LastKnownLocation,
		"backup_available": d.GetLastBackupDate() != nil,
	}

	// Set device to recovery mode
	previousStatus := d.Status
	d.Status = "recovery_mode"

	// Log recovery initiation
	fmt.Printf("Recovery Mode Initiated: %+v\n", recoverySession)

	// If stolen, attempt to track
	if d.IsStolen && d.FindMyDeviceEnabled {
		// In production, would trigger location tracking
		fmt.Println("Activating enhanced location tracking for stolen device recovery")
	}

	// If locked, prepare unlock sequence
	if d.IsLocked && d.RemoteLockEnabled {
		fmt.Printf("Device can be unlocked with recovery code: %s\n", recoveryCode)
	}

	// Store previous status for restoration
	recoverySession["previous_status"] = previousStatus

	return nil
}

// ExecuteEmergencyProtocol runs full emergency protection
func (d *Device) ExecuteEmergencyProtocol(reason string) error {
	protocol := []string{}

	// Step 1: Mark as stolen if applicable
	if strings.Contains(strings.ToLower(reason), "stolen") || strings.Contains(strings.ToLower(reason), "theft") {
		d.MarkAsStolen(reason)
		protocol = append(protocol, "Device marked as stolen")
	}

	// Step 2: Enable all security features
	if !d.FindMyDeviceEnabled {
		d.FindMyDeviceEnabled = true
		protocol = append(protocol, "Find My Device activated")
	}

	if !d.RemoteLockEnabled {
		d.RemoteLockEnabled = true
		protocol = append(protocol, "Remote lock enabled")
	}

	if !d.EncryptionEnabled {
		d.EncryptionEnabled = true
		protocol = append(protocol, "Encryption verified")
	}

	// Step 3: Lock device
	if !d.IsLocked {
		if err := d.InitiateEmergencyLock(); err == nil {
			protocol = append(protocol, "Device locked remotely")
		}
	}

	// Step 4: Update location if possible
	if d.FindMyDeviceEnabled {
		now := time.Now()
		d.LastLocationUpdate = &now
		protocol = append(protocol, "Location tracking activated")
	}

	// Step 5: Trigger security alert
	d.TriggerSecurityAlert(reason)
	protocol = append(protocol, "Security alert triggered")

	// Step 6: Check backup status
	lastBackup := d.GetLastBackupDate()
	if lastBackup != nil {
		protocol = append(protocol, fmt.Sprintf("Last backup: %s", lastBackup.Format("2006-01-02")))
	} else {
		protocol = append(protocol, "WARNING: No recent backup found")
	}

	// Step 7: Update blacklist if fraud suspected
	if strings.Contains(strings.ToLower(reason), "fraud") {
		d.BlacklistStatus = "checking"
		protocol = append(protocol, "Blacklist check initiated")
	}

	// Log emergency protocol execution
	fmt.Printf("Emergency Protocol Executed for Device %s:\n", d.ID)
	for i, step := range protocol {
		fmt.Printf("%d. %s\n", i+1, step)
	}

	return nil
}

// CreateEmergencyContactInfo generates emergency contact information
func (d *Device) CreateEmergencyContactInfo() map[string]interface{} {
	info := map[string]interface{}{
		"device_id":     d.ID,
		"imei":          d.IMEI,
		"serial_number": d.SerialNumber,
		"model":         fmt.Sprintf("%s %s", d.Brand, d.Model),
		"status":        d.Status,
		"is_stolen":     d.IsStolen,
		"is_locked":     d.IsLocked,
	}

	// Add location if available
	if d.LastKnownLocation != "" {
		info["last_known_location"] = d.LastKnownLocation
		if d.LastLocationUpdate != nil {
			info["location_timestamp"] = d.LastLocationUpdate
			info["hours_since_location"] = time.Since(*d.LastLocationUpdate).Hours()
		}
	}

	// Add security status
	security := map[string]bool{
		"find_my_device": d.FindMyDeviceEnabled,
		"remote_lock":    d.RemoteLockEnabled,
		"encryption":     d.EncryptionEnabled,
		"anti_theft":     d.AntiTheftAppInstalled,
		"can_be_wiped":   d.IsStolen && d.IsVerified,
		"can_be_locked":  d.RemoteLockEnabled && !d.IsLocked,
		"can_be_tracked": d.FindMyDeviceEnabled,
		"is_recoverable": d.IsRecoverable(),
	}
	info["security_capabilities"] = security

	// Add emergency contacts
	contacts := map[string]string{
		"emergency_hotline": "1-800-DEVICE-911",
		"theft_report_line": "1-800-STOLEN",
		"technical_support": "1-800-TECH-HELP",
		"insurance_claims":  "1-800-CLAIMS",
		"police_reference":  d.generatePoliceReference(),
	}
	info["emergency_contacts"] = contacts

	// Add recovery options
	recovery := []string{}
	if d.IsRecoverable() {
		if d.FindMyDeviceEnabled {
			recovery = append(recovery, "Track device location")
		}
		if d.RemoteLockEnabled && !d.IsLocked {
			recovery = append(recovery, "Lock device remotely")
		}
		if d.IsLocked && d.RemoteLockEnabled {
			recovery = append(recovery, "Unlock with recovery code")
		}
		if d.GetLastBackupDate() != nil {
			recovery = append(recovery, "Restore from backup")
		}
	}
	info["recovery_options"] = recovery

	// Add case number
	info["case_number"] = d.generateCaseNumber()
	info["generated_at"] = time.Now()

	return info
}

// generatePoliceReference creates a police reference number
func (d *Device) generatePoliceReference() string {
	// Generate based on device ID and timestamp
	timestamp := time.Now().Format("20060102")
	return fmt.Sprintf("POL-%s-%s", timestamp, d.ID.String()[:8])
}

// generateCaseNumber creates an emergency case number
func (d *Device) generateCaseNumber() string {
	// Generate unique case number
	timestamp := time.Now().Unix()
	return fmt.Sprintf("EMG-%d-%s", timestamp, d.ID.String()[:6])
}
