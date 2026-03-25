package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// UserGamification manages gamification elements and achievements
type UserGamification struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Points & Levels
	TotalPoints      int     `gorm:"default:0" json:"total_points"`
	CurrentPoints    int     `gorm:"default:0" json:"current_points"`
	SpentPoints      int     `gorm:"default:0" json:"spent_points"`
	Level            int     `gorm:"default:1" json:"level"`
	LevelProgress    float64 `gorm:"default:0" json:"level_progress"`
	NextLevelPoints  int     `gorm:"default:0" json:"next_level_points"`
	ExperiencePoints int     `gorm:"default:0" json:"experience_points"`

	// Achievements
	UnlockedAchievements []map[string]interface{} `gorm:"type:json" json:"unlocked_achievements"`
	LockedAchievements   []map[string]interface{} `gorm:"type:json" json:"locked_achievements"`
	AchievementProgress  map[string]float64       `gorm:"type:json" json:"achievement_progress"`
	TotalAchievements    int                      `gorm:"default:0" json:"total_achievements"`
	AchievementScore     float64                  `gorm:"default:0" json:"achievement_score"`

	// Badges
	EarnedBadges    []map[string]interface{} `gorm:"type:json" json:"earned_badges"`
	DisplayedBadges []string                 `gorm:"type:json" json:"displayed_badges"`
	RareBadges      []string                 `gorm:"type:json" json:"rare_badges"`
	BadgeCategories map[string]int           `gorm:"type:json" json:"badge_categories"`

	// Streaks
	CurrentStreak   int                      `gorm:"default:0" json:"current_streak"`
	LongestStreak   int                      `gorm:"default:0" json:"longest_streak"`
	StreakStartDate *time.Time               `json:"streak_start_date"`
	StreakTypes     map[string]int           `gorm:"type:json" json:"streak_types"`
	StreakRewards   []map[string]interface{} `gorm:"type:json" json:"streak_rewards"`

	// Challenges
	ActiveChallenges    []map[string]interface{} `gorm:"type:json" json:"active_challenges"`
	CompletedChallenges []map[string]interface{} `gorm:"type:json" json:"completed_challenges"`
	ChallengeProgress   map[string]float64       `gorm:"type:json" json:"challenge_progress"`
	DailyChallenges     []map[string]interface{} `gorm:"type:json" json:"daily_challenges"`
	WeeklyChallenges    []map[string]interface{} `gorm:"type:json" json:"weekly_challenges"`

	// Leaderboards
	GlobalRank         int                      `gorm:"default:0" json:"global_rank"`
	RegionalRank       int                      `gorm:"default:0" json:"regional_rank"`
	FriendsRank        int                      `gorm:"default:0" json:"friends_rank"`
	LeaderboardScores  map[string]int           `gorm:"type:json" json:"leaderboard_scores"`
	CompetitionHistory []map[string]interface{} `gorm:"type:json" json:"competition_history"`

	// Rewards
	UnclaimedRewards []map[string]interface{} `gorm:"type:json" json:"unclaimed_rewards"`
	ClaimedRewards   []map[string]interface{} `gorm:"type:json" json:"claimed_rewards"`
	RewardHistory    []map[string]interface{} `gorm:"type:json" json:"reward_history"`
	TotalRewardValue decimal.Decimal          `gorm:"type:decimal(10,2)" json:"total_reward_value"`
}

