package device

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// DeviceCrossRelationship represents relationships between devices
type DeviceCrossRelationship struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Device Family/Ecosystem
	DeviceFamilyID string         `gorm:"type:varchar(100)" json:"device_family_id"`
	EcosystemType  string         `gorm:"type:varchar(50)" json:"ecosystem_type"` // apple, google, samsung
	PrimaryDevice  bool           `gorm:"type:boolean;default:false" json:"primary_device"`
	LinkedDevices  datatypes.JSON `gorm:"type:json" json:"linked_devices"` // []LinkedDevice

	// Companion Devices
	CompanionDevices datatypes.JSON `gorm:"type:json" json:"companion_devices"`      // []CompanionDevice
	PairedDevices    datatypes.JSON `gorm:"type:json" json:"paired_devices"`         // []PairedDevice
	ConnectionType   string         `gorm:"type:varchar(50)" json:"connection_type"` // bluetooth, wifi, cellular
	LastSyncTime     *time.Time     `gorm:"type:timestamp" json:"last_sync_time,omitempty"`
	SyncStatus       string         `gorm:"type:varchar(50)" json:"sync_status"`

	// Device Bundles/Kits
	BundleID         string         `gorm:"type:varchar(100)" json:"bundle_id"`
	BundleName       string         `gorm:"type:varchar(255)" json:"bundle_name"`
	BundleComponents datatypes.JSON `gorm:"type:json" json:"bundle_components"` // []BundleComponent
	BundleDiscount   float64        `gorm:"type:decimal(5,2)" json:"bundle_discount"`
	BundleValue      float64        `gorm:"type:decimal(15,2)" json:"bundle_value"`

	// Shared Accessories
	SharedAccessories      datatypes.JSON `gorm:"type:json" json:"shared_accessories"`      // []Accessory
	AccessoryCompatibility datatypes.JSON `gorm:"type:json" json:"accessory_compatibility"` // map[string]bool

	// Cross-Device Warranties
	SharedWarranty    bool    `gorm:"type:boolean;default:false" json:"shared_warranty"`
	MasterWarrantyID  string  `gorm:"type:varchar(100)" json:"master_warranty_id"`
	CombinedCoverage  float64 `gorm:"type:decimal(15,2)" json:"combined_coverage"`
	GroupDiscountRate float64 `gorm:"type:decimal(5,2)" json:"group_discount_rate"`

	// Family Sharing
	FamilySharingEnabled bool           `gorm:"type:boolean;default:false" json:"family_sharing_enabled"`
	FamilyGroupID        string         `gorm:"type:varchar(100)" json:"family_group_id"`
	FamilyMembers        datatypes.JSON `gorm:"type:json" json:"family_members"`  // []FamilyMember
	SharedBenefits       datatypes.JSON `gorm:"type:json" json:"shared_benefits"` // []Benefit

	// Device Hierarchy
	ParentDeviceID *uuid.UUID     `gorm:"type:uuid" json:"parent_device_id,omitempty"`
	ChildDevices   datatypes.JSON `gorm:"type:json" json:"child_devices"` // []uuid.UUID
	HierarchyLevel int            `gorm:"type:int" json:"hierarchy_level"`
	DeviceRole     string         `gorm:"type:varchar(50)" json:"device_role"` // master, slave, peer

	// Inter-Device Communication
	CommunicationProtocol string         `gorm:"type:varchar(100)" json:"communication_protocol"`
	DataSharingEnabled    bool           `gorm:"type:boolean;default:false" json:"data_sharing_enabled"`
	SharedDataTypes       datatypes.JSON `gorm:"type:json" json:"shared_data_types"` // []string
	LastDataExchange      *time.Time     `gorm:"type:timestamp" json:"last_data_exchange,omitempty"`

	// Status
	RelationshipStatus string    `gorm:"type:varchar(50)" json:"relationship_status"`
	CreatedAt          time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// DeviceGamification represents gamification and engagement features
