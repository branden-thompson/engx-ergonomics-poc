package chaos

import (
	"sync"
	"time"
)

// SkillLevel represents the user's assessed skill level
type SkillLevel int

const (
	Novice SkillLevel = iota
	Intermediate
	Advanced
	Expert
)

// BehaviorTracker tracks and analyzes user behavior patterns
type BehaviorTracker struct {
	sessions          []Session
	currentSession    *Session
	skillLevel        SkillLevel
	competenceMetrics *CompetenceMetrics
	adaptationHistory []AdaptationEvent
	mutex             sync.RWMutex
}

// Session represents a user session with chaos scenarios
type Session struct {
	ID                string
	StartTime         time.Time
	EndTime           time.Time
	Actions           []UserAction
	Scenarios         []ScenarioInteraction
	SkillAssessment   *SkillAssessment
}

// ScenarioInteraction represents how a user interacted with a chaos scenario
type ScenarioInteraction struct {
	ScenarioType      string
	StartTime         time.Time
	EndTime           time.Time
	Resolved          bool
	ResolutionTime    time.Duration
	AttemptsRequired  int
	HelpRequested     bool
	RecoveryPath      string
}

// CompetenceMetrics tracks detailed user competence measurements
type CompetenceMetrics struct {
	SuccessRate           float64       `json:"success_rate"`
	AverageResolutionTime time.Duration `json:"average_resolution_time"`
	HelpRequestFrequency  float64       `json:"help_request_frequency"`
	RetryPatterns         []RetryPattern `json:"retry_patterns"`
	LearningVelocity      float64       `json:"learning_velocity"`
	ConfidenceScore       float64       `json:"confidence_score"`
	ProblemAreas          []string      `json:"problem_areas"`
	Strengths             []string      `json:"strengths"`
}

// RetryPattern represents common user retry behaviors
type RetryPattern struct {
	TriggerScenario   string        `json:"trigger_scenario"`
	AverageRetries    float64       `json:"average_retries"`
	RetryInterval     time.Duration `json:"retry_interval"`
	SuccessfulPattern bool          `json:"successful_pattern"`
}

// AdaptationEvent represents a difficulty adaptation event
type AdaptationEvent struct {
	Timestamp         time.Time          `json:"timestamp"`
	PreviousLevel     AggressivenessLevel `json:"previous_level"`
	NewLevel          AggressivenessLevel `json:"new_level"`
	Reason            string             `json:"reason"`
	UserMetrics       *CompetenceMetrics `json:"user_metrics"`
}

// SkillAssessment represents a comprehensive skill assessment
type SkillAssessment struct {
	CurrentLevel      SkillLevel `json:"current_level"`
	Confidence        float64    `json:"confidence"`
	LearningRate      float64    `json:"learning_rate"`
	ProblemAreas      []string   `json:"problem_areas"`
	Strengths         []string   `json:"strengths"`
	RecommendedLevel  SkillLevel `json:"recommended_level"`
	AssessmentDate    time.Time  `json:"assessment_date"`
}

// SkillTest represents a learning assessment test
type SkillTest struct {
	TestID          string            `json:"test_id"`
	LearningGoals   []string         `json:"learning_goals"`
	Questions       []TestQuestion   `json:"questions"`
	PassingScore    float64          `json:"passing_score"`
	TimeLimit       time.Duration    `json:"time_limit"`
}

// TestQuestion represents a single assessment question
type TestQuestion struct {
	ID              string   `json:"id"`
	Question        string   `json:"question"`
	QuestionType    string   `json:"question_type"` // multiple_choice, practical, explanation
	Options         []string `json:"options,omitempty"`
	CorrectAnswer   string   `json:"correct_answer"`
	Points          int      `json:"points"`
	Category        string   `json:"category"`
}