// UserEducation tracks educational progress and insurance literacy
type UserEducation struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Insurance Literacy
	InsuranceLiteracyScore float64            `gorm:"default:0" json:"insurance_literacy_score"`
	KnowledgeLevel         string             `gorm:"type:varchar(20)" json:"knowledge_level"`
	CompletedTopics        []string           `gorm:"type:json" json:"completed_topics"`
	PendingTopics          []string           `gorm:"type:json" json:"pending_topics"`
	TopicProgress          map[string]float64 `gorm:"type:json" json:"topic_progress"`

	// Courses & Modules
	EnrolledCourses    []map[string]interface{} `gorm:"type:json" json:"enrolled_courses"`
	CompletedCourses   []map[string]interface{} `gorm:"type:json" json:"completed_courses"`
	CourseProgress     map[string]float64       `gorm:"type:json" json:"course_progress"`
	CertificatesEarned []map[string]interface{} `gorm:"type:json" json:"certificates_earned"`

	// Learning Preferences
	PreferredLearningStyle string            `gorm:"type:varchar(50)" json:"preferred_learning_style"`
	PreferredContentType   string            `gorm:"type:varchar(50)" json:"preferred_content_type"`
	LearningPace           string            `gorm:"type:varchar(20)" json:"learning_pace"`
	StudySchedule          map[string]string `gorm:"type:json" json:"study_schedule"`

	// Assessments
	AssessmentsTaken []map[string]interface{} `gorm:"type:json" json:"assessments_taken"`
	AssessmentScores map[string]float64       `gorm:"type:json" json:"assessment_scores"`
	SkillLevels      map[string]int           `gorm:"type:json" json:"skill_levels"`
	CompetencyMatrix map[string]float64       `gorm:"type:json" json:"competency_matrix"`

	// Resources
	AccessedResources   []string `gorm:"type:json" json:"accessed_resources"`
	BookmarkedContent   []string `gorm:"type:json" json:"bookmarked_content"`
	DownloadedMaterials []string `gorm:"type:json" json:"downloaded_materials"`
	SharedContent       []string `gorm:"type:json" json:"shared_content"`

	// Engagement
	LearningStreak       int        `gorm:"default:0" json:"learning_streak"`
	TotalLearningHours   float64    `gorm:"default:0" json:"total_learning_hours"`
	LastLearningActivity *time.Time `json:"last_learning_activity"`
	EngagementScore      float64    `gorm:"default:0" json:"engagement_score"`

	// Achievements
	EducationMilestones []map[string]interface{} `gorm:"type:json" json:"education_milestones"`
	EducationBadges     []string                 `gorm:"type:json" json:"education_badges"`
	ProficiencyLevels   map[string]string        `gorm:"type:json" json:"proficiency_levels"`
}

// UserSocial manages social interactions and community engagement
type UserSocial struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Social Profile
	Username          string `gorm:"type:varchar(50);uniqueIndex" json:"username"`
	DisplayName       string `gorm:"type:varchar(100)" json:"display_name"`
	Bio               string `gorm:"type:text" json:"bio"`
	Avatar            string `gorm:"type:varchar(255)" json:"avatar"`
	ProfileVisibility string `gorm:"type:varchar(20)" json:"profile_visibility"`
	VerifiedStatus    bool   `gorm:"default:false" json:"verified_status"`

	// Connections
	FollowersCount   int         `gorm:"default:0" json:"followers_count"`
	FollowingCount   int         `gorm:"default:0" json:"following_count"`
	ConnectionsCount int         `gorm:"default:0" json:"connections_count"`
	Followers        []uuid.UUID `gorm:"type:json" json:"followers"`
	Following        []uuid.UUID `gorm:"type:json" json:"following"`
	BlockedUsers     []uuid.UUID `gorm:"type:json" json:"blocked_users"`

	// Community Engagement
	CommunitiesJoined []map[string]interface{} `gorm:"type:json" json:"communities_joined"`
	CommunityRole     map[string]string        `gorm:"type:json" json:"community_role"`
	PostsCreated      int                      `gorm:"default:0" json:"posts_created"`
	CommentsCount     int                      `gorm:"default:0" json:"comments_count"`
	LikesGiven        int                      `gorm:"default:0" json:"likes_given"`
	LikesReceived     int                      `gorm:"default:0" json:"likes_received"`
	SharesCount       int                      `gorm:"default:0" json:"shares_count"`

	// Reputation
	ReputationScore    float64  `gorm:"default:0" json:"reputation_score"`
	HelpfulnessScore   float64  `gorm:"default:0" json:"helpfulness_score"`
	ExpertiseAreas     []string `gorm:"type:json" json:"expertise_areas"`
	ContributionLevel  string   `gorm:"type:varchar(20)" json:"contribution_level"`
	TrustedContributor bool     `gorm:"default:false" json:"trusted_contributor"`

	// Content Sharing
	SharedContent  []map[string]interface{} `gorm:"type:json" json:"shared_content"`
	SavedPosts     []string                 `gorm:"type:json" json:"saved_posts"`
	ContentViews   int                      `gorm:"default:0" json:"content_views"`
	EngagementRate float64                  `gorm:"default:0" json:"engagement_rate"`

	// Groups & Forums
	GroupMemberships   []map[string]interface{} `gorm:"type:json" json:"group_memberships"`
	ForumActivity      map[string]int           `gorm:"type:json" json:"forum_activity"`
	DiscussionsStarted int                      `gorm:"default:0" json:"discussions_started"`
	SolutionsProvided  int                      `gorm:"default:0" json:"solutions_provided"`

	// Social Features
	PrivacySettings      map[string]bool `gorm:"type:json" json:"privacy_settings"`
	NotificationSettings map[string]bool `gorm:"type:json" json:"notification_settings"`
	ActivityStatus       string          `gorm:"type:varchar(20)" json:"activity_status"`
	LastActiveTime       *time.Time      `json:"last_active_time"`
}

