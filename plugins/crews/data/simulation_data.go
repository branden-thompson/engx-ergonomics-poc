package data

import (
	"fmt"
	"strings"
	"time"

	"github.com/bthompso/engx-ergonomics-poc/plugins/crews/models"
)

// SimulationDataStore provides in-memory data for crew command simulation
type SimulationDataStore struct {
	crews       map[string]*models.Crew
	assets      map[string]*models.AssetOwnership
	userIndex   map[string][]string // user -> crew IDs
	assetIndex  map[string]string   // asset -> crew ID
}

// NewSimulationDataStore creates and populates a new data store
func NewSimulationDataStore() *SimulationDataStore {
	store := &SimulationDataStore{
		crews:      make(map[string]*models.Crew),
		assets:     make(map[string]*models.AssetOwnership),
		userIndex:  make(map[string][]string),
		assetIndex: make(map[string]string),
	}

	store.populateData()
	return store
}

// populateData creates realistic simulation data
func (s *SimulationDataStore) populateData() {
	now := time.Now()

	// Create sample crews
	crews := []*models.Crew{
		{
			ID:          "CREW-1234",
			VanityName:  "Web Platform Team",
			Description: "Manages web application platform infrastructure and frontend services",
			CreatedAt:   now.AddDate(0, -6, 0),
			CreatedBy:   "bthompso",
			Status:      models.CrewStatusActive,
			LastModified: now.AddDate(0, 0, -5),
			ModifiedBy:  "bthompso",
			Members: []models.Member{
				{
					UserID:     "bthompso",
					Email:      "bthompso@company.com",
					FullName:   "Branden Thompson",
					Level:      "IC5",
					Role:       models.RoleOwner,
					JoinedAt:   now.AddDate(0, -6, 0),
					AddedBy:    "system",
					IsOnCall:   false,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 2),
				},
				{
					UserID:     "jlawrenc",
					Email:      "jlawrenc@company.com",
					FullName:   "Josh Lawrence",
					Level:      "MR3",
					Role:       models.RoleAdmin,
					JoinedAt:   now.AddDate(0, -4, 0),
					AddedBy:    "bthompso",
					IsOnCall:   false,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 6),
				},
				{
					UserID:     "osurtiz",
					Email:      "osurtiz@company.com",
					FullName:   "Olivier Surtiz",
					Level:      "IC5",
					Role:       models.RoleAdmin,
					JoinedAt:   now.AddDate(0, -3, 0),
					AddedBy:    "bthompso",
					IsOnCall:   true,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 4),
				},
				{
					UserID:     "hbacot",
					Email:      "hbacot@company.com",
					FullName:   "Hunter Bacot",
					Level:      "IC3",
					Role:       models.RoleMember,
					JoinedAt:   now.AddDate(0, -2, 0),
					AddedBy:    "jlawrenc",
					IsOnCall:   true,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 8),
				},
				{
					UserID:     "mlondhe",
					Email:      "mlondhe@company.com",
					FullName:   "Mruganka Rajendra Londhe",
					Level:      "IC3",
					Role:       models.RoleMember,
					JoinedAt:   now.AddDate(0, -3, 0),
					AddedBy:    "jlawrenc",
					IsOnCall:   false,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 16),
				},
				{
					UserID:     "pmudugan",
					Email:      "pmudugan@company.com",
					FullName:   "Pooja Muduganti",
					Level:      "IC4",
					Role:       models.RoleMember,
					JoinedAt:   now.AddDate(0, -2, 0),
					AddedBy:    "bthompso",
					IsOnCall:   false,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 12),
				},
				{
					UserID:     "cheli",
					Email:      "cheli@company.com",
					FullName:   "Chennan Li",
					Level:      "IC4",
					Role:       models.RoleMember,
					JoinedAt:   now.AddDate(0, -1, 0),
					AddedBy:    "osurtiz",
					IsOnCall:   false,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 20),
				},
				{
					UserID:     "pvalluri",
					Email:      "pvalluri@company.com",
					FullName:   "Prince Valluri",
					Level:      "IC6",
					Role:       models.RoleAuto,
					JoinedAt:   now.AddDate(0, -4, 0),
					AddedBy:    "system",
					IsOnCall:   false,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 36),
				},
				{
					UserID:     "kamandaliya",
					Email:      "kamandaliya@company.com",
					FullName:   "Kavya Mandaliya",
					Level:      "IC3",
					Role:       models.RoleAuto,
					JoinedAt:   now.AddDate(0, -5, 0),
					AddedBy:    "system",
					IsOnCall:   false,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 48),
				},
				{
					UserID:     "ayegappan",
					Email:      "ayegappan@company.com",
					FullName:   "Arun Yeggapan",
					Level:      "IC4",
					Role:       models.RoleAuto,
					JoinedAt:   now.AddDate(0, -3, 0),
					AddedBy:    "system",
					IsOnCall:   false,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 28),
				},
				{
					UserID:     "tempUser",
					Email:      "tempUser@company.com",
					FullName:   "Sample Temp Member",
					Level:      "IC2",
					Role:       models.RoleTemp,
					JoinedAt:   now.AddDate(0, -1, 0),
					AddedBy:    "bthompso",
					IsOnCall:   false,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 72),
				},
			},
			OnCallSchedule: models.OnCallSchedule{
				CurrentOnCall:  []string{"osurtiz", "hbacot"},
				RotationType:   "weekly",
				EscalationPath: []string{"jlawrenc", "bthompso"},
				Enabled:        true,
				LastUpdated:    now.Add(-time.Hour * 24),
				UpdatedBy:      "bthompso",
			},
			OwnedAssets: []string{
				"asset://web-app/dashboard",
				"asset://web-app/admin-portal",
				"asset://service/api-gateway",
				"asset://service/user-auth",
				"asset://database/user-store",
				"asset://database/session-store",
			},
		},
		{
			ID:          "CREW-2345",
			VanityName:  "Security Response Team",
			Description: "24/7 security monitoring and incident response for all company infrastructure",
			CreatedAt:   now.AddDate(-1, 0, 0),
			CreatedBy:   "sec.admin",
			Status:      models.CrewStatusActive,
			Members: []models.Member{
				{
					UserID:     "sec.admin",
					Email:      "sec.admin@company.com",
					FullName:   "Security Admin",
					Level:      "MR4",
					Role:       models.RoleOwner,
					JoinedAt:   now.AddDate(-1, 0, 0),
					AddedBy:    "system",
					IsOnCall:   false,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 1),
				},
				{
					UserID:     "bthompso",
					Email:      "bthompso@company.com",
					FullName:   "Branden Thompson",
					Level:      "IC5",
					Role:       models.RoleMember,
					JoinedAt:   now.AddDate(0, -3, 0),
					AddedBy:    "sec.admin",
					IsOnCall:   false,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 2),
				},
				{
					UserID:     "security.analyst",
					Email:      "security.analyst@company.com",
					FullName:   "Security Analyst",
					Level:      "IC4",
					Role:       models.RoleAdmin,
					JoinedAt:   now.AddDate(0, -8, 0),
					AddedBy:    "sec.admin",
					IsOnCall:   true,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Minute * 30),
				},
			},
			OnCallSchedule: models.OnCallSchedule{
				CurrentOnCall:  []string{"security.analyst"},
				RotationType:   "daily",
				EscalationPath: []string{"sec.admin"},
				Enabled:        true,
				LastUpdated:    now.Add(-time.Hour * 12),
				UpdatedBy:      "sec.admin",
			},
			OwnedAssets: []string{
				"asset://service/security-scanner",
				"asset://service/threat-detection",
				"asset://database/security-logs",
				"asset://service/incident-management",
			},
		},
		{
			ID:          "CREW-3456",
			VanityName:  "Database Operations",
			Description: "Database administration, performance optimization, and data infrastructure management",
			CreatedAt:   now.AddDate(0, -10, 0),
			CreatedBy:   "dba.lead",
			Status:      models.CrewStatusActive,
			Members: []models.Member{
				{
					UserID:     "dba.lead",
					Email:      "dba.lead@company.com",
					FullName:   "DBA Lead",
					Level:      "MR5",
					Role:       models.RoleOwner,
					JoinedAt:   now.AddDate(0, -10, 0),
					AddedBy:    "system",
					IsOnCall:   false,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 4),
				},
				{
					UserID:     "bthompso",
					Email:      "bthompso@company.com",
					FullName:   "Branden Thompson",
					Level:      "IC5",
					Role:       models.RoleAdmin,
					JoinedAt:   now.AddDate(0, -5, 0),
					AddedBy:    "dba.lead",
					IsOnCall:   false,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 2),
				},
			},
			OnCallSchedule: models.OnCallSchedule{
				CurrentOnCall:  []string{"dba.lead"},
				RotationType:   "weekly",
				EscalationPath: []string{"bthompso"},
				Enabled:        true,
				LastUpdated:    now.Add(-time.Hour * 48),
				UpdatedBy:      "dba.lead",
			},
			OwnedAssets: []string{
				"asset://database/primary-cluster",
				"asset://database/analytics-warehouse",
				"asset://database/backup-systems",
				"asset://service/db-monitoring",
			},
		},
		{
			ID:          "CREW-4567",
			VanityName:  "Mobile Development",
			Description: "iOS and Android application development and deployment",
			CreatedAt:   now.AddDate(0, -3, 0),
			CreatedBy:   "mobile.lead",
			Status:      models.CrewStatusActive,
			Members: []models.Member{
				{
					UserID:     "mobile.lead",
					Email:      "mobile.lead@company.com",
					FullName:   "Mobile Lead",
					Level:      "MR3",
					Role:       models.RoleOwner,
					JoinedAt:   now.AddDate(0, -3, 0),
					AddedBy:    "system",
					IsOnCall:   false,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 8),
				},
				{
					UserID:     "ios.dev",
					Email:      "ios.dev@company.com",
					FullName:   "iOS Developer",
					Level:      "IC4",
					Role:       models.RoleMember,
					JoinedAt:   now.AddDate(0, -2, 0),
					AddedBy:    "mobile.lead",
					IsOnCall:   false,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 16),
				},
			},
			OnCallSchedule: models.OnCallSchedule{
				CurrentOnCall:  []string{},
				RotationType:   "none",
				EscalationPath: []string{"mobile.lead"},
				Enabled:        false,
				LastUpdated:    now.Add(-time.Hour * 72),
				UpdatedBy:      "mobile.lead",
			},
			OwnedAssets: []string{
				"asset://mobile-app/ios-app",
				"asset://mobile-app/android-app",
				"asset://service/mobile-backend",
			},
		},
		{
			ID:          "CREW-0014",
			VanityName:  "Universal Ownership",
			Description: "Cross-functional ownership coordination and governance",
			CreatedAt:   now.AddDate(0, -8, 0),
			CreatedBy:   "governance.lead",
			Status:      models.CrewStatusActive,
			Members: []models.Member{
				{
					UserID:     "governance.lead",
					Email:      "governance.lead@company.com",
					FullName:   "Governance Lead",
					Level:      "MR6",
					Role:       models.RoleOwner,
					JoinedAt:   now.AddDate(0, -8, 0),
					AddedBy:    "system",
					IsOnCall:   false,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 1),
				},
				{
					UserID:     "bthompso",
					Email:      "bthompso@company.com",
					FullName:   "Branden Thompson",
					Level:      "IC5",
					Role:       models.RoleAdmin,
					JoinedAt:   now.AddDate(0, -6, 0),
					AddedBy:    "governance.lead",
					IsOnCall:   true, // Primary on-call
					OnCallType: models.OnCallPrimary,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 2),
				},
			},
			OnCallSchedule: models.OnCallSchedule{
				CurrentOnCall:  []string{"bthompso"},
				RotationType:   "monthly",
				EscalationPath: []string{"governance.lead"},
				Enabled:        true,
				LastUpdated:    now.Add(-time.Hour * 48),
				UpdatedBy:      "governance.lead",
			},
			OwnedAssets: []string{
				"asset://governance/ownership-matrix",
				"asset://governance/escalation-paths",
			},
		},
		{
			ID:          "CREW-2345",
			VanityName:  "Design Foundations",
			Description: "Core design system components and standards",
			CreatedAt:   now.AddDate(0, -7, 0),
			CreatedBy:   "design.lead",
			Status:      models.CrewStatusActive,
			Members: []models.Member{
				{
					UserID:     "design.lead",
					Email:      "design.lead@company.com",
					FullName:   "Design Lead",
					Level:      "IC7",
					Role:       models.RoleOwner,
					JoinedAt:   now.AddDate(0, -7, 0),
					AddedBy:    "system",
					IsOnCall:   false,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 3),
				},
				{
					UserID:     "bthompso",
					Email:      "bthompso@company.com",
					FullName:   "Branden Thompson",
					Level:      "IC5",
					Role:       models.RoleMember,
					JoinedAt:   now.AddDate(0, -5, 0),
					AddedBy:    "design.lead",
					IsOnCall:   false,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 2),
				},
			},
			OnCallSchedule: models.OnCallSchedule{
				CurrentOnCall:  []string{},
				RotationType:   "none",
				EscalationPath: []string{"design.lead"},
				Enabled:        false,
				LastUpdated:    now.Add(-time.Hour * 168),
				UpdatedBy:      "design.lead",
			},
			OwnedAssets: []string{
				"asset://design/foundation-tokens",
				"asset://design/component-library",
			},
		},
		{
			ID:          "CREW-4567",
			VanityName:  "Design Systems",
			Description: "Advanced design system tooling and implementation",
			CreatedAt:   now.AddDate(0, -6, 0),
			CreatedBy:   "design.systems.lead",
			Status:      models.CrewStatusActive,
			Members: []models.Member{
				{
					UserID:     "design.systems.lead",
					Email:      "design.systems.lead@company.com",
					FullName:   "Design Systems Lead",
					Level:      "IC6",
					Role:       models.RoleOwner,
					JoinedAt:   now.AddDate(0, -6, 0),
					AddedBy:    "system",
					IsOnCall:   false,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 1),
				},
				{
					UserID:     "bthompso",
					Email:      "bthompso@company.com",
					FullName:   "Branden Thompson",
					Level:      "IC5",
					Role:       models.RoleAdmin,
					JoinedAt:   now.AddDate(0, -4, 0),
					AddedBy:    "design.systems.lead",
					IsOnCall:   false,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 2),
				},
			},
			OnCallSchedule: models.OnCallSchedule{
				CurrentOnCall:  []string{},
				RotationType:   "none",
				EscalationPath: []string{"design.systems.lead"},
				Enabled:        false,
				LastUpdated:    now.Add(-time.Hour * 96),
				UpdatedBy:      "design.systems.lead",
			},
			OwnedAssets: []string{
				"asset://design/advanced-components",
				"asset://design/automation-tools",
			},
		},
		{
			ID:          "CREW-5678",
			VanityName:  "Design Engineering",
			Description: "Bridge between design and engineering implementation",
			CreatedAt:   now.AddDate(0, -5, 0),
			CreatedBy:   "design.eng.lead",
			Status:      models.CrewStatusActive,
			Members: []models.Member{
				{
					UserID:     "design.eng.lead",
					Email:      "design.eng.lead@company.com",
					FullName:   "Design Engineering Lead",
					Level:      "IC8",
					Role:       models.RoleOwner,
					JoinedAt:   now.AddDate(0, -5, 0),
					AddedBy:    "system",
					IsOnCall:   false,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 2),
				},
				{
					UserID:     "bthompso",
					Email:      "bthompso@company.com",
					FullName:   "Branden Thompson",
					Level:      "IC5",
					Role:       models.RoleAuto,
					JoinedAt:   now.AddDate(0, -3, 0),
					AddedBy:    "design.eng.lead",
					IsOnCall:   true, // Secondary on-call
					OnCallType: models.OnCallBackup,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 2),
				},
			},
			OnCallSchedule: models.OnCallSchedule{
				CurrentOnCall:  []string{"bthompso"},
				RotationType:   "bi-weekly",
				EscalationPath: []string{"design.eng.lead"},
				Enabled:        true,
				LastUpdated:    now.Add(-time.Hour * 72),
				UpdatedBy:      "design.eng.lead",
			},
			OwnedAssets: []string{
				"asset://design/implementation-guides",
				"asset://design/code-generation-tools",
			},
		},
	}

	// Store crews and build indices
	for _, crew := range crews {
		s.crews[crew.ID] = crew

		// Build user index
		for _, member := range crew.Members {
			if member.IsActive() {
				s.userIndex[member.UserID] = append(s.userIndex[member.UserID], crew.ID)
				// Also index by email
				emailUser := strings.Split(member.Email, "@")[0]
				s.userIndex[emailUser] = append(s.userIndex[emailUser], crew.ID)
			}
		}

		// Build asset index
		for _, assetURN := range crew.OwnedAssets {
			s.assetIndex[assetURN] = crew.ID
		}
	}

	// Create asset ownership records
	allAssets := []string{
		"asset://web-app/dashboard",
		"asset://web-app/admin-portal",
		"asset://service/api-gateway",
		"asset://service/user-auth",
		"asset://database/user-store",
		"asset://database/session-store",
		"asset://service/security-scanner",
		"asset://service/threat-detection",
		"asset://database/security-logs",
		"asset://service/incident-management",
		"asset://database/primary-cluster",
		"asset://database/analytics-warehouse",
		"asset://database/backup-systems",
		"asset://service/db-monitoring",
		"asset://mobile-app/ios-app",
		"asset://mobile-app/android-app",
		"asset://service/mobile-backend",
	}

	for _, assetURN := range allAssets {
		if ownerCrewID, exists := s.assetIndex[assetURN]; exists {
			// Extract asset name and type from URN
			parts := strings.Split(assetURN, "/")
			assetType := strings.TrimPrefix(parts[0], "asset://")
			assetName := strings.Join(parts[1:], "/")

			s.assets[assetURN] = &models.AssetOwnership{
				AssetURN:     assetURN,
				AssetName:    assetName,
				AssetType:    assetType,
				OwnerCrewID:  ownerCrewID,
				CreatedAt:    now.AddDate(0, -6, 0),
				LastModified: now.Add(-time.Hour * 24),
				ModifiedBy:   "system",
				Status:       models.AssetActive,
				Metadata: map[string]string{
					"environment": "production",
					"criticality": "high",
				},
			}
		}
	}
}

