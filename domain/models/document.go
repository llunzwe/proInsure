package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"smartsure/internal/domain/models/document"
)

// Document represents a stored document with comprehensive metadata and lifecycle management
type Document struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	// Embedded structs for organization
	document.DocumentMetadata
	document.DocumentStorage
	document.DocumentSecurity
	document.DocumentLifecycle
	document.DocumentRelationships
	document.DocumentCompliance
	document.DocumentProcessing
	document.DocumentWorkflow
	document.DocumentSharing
	document.DocumentAnalytics

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User       *User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Policy     *Policy           `gorm:"foreignKey:PolicyID" json:"policy,omitempty"`
	Claim      *Claim            `gorm:"foreignKey:ClaimID" json:"claim,omitempty"`
	Payment    *Payment          `gorm:"foreignKey:PaymentID" json:"payment,omitempty"`
	Device     *Device           `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	Repair     *RepairBooking   `gorm:"foreignKey:RepairID" json:"repair,omitempty"`
	Vendor     *Vendor           `gorm:"foreignKey:VendorID" json:"vendor,omitempty"`
	Technician *Technician       `gorm:"foreignKey:TechnicianID" json:"technician,omitempty"`
	Corporate  *CorporateAccount `gorm:"foreignKey:CorporateID" json:"corporate,omitempty"`
	Signatures []ESignature      `gorm:"foreignKey:DocumentID" json:"signatures,omitempty"`
	Versions   []Document        `gorm:"foreignKey:ParentDocumentID" json:"versions,omitempty"`
	AccessLogs []DocumentAccess  `gorm:"foreignKey:DocumentID" json:"access_logs,omitempty"`
	Parent     *Document         `gorm:"foreignKey:ParentDocumentID" json:"parent,omitempty"`
}

// ESignature wraps the base signature model with relationships
type ESignature struct {
	document.ESignature

	// Relationships
	Document *Document `gorm:"foreignKey:DocumentID" json:"document,omitempty"`
	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// DocumentTemplate wraps the base template model with relationships
type DocumentTemplate struct {
	document.DocumentTemplate

	// Relationships
	CreatedByUser       *User                `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
	DocumentGenerations []DocumentGeneration `gorm:"foreignKey:TemplateID" json:"generations,omitempty"`
}