// UserRewards manages reward programs and redemptions
type UserRewards struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Reward Points
	TotalPointsEarned int        `gorm:"default:0" json:"total_points_earned"`
	AvailablePoints   int        `gorm:"default:0" json:"available_points"`
	PendingPoints     int        `gorm:"default:0" json:"pending_points"`
	RedeemedPoints    int        `gorm:"default:0" json:"redeemed_points"`
	ExpiringPoints    int        `gorm:"default:0" json:"expiring_points"`
	PointsExpiryDate  *time.Time `json:"points_expiry_date"`

	// Reward Categories
	CashbackEarned  decimal.Decimal `gorm:"type:decimal(10,2)" json:"cashback_earned"`
	DiscountsEarned decimal.Decimal `gorm:"type:decimal(10,2)" json:"discounts_earned"`
	VouchersEarned  int             `gorm:"default:0" json:"vouchers_earned"`
	GiftCardsEarned int             `gorm:"default:0" json:"gift_cards_earned"`

	// Redemption History
	RedemptionHistory    []map[string]interface{} `gorm:"type:json" json:"redemption_history"`
	LastRedemptionDate   *time.Time               `json:"last_redemption_date"`
	TotalRedemptionValue decimal.Decimal          `gorm:"type:decimal(15,2)" json:"total_redemption_value"`
	PreferredRedemption  string                   `gorm:"type:varchar(50)" json:"preferred_redemption"`

	// Catalog Access
	AvailableRewards    []map[string]interface{} `gorm:"type:json" json:"available_rewards"`
	LockedRewards       []map[string]interface{} `gorm:"type:json" json:"locked_rewards"`
	ExclusiveRewards    []map[string]interface{} `gorm:"type:json" json:"exclusive_rewards"`
	PersonalizedRewards []map[string]interface{} `gorm:"type:json" json:"personalized_rewards"`

	// Partner Rewards
	PartnerPoints      map[string]int           `gorm:"type:json" json:"partner_points"`
	PartnerRedemptions []map[string]interface{} `gorm:"type:json" json:"partner_redemptions"`
	TransferablePoints bool                     `gorm:"default:false" json:"transferable_points"`

	// Bonus Programs
	BonusMultiplier     float64                  `gorm:"default:1.0" json:"bonus_multiplier"`
	ActiveBonuses       []map[string]interface{} `gorm:"type:json" json:"active_bonuses"`
	QualifiedPromotions []string                 `gorm:"type:json" json:"qualified_promotions"`
	SpecialOffers       []map[string]interface{} `gorm:"type:json" json:"special_offers"`
}