// GetCrew retrieves a crew by ID
func (s *SimulationDataStore) GetCrew(id string) (*models.Crew, error) {
	crew, exists := s.crews[id]
	if !exists {
		return nil, fmt.Errorf("crew %s not found", id)
	}
	return crew, nil
}

// GetCrewsByUser retrieves all crews for a user
func (s *SimulationDataStore) GetCrewsByUser(userID string) ([]*models.Crew, error) {
	crewIDs, exists := s.userIndex[userID]
	if !exists {
		return []*models.Crew{}, nil
	}

	var crews []*models.Crew
	for _, crewID := range crewIDs {
		if crew, exists := s.crews[crewID]; exists {
			crews = append(crews, crew)
		}
	}

	return crews, nil
}

// GetCrewByAsset retrieves the owning crew for an asset
func (s *SimulationDataStore) GetCrewByAsset(assetURN string) (*models.Crew, error) {
	crewID, exists := s.assetIndex[assetURN]
	if !exists {
		return nil, fmt.Errorf("asset %s not found or has no owner", assetURN)
	}

	return s.GetCrew(crewID)
}

// GetAssetsByCrewID retrieves all assets owned by a crew
func (s *SimulationDataStore) GetAssetsByCrewID(crewID string) ([]*models.AssetOwnership, error) {
	var assets []*models.AssetOwnership
	for _, asset := range s.assets {
		if asset.OwnerCrewID == crewID {
			assets = append(assets, asset)
		}
	}
	return assets, nil
}