type DeviceGamification struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Care Achievements
	CareAchievements     datatypes.JSON `gorm:"type:json" json:"care_achievements"` // []Achievement
	TotalAchievements    int            `gorm:"type:int" json:"total_achievements"`
	UnlockedAchievements int            `gorm:"type:int" json:"unlocked_achievements"`
	AchievementPoints    int            `gorm:"type:int" json:"achievement_points"`

	// Safety Scores & Badges
	SafetyScore  float64        `gorm:"type:decimal(5,2)" json:"safety_score"`
	SafetyBadges datatypes.JSON `gorm:"type:json" json:"safety_badges"` // []Badge
	SafetyLevel  int            `gorm:"type:int" json:"safety_level"`
	SafetyRank   string         `gorm:"type:varchar(50)" json:"safety_rank"`

	// Maintenance Streaks
	MaintenanceStreak int            `gorm:"type:int" json:"maintenance_streak_days"`
	LongestStreak     int            `gorm:"type:int" json:"longest_streak_days"`
	StreakStartDate   *time.Time     `gorm:"type:timestamp" json:"streak_start_date,omitempty"`
	StreakRewards     datatypes.JSON `gorm:"type:json" json:"streak_rewards"` // []Reward

	// Referral Rewards
	ReferralCode        string         `gorm:"type:varchar(100);unique" json:"referral_code"`
	ReferralCount       int            `gorm:"type:int" json:"referral_count"`
	SuccessfulReferrals int            `gorm:"type:int" json:"successful_referrals"`
	ReferralRewards     datatypes.JSON `gorm:"type:json" json:"referral_rewards"` // []Reward
	ReferralPoints      int            `gorm:"type:int" json:"referral_points"`

	// Community Challenges
	ActiveChallenges    datatypes.JSON `gorm:"type:json" json:"active_challenges"`    // []Challenge
	CompletedChallenges datatypes.JSON `gorm:"type:json" json:"completed_challenges"` // []Challenge
	ChallengePoints     int            `gorm:"type:int" json:"challenge_points"`
	ChallengeRank       int            `gorm:"type:int" json:"challenge_rank"`

	// Experience & Levels
	ExperiencePoints    int     `gorm:"type:int" json:"experience_points"`
	CurrentLevel        int     `gorm:"type:int" json:"current_level"`
	NextLevelPoints     int     `gorm:"type:int" json:"next_level_points"`
	ProgressToNextLevel float64 `gorm:"type:decimal(5,2)" json:"progress_to_next_level"`

	// Leaderboards
	GlobalRank          int    `gorm:"type:int" json:"global_rank"`
	RegionalRank        int    `gorm:"type:int" json:"regional_rank"`
	FriendsRank         int    `gorm:"type:int" json:"friends_rank"`
	LeaderboardCategory string `gorm:"type:varchar(50)" json:"leaderboard_category"`

	// Rewards & Incentives
	TotalRewards     int            `gorm:"type:int" json:"total_rewards"`
	RedeemedRewards  int            `gorm:"type:int" json:"redeemed_rewards"`
	AvailableRewards datatypes.JSON `gorm:"type:json" json:"available_rewards"` // []Reward
	RewardPoints     int            `gorm:"type:int" json:"reward_points"`

	// Milestones
	Milestones        datatypes.JSON `gorm:"type:json" json:"milestones"` // []Milestone
	NextMilestone     datatypes.JSON `gorm:"type:json" json:"next_milestone"`
	MilestoneProgress float64        `gorm:"type:decimal(5,2)" json:"milestone_progress"`

	// Status
	GamificationStatus string    `gorm:"type:varchar(50)" json:"gamification_status"`
	LastActivityDate   time.Time `gorm:"type:timestamp" json:"last_activity_date"`
	CreatedAt          time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// DeviceEducation represents educational and training features
