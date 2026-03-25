package policy

import (
	"encoding/json"
	"time"

	"smartsure/pkg/database"

	"github.com/google/uuid"
)

// PolicyInternationalCoverage represents international and travel coverage for a policy
type PolicyInternationalCoverage struct {
	database.BaseModel
	PolicyID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"policy_id"`

	// Coverage Settings
	InternationalCoverage bool   `gorm:"default:false" json:"international_coverage"`
	CoveredCountries      string `gorm:"type:json" json:"covered_countries"` // JSON array of country codes
	ExcludedCountries     string `gorm:"type:json" json:"excluded_countries"`
	WorldwideCoverage     bool   `gorm:"default:false" json:"worldwide_coverage"`

	// Travel Limits
	TravelDaysLimit      int `gorm:"default:90" json:"travel_days_limit"` // Per year
	ConsecutiveDaysLimit int `gorm:"default:30" json:"consecutive_days_limit"`
	TravelDaysUsed       int `json:"travel_days_used"`
	CurrentTripDays      int `json:"current_trip_days"`

	// Additional Coverage
	RoamingChargesCover bool    `gorm:"default:false" json:"roaming_charges_cover"`
	RoamingLimit        float64 `json:"roaming_limit"`
	EmergencyAssistance bool    `gorm:"default:false" json:"emergency_assistance"`
	MedicalEvacuation   bool    `gorm:"default:false" json:"medical_evacuation"`

	// Emergency Services
	EmergencyHotline       string `json:"emergency_hotline"`
	LocalAssistanceNetwork string `gorm:"type:json" json:"local_assistance_network"`
	TranslationService     bool   `gorm:"default:false" json:"translation_service"`

	// Travel History
	LastTravelStartDate *time.Time `json:"last_travel_start_date"`
	LastTravelEndDate   *time.Time `json:"last_travel_end_date"`
	TotalTrips          int        `json:"total_trips"`

	// Status
	IsActive        bool   `gorm:"default:true" json:"is_active"`
	CurrentlyAbroad bool   `gorm:"default:false" json:"currently_abroad"`
	CurrentCountry  string `json:"current_country"`

	// Relationships
	// Note: Policy relationship is handled through embedding in the main Policy struct
}

// TableName returns the table name
func (PolicyInternationalCoverage) TableName() string {
	return "policy_international_coverages"
}

// IsValidInCountry checks if policy is valid in a specific country
func (pic *PolicyInternationalCoverage) IsValidInCountry(countryCode string) bool {
	if !pic.IsActive || !pic.InternationalCoverage {
		return false
	}

	if pic.WorldwideCoverage {
		// Check excluded countries
		var excluded []string
		if err := json.Unmarshal([]byte(pic.ExcludedCountries), &excluded); err == nil {
			for _, exc := range excluded {
				if exc == countryCode {
					return false
				}
			}
		}
		return true
	}

	// Check covered countries
	var covered []string
	if err := json.Unmarshal([]byte(pic.CoveredCountries), &covered); err == nil {
		for _, cov := range covered {
			if cov == countryCode {
				return true
			}
		}
	}

	return false
}

// GetRemainingTravelDays returns remaining travel days
func (pic *PolicyInternationalCoverage) GetRemainingTravelDays() int {
	if !pic.InternationalCoverage || pic.TravelDaysLimit <= 0 {
		return 0
	}
	return pic.TravelDaysLimit - pic.TravelDaysUsed
}

// CanStartNewTrip checks if a new trip can be started
func (pic *PolicyInternationalCoverage) CanStartNewTrip(tripDays int) bool {
	if !pic.IsActive || !pic.InternationalCoverage {
		return false
	}

	if pic.CurrentlyAbroad {
		return false // Already on a trip
	}

	remainingDays := pic.GetRemainingTravelDays()
	if tripDays > remainingDays {
		return false
	}

	if tripDays > pic.ConsecutiveDaysLimit {
		return false
	}

	return true
}

// StartTrip records the start of a new trip
func (pic *PolicyInternationalCoverage) StartTrip(countryCode string) {
	now := time.Now()
	pic.LastTravelStartDate = &now
	pic.CurrentlyAbroad = true
	pic.CurrentCountry = countryCode
	pic.CurrentTripDays = 0
	pic.TotalTrips++
}

// EndTrip records the end of a trip
func (pic *PolicyInternationalCoverage) EndTrip() {
	if !pic.CurrentlyAbroad {
		return
	}

	now := time.Now()
	pic.LastTravelEndDate = &now
	pic.CurrentlyAbroad = false
	pic.TravelDaysUsed += pic.CurrentTripDays
	pic.CurrentTripDays = 0
	pic.CurrentCountry = ""
}

// UpdateTripDays updates the current trip day count
func (pic *PolicyInternationalCoverage) UpdateTripDays() {
	if pic.CurrentlyAbroad && pic.LastTravelStartDate != nil {
		days := int(time.Since(*pic.LastTravelStartDate).Hours() / 24)
		pic.CurrentTripDays = days
	}
}

// HasEmergencyServices checks if emergency services are available
func (pic *PolicyInternationalCoverage) HasEmergencyServices() bool {
	return pic.IsActive &&
		pic.InternationalCoverage &&
		(pic.EmergencyAssistance || pic.MedicalEvacuation)
}

// IsExceedingConsecutiveLimit checks if consecutive days limit is exceeded
func (pic *PolicyInternationalCoverage) IsExceedingConsecutiveLimit() bool {
	pic.UpdateTripDays()
	return pic.CurrentTripDays > pic.ConsecutiveDaysLimit
}

// GetCoverageRegion returns the coverage region type
func (pic *PolicyInternationalCoverage) GetCoverageRegion() string {
	if pic.WorldwideCoverage {
		return "worldwide"
	}

	var covered []string
	if err := json.Unmarshal([]byte(pic.CoveredCountries), &covered); err == nil {
		if len(covered) > 10 {
			return "multi-region"
		} else if len(covered) > 1 {
			return "selected-countries"
		} else if len(covered) == 1 {
			return "single-country"
		}
	}

	return "domestic-only"
}