// GetAsset retrieves asset information
func (s *SimulationDataStore) GetAsset(assetURN string) (*models.AssetOwnership, error) {
	asset, exists := s.assets[assetURN]
	if !exists {
		return nil, fmt.Errorf("asset %s not found", assetURN)
	}
	return asset, nil
}

// SearchCrews searches crews by name or description
func (s *SimulationDataStore) SearchCrews(query string) ([]*models.Crew, error) {
	var results []*models.Crew
	query = strings.ToLower(query)

	for _, crew := range s.crews {
		if strings.Contains(strings.ToLower(crew.VanityName), query) ||
			strings.Contains(strings.ToLower(crew.Description), query) ||
			strings.Contains(strings.ToLower(crew.ID), query) {
			results = append(results, crew)
		}
	}

	return results, nil
}

// GetAllCrews retrieves all crews
func (s *SimulationDataStore) GetAllCrews() ([]*models.Crew, error) {
	var crews []*models.Crew
	for _, crew := range s.crews {
		crews = append(crews, crew)
	}
	return crews, nil
}

// CreateCrew adds a new crew (simulation)
func (s *SimulationDataStore) CreateCrew(crew *models.Crew) error {
	if _, exists := s.crews[crew.ID]; exists {
		return fmt.Errorf("crew %s already exists", crew.ID)
	}

	s.crews[crew.ID] = crew

	// Update indices
	for _, member := range crew.Members {
		if member.IsActive() {
			s.userIndex[member.UserID] = append(s.userIndex[member.UserID], crew.ID)
		}
	}

	for _, assetURN := range crew.OwnedAssets {
		s.assetIndex[assetURN] = crew.ID
	}

	return nil
}

// UpdateCrew updates an existing crew (simulation)
func (s *SimulationDataStore) UpdateCrew(crew *models.Crew) error {
	if _, exists := s.crews[crew.ID]; !exists {
		return fmt.Errorf("crew %s not found", crew.ID)
	}

	s.crews[crew.ID] = crew
	return nil
}

// Helper function to normalize user identifiers
func (s *SimulationDataStore) normalizeUserID(userID string) string {
	// Handle email addresses by extracting username
	if strings.Contains(userID, "@") {
		return strings.Split(userID, "@")[0]
	}
	return userID
}