type DeviceEducation struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	DeviceID uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`

	// Care Tutorials
	ViewedTutorials      datatypes.JSON `gorm:"type:json" json:"viewed_tutorials"` // []Tutorial
	CompletedTutorials   int            `gorm:"type:int" json:"completed_tutorials"`
	TutorialProgress     datatypes.JSON `gorm:"type:json" json:"tutorial_progress"`     // map[string]float64
	RecommendedTutorials datatypes.JSON `gorm:"type:json" json:"recommended_tutorials"` // []Tutorial

	// Insurance Education
	InsuranceLessons datatypes.JSON `gorm:"type:json" json:"insurance_lessons"` // []Lesson
	CompletedLessons int            `gorm:"type:int" json:"completed_lessons"`
	EducationLevel   string         `gorm:"type:varchar(50)" json:"education_level"` // beginner, intermediate, expert
	QuizScores       datatypes.JSON `gorm:"type:json" json:"quiz_scores"`            // map[string]float64

	// Safety Tips & Guides
	SafetyTips     datatypes.JSON `gorm:"type:json" json:"safety_tips"` // []SafetyTip
	ViewedTips     int            `gorm:"type:int" json:"viewed_tips"`
	AppliedTips    int            `gorm:"type:int" json:"applied_tips"`
	CustomizedTips datatypes.JSON `gorm:"type:json" json:"customized_tips"` // []Tip

	// Feature Discovery
	DiscoveredFeatures   datatypes.JSON `gorm:"type:json" json:"discovered_features"`    // []Feature
	FeatureUsageTracking datatypes.JSON `gorm:"type:json" json:"feature_usage_tracking"` // map[string]Usage
	UnusedFeatures       datatypes.JSON `gorm:"type:json" json:"unused_features"`        // []Feature
	FeatureTutorials     datatypes.JSON `gorm:"type:json" json:"feature_tutorials"`      // map[string]Tutorial

	// User Proficiency
	ProficiencyScore float64        `gorm:"type:decimal(5,2)" json:"proficiency_score"`
	ProficiencyLevel string         `gorm:"type:varchar(50)" json:"proficiency_level"`
	SkillsAssessment datatypes.JSON `gorm:"type:json" json:"skills_assessment"` // Assessment
	ImprovementAreas datatypes.JSON `gorm:"type:json" json:"improvement_areas"` // []Area

	// Learning Paths
	EnrolledPaths       datatypes.JSON `gorm:"type:json" json:"enrolled_paths"`  // []LearningPath
	CompletedPaths      datatypes.JSON `gorm:"type:json" json:"completed_paths"` // []LearningPath
	CurrentPathProgress float64        `gorm:"type:decimal(5,2)" json:"current_path_progress"`
	NextLesson          datatypes.JSON `gorm:"type:json" json:"next_lesson"`

	// Certifications
	EarnedCertifications datatypes.JSON `gorm:"type:json" json:"earned_certifications"` // []Certification
	CertificationPoints  int            `gorm:"type:int" json:"certification_points"`
	NextCertification    datatypes.JSON `gorm:"type:json" json:"next_certification"`

	// Help & Support Access
	HelpArticlesViewed    int            `gorm:"type:int" json:"help_articles_viewed"`
	SupportTickets        int            `gorm:"type:int" json:"support_tickets"`
	FAQAccess             datatypes.JSON `gorm:"type:json" json:"faq_access"` // []FAQAccess
	VideoTutorialsWatched int            `gorm:"type:int" json:"video_tutorials_watched"`

	// Knowledge Base
	SavedArticles     datatypes.JSON `gorm:"type:json" json:"saved_articles"`     // []Article
	BookmarkedContent datatypes.JSON `gorm:"type:json" json:"bookmarked_content"` // []Content
	PersonalNotes     datatypes.JSON `gorm:"type:json" json:"personal_notes"`     // []Note

	// Status
	EducationStatus  string    `gorm:"type:varchar(50)" json:"education_status"`
	LastLearningDate time.Time `gorm:"type:timestamp" json:"last_learning_date"`
	CreatedAt        time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// =====================================
// METHODS
// =====================================

// IsInEcosystem checks if device is part of an ecosystem
func (dcr *DeviceCrossRelationship) IsInEcosystem() bool {
	return dcr.DeviceFamilyID != "" || len(dcr.LinkedDevices) > 0
}

// HasCompanionDevices checks if device has companions
func (dcr *DeviceCrossRelationship) HasCompanionDevices() bool {
	return dcr.CompanionDevices != nil || dcr.PairedDevices != nil
}

// IsBundled checks if device is part of a bundle
func (dcr *DeviceCrossRelationship) IsBundled() bool {
	return dcr.BundleID != "" && dcr.BundleComponents != nil
}

// GetAchievementLevel calculates achievement level
func (dg *DeviceGamification) GetAchievementLevel() string {
	ratio := float64(dg.UnlockedAchievements) / float64(dg.TotalAchievements)
	if ratio >= 0.9 {
		return "master"
	} else if ratio >= 0.7 {
		return "expert"
	} else if ratio >= 0.5 {
		return "advanced"
	} else if ratio >= 0.3 {
		return "intermediate"
	}
	return "beginner"
}

// HasActiveStreak checks if user has active maintenance streak
func (dg *DeviceGamification) HasActiveStreak() bool {
	return dg.MaintenanceStreak > 0 && dg.StreakStartDate != nil
}

// IsTopPerformer checks if user is in top rankings
func (dg *DeviceGamification) IsTopPerformer() bool {
	return dg.GlobalRank <= 100 || dg.RegionalRank <= 10
}

// GetEducationCompletionRate calculates education completion
func (de *DeviceEducation) GetEducationCompletionRate() float64 {
	if de.CompletedTutorials == 0 {
		return 0
	}
	totalContent := de.CompletedTutorials + de.CompletedLessons
	if totalContent == 0 {
		return 0
	}
	// Use CompletedTutorials as the numerator since ViewedTutorials is JSON
	return float64(de.CompletedTutorials) / float64(totalContent) * 100
}

// IsProficient checks if user is proficient
func (de *DeviceEducation) IsProficient() bool {
	return de.ProficiencyScore >= 80 || de.ProficiencyLevel == "expert"
}

// NeedsTraining checks if user needs additional training
func (de *DeviceEducation) NeedsTraining() bool {
	return de.ProficiencyScore < 50 || de.ImprovementAreas != nil
}