// ExpectedAction represents an expected user action for recovery
type ExpectedAction struct {
	ActionType      string        `json:"action_type"`
	Description     string        `json:"description"`
	Command         string        `json:"command,omitempty"`
	TimeoutDuration time.Duration `json:"timeout_duration"`
	Required        bool          `json:"required"`
	Points          int           `json:"points"`
}

// RecoveryPath represents alternative recovery approaches
type RecoveryPath struct {
	PathID          string            `json:"path_id"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	Actions         []ExpectedAction  `json:"actions"`
	DifficultyLevel SkillLevel        `json:"difficulty_level"`
	EstimatedTime   time.Duration     `json:"estimated_time"`
}

// ResourceType represents types of resources that can be affected
type ResourceType int

const (
	FileSystem ResourceType = iota
	Network
	Memory
	CPU
	Disk
	Permissions
)

// UserResponse represents how a user responded to a chaos scenario
type UserResponse struct {
	ScenarioID        string            `json:"scenario_id"`
	ResponseStartTime time.Time         `json:"response_start_time"`
	ResponseEndTime   time.Time         `json:"response_end_time"`
	Actions           []UserAction      `json:"actions"`
	Successful        bool              `json:"successful"`
	RecoveryPath      string            `json:"recovery_path"`
	HelpRequested     bool              `json:"help_requested"`
}

// NewBehaviorTracker creates a new behavior tracker
func NewBehaviorTracker() *BehaviorTracker {
	return &BehaviorTracker{
		sessions:          make([]Session, 0),
		currentSession:    nil,
		skillLevel:        Intermediate, // Start with intermediate assumption
		competenceMetrics: &CompetenceMetrics{
			SuccessRate:           0.5, // Neutral starting point
			AverageResolutionTime: 5 * time.Minute,
			HelpRequestFrequency:  0.2,
			RetryPatterns:         make([]RetryPattern, 0),
			LearningVelocity:      0.5,
			ConfidenceScore:       0.5,
			ProblemAreas:          make([]string, 0),
			Strengths:             make([]string, 0),
		},
		adaptationHistory: make([]AdaptationEvent, 0),
	}
}

// StartSession starts a new user session
func (bt *BehaviorTracker) StartSession() string {
	bt.mutex.Lock()
	defer bt.mutex.Unlock()

	sessionID := generateSessionID()
	session := &Session{
		ID:        sessionID,
		StartTime: time.Now(),
		Actions:   make([]UserAction, 0),
		Scenarios: make([]ScenarioInteraction, 0),
	}

	bt.currentSession = session
	return sessionID
}

// EndSession ends the current session
func (bt *BehaviorTracker) EndSession() {
	bt.mutex.Lock()
	defer bt.mutex.Unlock()

	if bt.currentSession != nil {
		bt.currentSession.EndTime = time.Now()
		bt.sessions = append(bt.sessions, *bt.currentSession)
		bt.currentSession = nil

		// Update competence metrics based on session
		bt.updateCompetenceMetrics()
	}
}

// RecordAction records a user action
func (bt *BehaviorTracker) RecordAction(action UserAction) error {
	bt.mutex.Lock()
	defer bt.mutex.Unlock()

	if bt.currentSession == nil {
		// Auto-start session if none exists
		bt.StartSession()
	}

	bt.currentSession.Actions = append(bt.currentSession.Actions, action)

	// Update metrics in real-time
	bt.updateRealTimeMetrics(action)

	return nil
}

// GetCurrentPattern returns the current behavior pattern
func (bt *BehaviorTracker) GetCurrentPattern() *BehaviorPattern {
	bt.mutex.RLock()
	defer bt.mutex.RUnlock()

	// Calculate recent success rate (last 10 actions)
	recentSuccessRate := bt.calculateRecentSuccessRate(10)

	// Calculate average resolution time
	avgResolutionTime := bt.calculateAverageResolutionTime()

	// Calculate help request frequency
	helpFrequency := bt.calculateHelpRequestFrequency()

	// Determine if user shows frustration
	showsFrustration := bt.detectFrustrationPattern()

	// Calculate confidence level
	confidenceLevel := bt.calculateConfidenceLevel()

	return &BehaviorPattern{
		SkillLevel:            bt.skillLevel,
		RecentSuccessRate:     recentSuccessRate,
		AverageResolutionTime: avgResolutionTime,
		HelpRequestFrequency:  helpFrequency,
		RetryPatterns:         bt.competenceMetrics.RetryPatterns,
		ShowsFrustration:      showsFrustration,
		ConfidenceLevel:       confidenceLevel,
	}
}

// calculateRecentSuccessRate calculates success rate for recent actions
func (bt *BehaviorTracker) calculateRecentSuccessRate(count int) float64 {
	if bt.currentSession == nil || len(bt.currentSession.Actions) == 0 {
		return bt.competenceMetrics.SuccessRate
	}

	actions := bt.currentSession.Actions
	start := len(actions) - count
	if start < 0 {
		start = 0
	}

	recentActions := actions[start:]
	if len(recentActions) == 0 {
		return bt.competenceMetrics.SuccessRate
	}

	successCount := 0
	for _, action := range recentActions {
		if action.Success {
			successCount++
		}
	}

	return float64(successCount) / float64(len(recentActions))
}

// calculateAverageResolutionTime calculates average time to resolve issues
func (bt *BehaviorTracker) calculateAverageResolutionTime() time.Duration {
	if bt.currentSession == nil || len(bt.currentSession.Scenarios) == 0 {
		return bt.competenceMetrics.AverageResolutionTime
	}

	totalTime := time.Duration(0)
	resolvedCount := 0

	for _, scenario := range bt.currentSession.Scenarios {
		if scenario.Resolved {
			totalTime += scenario.ResolutionTime
			resolvedCount++
		}
	}

	if resolvedCount == 0 {
		return bt.competenceMetrics.AverageResolutionTime
	}

	return totalTime / time.Duration(resolvedCount)
}

// calculateHelpRequestFrequency calculates how often user requests help
func (bt *BehaviorTracker) calculateHelpRequestFrequency() float64 {
	if bt.currentSession == nil || len(bt.currentSession.Actions) == 0 {
		return bt.competenceMetrics.HelpRequestFrequency
	}

	helpRequests := 0
	totalActions := len(bt.currentSession.Actions)

	for _, action := range bt.currentSession.Actions {
		if action.ActionType == HelpRequest {
			helpRequests++
		}
	}

	return float64(helpRequests) / float64(totalActions)
}

// detectFrustrationPattern detects if user shows signs of frustration
func (bt *BehaviorTracker) detectFrustrationPattern() bool {
	if bt.currentSession == nil || len(bt.currentSession.Actions) < 5 {
		return false
	}

	// Look for patterns indicating frustration:
	// 1. Multiple rapid retry attempts
	// 2. Frequent help requests
	// 3. Long gaps between actions (confusion)

	recentActions := bt.getRecentActions(10)
	if len(recentActions) < 5 {
		return false
	}

	// Check for rapid retries
	rapidRetries := 0
	for i := 1; i < len(recentActions); i++ {
		timeDiff := recentActions[i].Timestamp.Sub(recentActions[i-1].Timestamp)
		if timeDiff < 5*time.Second && recentActions[i].ActionType == RetryAttempt {
			rapidRetries++
		}
	}

	// Check help request frequency in recent actions
	helpRequests := 0
	for _, action := range recentActions {
		if action.ActionType == HelpRequest {
			helpRequests++
		}
	}

	// Frustration indicators
	return rapidRetries >= 3 || float64(helpRequests)/float64(len(recentActions)) > 0.4
}

// calculateConfidenceLevel calculates user confidence based on behavior
func (bt *BehaviorTracker) calculateConfidenceLevel() float64 {
	if bt.currentSession == nil {
		return bt.competenceMetrics.ConfidenceScore
	}

	// Factors that indicate confidence:
	// - Success rate
	// - Speed of resolution
	// - Low help request frequency
	// - Consistent performance

	successRate := bt.calculateRecentSuccessRate(10)
	helpFrequency := bt.calculateHelpRequestFrequency()

	// Base confidence on success rate
	confidence := successRate

	// Adjust for help request frequency (lower frequency = higher confidence)
	confidence *= (1.0 - helpFrequency*0.5)

	// Adjust for consistency (less variance = higher confidence)
	confidence *= bt.calculateConsistencyFactor()

	return clampFloat64(confidence, 0.0, 1.0)
}

// calculateConsistencyFactor calculates how consistent user performance is
func (bt *BehaviorTracker) calculateConsistencyFactor() float64 {
	if bt.currentSession == nil || len(bt.currentSession.Actions) < 5 {
		return 0.5 // Neutral for insufficient data
	}

	recentActions := bt.getRecentActions(10)
	if len(recentActions) < 5 {
		return 0.5
	}

	// Calculate variance in action success
	successes := make([]float64, len(recentActions))
	for i, action := range recentActions {
		if action.Success {
			successes[i] = 1.0
		} else {
			successes[i] = 0.0
		}
	}

	variance := calculateVariance(successes)

	// Lower variance = higher consistency = higher confidence factor
	// Variance ranges from 0 (perfect consistency) to 0.25 (maximum inconsistency)
	consistencyFactor := 1.0 - (variance * 4.0) // Scale to 0-1 range

	return clampFloat64(consistencyFactor, 0.0, 1.0)
}

// getRecentActions gets the most recent actions
func (bt *BehaviorTracker) getRecentActions(count int) []UserAction {
	if bt.currentSession == nil || len(bt.currentSession.Actions) == 0 {
		return make([]UserAction, 0)
	}

	actions := bt.currentSession.Actions
	start := len(actions) - count
	if start < 0 {
		start = 0
	}

	return actions[start:]
}

// updateRealTimeMetrics updates metrics as actions are recorded
func (bt *BehaviorTracker) updateRealTimeMetrics(action UserAction) {
	// Update skill assessment based on action patterns
	bt.updateSkillAssessment(action)

	// Update competence metrics
	bt.updateCompetenceMetrics()
}

// updateSkillAssessment updates the user's skill assessment
func (bt *BehaviorTracker) updateSkillAssessment(action UserAction) {
	// Simple skill level assessment based on action patterns
	// This would be more sophisticated in a full implementation

	recentSuccessRate := bt.calculateRecentSuccessRate(5)
	avgResolutionTime := bt.calculateAverageResolutionTime()
	helpFrequency := bt.calculateHelpRequestFrequency()

	// Skill level assessment logic
	if recentSuccessRate > 0.9 && avgResolutionTime < 2*time.Minute && helpFrequency < 0.1 {
		bt.skillLevel = Expert
	} else if recentSuccessRate > 0.7 && avgResolutionTime < 5*time.Minute && helpFrequency < 0.2 {
		bt.skillLevel = Advanced
	} else if recentSuccessRate > 0.5 && helpFrequency < 0.4 {
		bt.skillLevel = Intermediate
	} else {
		bt.skillLevel = Novice
	}
}

// updateCompetenceMetrics updates overall competence metrics
func (bt *BehaviorTracker) updateCompetenceMetrics() {
	if bt.currentSession == nil {
		return
	}

	// Update success rate
	bt.competenceMetrics.SuccessRate = bt.calculateRecentSuccessRate(20)

	// Update average resolution time
	bt.competenceMetrics.AverageResolutionTime = bt.calculateAverageResolutionTime()

	// Update help request frequency
	bt.competenceMetrics.HelpRequestFrequency = bt.calculateHelpRequestFrequency()

	// Update confidence score
	bt.competenceMetrics.ConfidenceScore = bt.calculateConfidenceLevel()

	// Update learning velocity (simplified calculation)
	bt.updateLearningVelocity()
}

// updateLearningVelocity calculates how quickly the user is improving
func (bt *BehaviorTracker) updateLearningVelocity() {
	if len(bt.sessions) < 2 {
		bt.competenceMetrics.LearningVelocity = 0.5 // Neutral
		return
	}

	// Compare current session metrics with previous sessions
	currentSuccess := bt.calculateRecentSuccessRate(10)

	// Get average success rate from previous sessions
	totalPreviousSuccess := 0.0
	validSessions := 0

	for i := len(bt.sessions) - 3; i < len(bt.sessions)-1 && i >= 0; i++ {
		session := bt.sessions[i]
		if len(session.Actions) > 0 {
			sessionSuccess := calculateSessionSuccessRate(session)
			totalPreviousSuccess += sessionSuccess
			validSessions++
		}
	}

	if validSessions == 0 {
		bt.competenceMetrics.LearningVelocity = 0.5
		return
	}

	avgPreviousSuccess := totalPreviousSuccess / float64(validSessions)
	improvement := currentSuccess - avgPreviousSuccess

	// Convert improvement to velocity score (0.0 to 1.0)
	velocity := 0.5 + improvement // 0.5 is neutral, above is positive, below is negative
	bt.competenceMetrics.LearningVelocity = clampFloat64(velocity, 0.0, 1.0)
}

// Reset resets the behavior tracker
func (bt *BehaviorTracker) Reset() {
	bt.mutex.Lock()
	defer bt.mutex.Unlock()

	bt.sessions = make([]Session, 0)
	bt.currentSession = nil
	bt.skillLevel = Intermediate
	bt.competenceMetrics = &CompetenceMetrics{
		SuccessRate:           0.5,
		AverageResolutionTime: 5 * time.Minute,
		HelpRequestFrequency:  0.2,
		RetryPatterns:         make([]RetryPattern, 0),
		LearningVelocity:      0.5,
		ConfidenceScore:       0.5,
		ProblemAreas:          make([]string, 0),
		Strengths:             make([]string, 0),
	}
	bt.adaptationHistory = make([]AdaptationEvent, 0)
}

// Helper functions

// generateSessionID generates a unique session ID
func generateSessionID() string {
	return time.Now().Format("20060102-150405") + "-session"
}

// calculateSessionSuccessRate calculates success rate for a completed session
func calculateSessionSuccessRate(session Session) float64 {
	if len(session.Actions) == 0 {
		return 0.0
	}

	successCount := 0
	for _, action := range session.Actions {
		if action.Success {
			successCount++
		}
	}

	return float64(successCount) / float64(len(session.Actions))
}

// calculateVariance calculates variance of a slice of float64 values
func calculateVariance(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}

	// Calculate mean
	sum := 0.0
	for _, value := range values {
		sum += value
	}
	mean := sum / float64(len(values))

	// Calculate variance
	varianceSum := 0.0
	for _, value := range values {
		diff := value - mean
		varianceSum += diff * diff
	}

	return varianceSum / float64(len(values))
}

// clampFloat64 ensures a float64 value is within the specified range
func clampFloat64(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// String methods for enums

func (sl SkillLevel) String() string {
	switch sl {
	case Novice:
		return "novice"
	case Intermediate:
		return "intermediate"
	case Advanced:
		return "advanced"
	case Expert:
		return "expert"
	default:
		return "unknown"
	}
}

func (rt ResourceType) String() string {
	switch rt {
	case FileSystem:
		return "filesystem"
	case Network:
		return "network"
	case Memory:
		return "memory"
	case CPU:
		return "cpu"
	case Disk:
		return "disk"
	case Permissions:
		return "permissions"
	default:
		return "unknown"
	}
}