// UserChallenges manages user challenges and competitions
type UserChallenges struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`

	// Challenge Details
	ChallengeID          string    `gorm:"type:varchar(100)" json:"challenge_id"`
	ChallengeType        string    `gorm:"type:varchar(50)" json:"challenge_type"`
	ChallengeName        string    `gorm:"type:varchar(100)" json:"challenge_name"`
	ChallengeDescription string    `gorm:"type:text" json:"challenge_description"`
	StartDate            time.Time `json:"start_date"`
	EndDate              time.Time `json:"end_date"`
	Status               string    `gorm:"type:varchar(20)" json:"status"`

	// Progress
	CurrentProgress      float64                  `gorm:"default:0" json:"current_progress"`
	TargetGoal           float64                  `gorm:"default:0" json:"target_goal"`
	CompletionPercentage float64                  `gorm:"default:0" json:"completion_percentage"`
	Milestones           []map[string]interface{} `gorm:"type:json" json:"milestones"`
	CheckpointsCompleted int                      `gorm:"default:0" json:"checkpoints_completed"`

	// Competition
	CompetitorCount    int                `gorm:"default:0" json:"competitor_count"`
	CurrentRank        int                `gorm:"default:0" json:"current_rank"`
	BestRank           int                `gorm:"default:0" json:"best_rank"`
	CompetitorProgress map[string]float64 `gorm:"type:json" json:"competitor_progress"`

	// Rewards
	RewardPoints        int                      `gorm:"default:0" json:"reward_points"`
	RewardItems         []map[string]interface{} `gorm:"type:json" json:"reward_items"`
	BonusEarned         decimal.Decimal          `gorm:"type:decimal(10,2)" json:"bonus_earned"`
	AchievementUnlocked []string                 `gorm:"type:json" json:"achievement_unlocked"`

	// Team Challenges
	TeamID           *uuid.UUID `gorm:"type:uuid" json:"team_id"`
	TeamName         string     `gorm:"type:varchar(100)" json:"team_name"`
	TeamRole         string     `gorm:"type:varchar(50)" json:"team_role"`
	TeamContribution float64    `gorm:"default:0" json:"team_contribution"`
	TeamRank         int        `gorm:"default:0" json:"team_rank"`
}

// UserMissions tracks user missions and quests
type UserMissions struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Active Missions
	ActiveMissions    []map[string]interface{} `gorm:"type:json" json:"active_missions"`
	CompletedMissions []map[string]interface{} `gorm:"type:json" json:"completed_missions"`
	FailedMissions    []map[string]interface{} `gorm:"type:json" json:"failed_missions"`
	MissionProgress   map[string]float64       `gorm:"type:json" json:"mission_progress"`

	// Daily/Weekly/Monthly
	DailyMissions   []map[string]interface{} `gorm:"type:json" json:"daily_missions"`
	WeeklyMissions  []map[string]interface{} `gorm:"type:json" json:"weekly_missions"`
	MonthlyMissions []map[string]interface{} `gorm:"type:json" json:"monthly_missions"`
	SpecialMissions []map[string]interface{} `gorm:"type:json" json:"special_missions"`

	// Quest Lines
	ActiveQuests  []map[string]interface{} `gorm:"type:json" json:"active_quests"`
	QuestChapter  int                      `gorm:"default:1" json:"quest_chapter"`
	QuestStage    int                      `gorm:"default:1" json:"quest_stage"`
	StoryProgress float64                  `gorm:"default:0" json:"story_progress"`

	// Achievements
	MissionPoints      int `gorm:"default:0" json:"mission_points"`
	MissionStreak      int `gorm:"default:0" json:"mission_streak"`
	PerfectCompletions int `gorm:"default:0" json:"perfect_completions"`
	SpeedCompletions   int `gorm:"default:0" json:"speed_completions"`

	// Rewards
	MissionRewards  []map[string]interface{} `gorm:"type:json" json:"mission_rewards"`
	BonusRewards    []map[string]interface{} `gorm:"type:json" json:"bonus_rewards"`
	UnlockedContent []string                 `gorm:"type:json" json:"unlocked_content"`
}

// TableName returns the table name
func (UserGamification) TableName() string {
	return "user_gamification"
}

// TableName returns the table name
func (UserEducation) TableName() string {
	return "user_education"
}

// TableName returns the table name
func (UserSocial) TableName() string {
	return "user_social"
}

// TableName returns the table name
func (UserRewards) TableName() string {
	return "user_rewards"
}

// TableName returns the table name
func (UserChallenges) TableName() string {
	return "user_challenges"
}

// TableName returns the table name
func (UserMissions) TableName() string {
	return "user_missions"
}