// DocumentGeneration wraps the base generation model with relationships
type DocumentGeneration struct {
	document.DocumentGeneration

	// Relationships
	Template *DocumentTemplate `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
	Document *Document         `gorm:"foreignKey:DocumentID" json:"document,omitempty"`
	User     *User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// DocumentAccess wraps the base access model with relationships
type DocumentAccess struct {
	document.DocumentAccess

	// Relationships
	Document      *Document `gorm:"foreignKey:DocumentID" json:"document,omitempty"`
	User          *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	GrantedByUser *User     `gorm:"foreignKey:GrantedBy" json:"granted_by_user,omitempty"`
	RevokedByUser *User     `gorm:"foreignKey:RevokedBy" json:"revoked_by_user,omitempty"`
}

// BeforeCreate hook for UUID generation
func (d *Document) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	if d.Status == "" {
		d.Status = document.StatusDraft
	}
	if d.SecurityLevel == "" {
		d.SecurityLevel = document.SecurityInternal
	}
	if d.VerificationStatus == "" {
		d.VerificationStatus = document.VerificationNotRequired
	}
	return nil
}

// Business Logic Methods

// IsActive checks if document is active
func (d *Document) IsActive() bool {
	return d.Status == document.StatusActive &&
		d.DeletedAt.Time.IsZero() &&
		(d.ExpiresAt == nil || d.ExpiresAt.After(time.Now()))
}

// IsExpired checks if document has expired
func (d *Document) IsExpired() bool {
	return d.Status == document.StatusExpired ||
		(d.ExpiresAt != nil && d.ExpiresAt.Before(time.Now()))
}

// RequiresVerification checks if document requires verification
func (d *Document) RequiresVerification() bool {
	sensitiveTypes := []string{
		document.TypeLegalDocument,
		document.TypeContract,
		document.TypeIDProof,
		document.TypeCertificate,
	}
	for _, t := range sensitiveTypes {
		if d.Type == t {
			return true
		}
	}
	return d.SecurityLevel == document.SecurityRestricted ||
		d.SecurityLevel == document.SecurityConfidential
}

// CanAccess checks if a user can access the document
func (d *Document) CanAccess(userID uuid.UUID, accessType string) bool {
	// Owner always has access
	if d.UserID == userID {
		return true
	}

	// Check if document is public for view access
	if d.IsPublic && accessType == document.AccessView {
		return true
	}

	// Check access logs for permissions
	for _, access := range d.AccessLogs {
		if access.UserID == userID && access.IsActive {
			if access.ExpiresAt == nil || access.ExpiresAt.After(time.Now()) {
				if access.AccessType == accessType || access.Permission == document.PermissionOwner {
					return true
				}
			}
		}
	}

	return false
}

// GetSize returns human-readable file size
func (d *Document) GetSize() string {
	const unit = 1024
	if d.Size < unit {
		return fmt.Sprintf("%d B", d.Size)
	}
	div, exp := int64(unit), 0
	for n := d.Size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(d.Size)/float64(div), "KMGTPE"[exp])
}

// IsValidMimeType checks if the mime type is acceptable
func (d *Document) IsValidMimeType() bool {
	allowedTypes := []string{
		"application/pdf",
		"image/jpeg", "image/jpg", "image/png", "image/gif",
		"application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"text/plain", "text/html",
		"application/json", "application/xml",
	}

	for _, allowed := range allowedTypes {
		if strings.EqualFold(d.MimeType, allowed) {
			return true
		}
	}
	return false
}

// NeedsOCR checks if document needs OCR processing
func (d *Document) NeedsOCR() bool {
	if d.OCRProcessed {
		return false
	}

	imageTypes := []string{"image/jpeg", "image/jpg", "image/png", "image/tiff"}
	for _, imgType := range imageTypes {
		if strings.Contains(d.MimeType, imgType) {
			return true
		}
	}

	return d.MimeType == "application/pdf" && d.OCRText == ""
}

// IsCompliant checks if document meets compliance requirements
func (d *Document) IsCompliant() bool {
	if !d.ComplianceChecked {
		return false
	}

	return d.ComplianceStatus == "compliant" &&
		(!d.RequiresVerification() || d.IsVerified)
}

// HasValidSignatures checks if all required signatures are valid
func (d *Document) HasValidSignatures() bool {
	if len(d.Signatures) == 0 {
		return false
	}

	for _, sig := range d.Signatures {
		if sig.Status != document.SignatureSigned {
			return false
		}
		if sig.ExpiresAt != nil && sig.ExpiresAt.Before(time.Now()) {
			return false
		}
	}

	return true
}

// GetTags returns parsed tags
func (d *Document) GetTags() []string {
	if d.Tags == "" {
		return []string{}
	}

	var tags []string
	_ = json.Unmarshal([]byte(d.Tags), &tags)
	return tags
}

// AddTag adds a new tag
func (d *Document) AddTag(tag string) {
	tags := d.GetTags()

	// Check if tag already exists
	for _, existingTag := range tags {
		if strings.EqualFold(existingTag, tag) {
			return
		}
	}

	tags = append(tags, tag)
	tagsJSON, _ := json.Marshal(tags)
	d.Tags = string(tagsJSON)
}

// GetMetadata returns parsed metadata
func (d *Document) GetMetadata() map[string]interface{} {
	if d.Metadata == "" {
		return make(map[string]interface{})
	}

	var metadata map[string]interface{}
	_ = json.Unmarshal([]byte(d.Metadata), &metadata)
	return metadata
}

// SetMetadata sets metadata field
func (d *Document) SetMetadata(key string, value interface{}) {
	metadata := d.GetMetadata()
	metadata[key] = value
	metadataJSON, _ := json.Marshal(metadata)
	d.Metadata = string(metadataJSON)
}

// IncrementViewCount increments view count
func (d *Document) IncrementViewCount() {
	d.ViewCount++
	now := time.Now()
	d.LastViewedAt = &now
}

// IncrementDownloadCount increments download count
func (d *Document) IncrementDownloadCount() {
	d.DownloadCount++
	now := time.Now()
	d.LastDownloadedAt = &now
}

// CanDelete checks if document can be deleted
func (d *Document) CanDelete() bool {
	// Cannot delete if on legal hold
	if d.LegalHold {
		return false
	}

	// Cannot delete if part of active workflow
	if d.WorkflowStatus == "in_progress" {
		return false
	}

	// Cannot delete if has active signatures
	for _, sig := range d.Signatures {
		if sig.Status == document.SignatureSigned {
			return false
		}
	}

	return true
}

// Archive archives the document
func (d *Document) Archive(userID uuid.UUID) {
	d.Status = document.StatusArchived
	now := time.Now()
	d.ArchivedAt = &now
	d.ArchivedBy = &userID
}

// GetRetentionDays returns retention period in days
func (d *Document) GetRetentionDays() int {
	switch d.Category {
	case document.CategoryLegal:
		return document.RetentionLongTerm
	case document.CategoryFinancial:
		return document.RetentionLongTerm
	case document.CategoryCompliance:
		return document.RetentionPermanent
	case document.CategoryClaim:
		return document.RetentionMediumTerm
	default:
		return document.RetentionShortTerm
	}
}

// ShouldDispose checks if document should be disposed
func (d *Document) ShouldDispose() bool {
	if d.LegalHold || d.Status != document.StatusArchived {
		return false
	}

	if d.DisposalDate != nil && d.DisposalDate.Before(time.Now()) {
		return true
	}

	if d.ArchivedAt != nil {
		retentionDays := d.GetRetentionDays()
		expiryDate := d.ArchivedAt.AddDate(0, 0, retentionDays)
		return expiryDate.Before(time.Now())
	}

	return false
}

// GetSecurityClassification returns security classification
func (d *Document) GetSecurityClassification() string {
	if d.SecurityLevel != "" {
		return d.SecurityLevel
	}

	// Determine based on type and category
	if d.Type == document.TypeIDProof || d.Category == document.CategoryIdentification {
		return document.SecurityConfidential
	}

	if d.Category == document.CategoryLegal || d.Category == document.CategoryCompliance {
		return document.SecurityRestricted
	}

	if d.IsPublic {
		return document.SecurityPublic
	}

	return document.SecurityInternal
}

// GetRelatedDocumentIDs returns parsed related document IDs
func (d *Document) GetRelatedDocumentIDs() []uuid.UUID {
	if d.RelatedDocuments == "" {
		return []uuid.UUID{}
	}

	var ids []uuid.UUID
	_ = json.Unmarshal([]byte(d.RelatedDocuments), &ids)
	return ids
}

// CalculatePopularityScore calculates document popularity
func (d *Document) CalculatePopularityScore() float64 {
	score := 0.0

	// View weight
	score += float64(d.ViewCount) * 1.0

	// Download weight
	score += float64(d.DownloadCount) * 2.0

	// Share weight
	score += float64(d.ShareCount) * 3.0

	// Citation weight
	score += float64(d.CitationCount) * 5.0

	// Recency factor
	if d.LastViewedAt != nil {
		daysSinceView := time.Since(*d.LastViewedAt).Hours() / 24
		if daysSinceView < 7 {
			score *= 1.5
		} else if daysSinceView < 30 {
			score *= 1.2
		}
	}

	d.PopularityScore = score
	return score
}
