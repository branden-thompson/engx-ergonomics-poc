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
	catalog     map[string]*models.CatalogAsset  // Catalog assets by ID
	userIndex   map[string][]string // user -> crew IDs
	assetIndex  map[string]string   // asset -> crew ID
	catalogURNs map[string]string   // URN -> catalog asset ID
	catalogNames map[string]string  // vanity name -> catalog asset ID
}

// NewSimulationDataStore creates and populates a new data store
func NewSimulationDataStore() *SimulationDataStore {
	store := &SimulationDataStore{
		crews:        make(map[string]*models.Crew),
		assets:       make(map[string]*models.AssetOwnership),
		catalog:      make(map[string]*models.CatalogAsset),
		userIndex:    make(map[string][]string),
		assetIndex:   make(map[string]string),
		catalogURNs:  make(map[string]string),
		catalogNames: make(map[string]string),
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
			Type:        models.CrewTypeStandard,
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
					IsOnCall:   true, // Secondary on-call
					OnCallType: models.OnCallSecondary,
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
					IsOnCall:   true, // Primary on-call
					OnCallType: models.OnCallPrimary,
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
				Rotations: []models.OnCallRotation{
					{
						Name:           "Primary",
						OnCallMember:   "hbacot",
						RotationStarts: now.AddDate(0, 0, -3), // Started 3 days ago
						RotationEnds:   now.AddDate(0, 0, 4),  // Ends in 4 days
					},
					{
						Name:           "Secondary",
						OnCallMember:   "osurtiz",
						RotationStarts: now.AddDate(0, 0, -3), // Started 3 days ago
						RotationEnds:   now.AddDate(0, 0, 4),  // Ends in 4 days
					},
				},
				Enabled:     true,
				LastUpdated: now.Add(-time.Hour * 24),
				UpdatedBy:   "bthompso",
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
				Rotations: []models.OnCallRotation{},
				Enabled:   false,
				LastUpdated: now.Add(-time.Hour * 12),
				UpdatedBy:   "sec.admin",
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
			Type:        models.CrewTypeStandard,
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
				Rotations: []models.OnCallRotation{},
				Enabled:   false,
				LastUpdated: now.Add(-time.Hour * 48),
				UpdatedBy:   "dba.lead",
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
				Rotations: []models.OnCallRotation{},
				Enabled:   false,
				LastUpdated: now.Add(-time.Hour * 72),
				UpdatedBy:   "mobile.lead",
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
			Type:        models.CrewTypeVirtual,
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
				Rotations: []models.OnCallRotation{
					{
						Name:           "Primary",
						OnCallMember:   "bthompso",
						RotationStarts: now.AddDate(0, -1, 0), // Started 1 month ago
						RotationEnds:   now.AddDate(0, 1, 0),  // Ends in 1 month
					},
				},
				Enabled:     true,
				LastUpdated: now.Add(-time.Hour * 48),
				UpdatedBy:   "governance.lead",
			},
			OwnedAssets: []string{
				"asset://governance/ownership-matrix",
				"asset://governance/escalation-paths",
			},
		},
		{
			ID:          "CREW-2345",
			VanityName:  "Design Foundations",
			Type:        models.CrewTypeStandard,
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
				Rotations: []models.OnCallRotation{},
				Enabled:   false,
				LastUpdated: now.Add(-time.Hour * 168),
				UpdatedBy:   "design.lead",
			},
			OwnedAssets: []string{
				"asset://design/foundation-tokens",
				"asset://design/component-library",
			},
		},
		{
			ID:          "CREW-4567",
			VanityName:  "Design Systems",
			Type:        models.CrewTypeStandard,
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
				Rotations: []models.OnCallRotation{},
				Enabled:   false,
				LastUpdated: now.Add(-time.Hour * 96),
				UpdatedBy:   "design.systems.lead",
			},
			OwnedAssets: []string{
				"asset://design/advanced-components",
				"asset://design/automation-tools",
			},
		},
		{
			ID:          "CREW-5678",
			VanityName:  "Design Engineering",
			Type:        models.CrewTypeStandard,
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
					OnCallType: models.OnCallSecondary,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 2),
				},
			},
			OnCallSchedule: models.OnCallSchedule{
				Rotations: []models.OnCallRotation{
					{
						Name:           "Secondary",
						OnCallMember:   "bthompso",
						RotationStarts: now.AddDate(0, 0, -7), // Started 1 week ago
						RotationEnds:   now.AddDate(0, 0, 7),  // Ends in 1 week
					},
				},
				Enabled:     true,
				LastUpdated: now.Add(-time.Hour * 72),
				UpdatedBy:   "design.eng.lead",
			},
			OwnedAssets: []string{
				"asset://design/implementation-guides",
				"asset://design/code-generation-tools",
			},
		},
		{
			ID:          "CREW-7890",
			VanityName:  "Web Performance Team",
			Description: "Optimizes web application performance and user experience",
			Type:        models.CrewTypeStandard,
			CreatedAt:   now.AddDate(0, -4, 0),
			CreatedBy:   "perf.lead",
			Status:      models.CrewStatusActive,
			Members: []models.Member{
				{
					UserID:     "perf.lead",
					Email:      "perf.lead@company.com",
					FullName:   "Performance Lead",
					Level:      "IC7",
					Role:       models.RoleOwner,
					JoinedAt:   now.AddDate(0, -4, 0),
					AddedBy:    "system",
					IsOnCall:   false,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 3),
				},
				{
					UserID:     "webdev1",
					Email:      "webdev1@company.com",
					FullName:   "Sarah Web Developer",
					Level:      "IC4",
					Role:       models.RoleMember,
					JoinedAt:   now.AddDate(0, -3, 0),
					AddedBy:    "perf.lead",
					IsOnCall:   false,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 1),
				},
			},
			OnCallSchedule: models.OnCallSchedule{
				Rotations: []models.OnCallRotation{},
				Enabled:     false,
				LastUpdated: now.Add(-time.Hour * 96),
				UpdatedBy:   "perf.lead",
			},
			OwnedAssets: []string{
				"asset://web/performance-monitor",
				"asset://web/optimization-tools",
			},
		},
		{
			ID:          "CREW-8901",
			VanityName:  "Web Security Team",
			Description: "Ensures web application security and vulnerability management",
			Type:        models.CrewTypeStandard,
			CreatedAt:   now.AddDate(0, -3, 0),
			CreatedBy:   "websec.lead",
			Status:      models.CrewStatusActive,
			Members: []models.Member{
				{
					UserID:     "websec.lead",
					Email:      "websec.lead@company.com",
					FullName:   "Web Security Lead",
					Level:      "IC6",
					Role:       models.RoleOwner,
					JoinedAt:   now.AddDate(0, -3, 0),
					AddedBy:    "system",
					IsOnCall:   true,
					OnCallType: models.OnCallPrimary,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 2),
				},
				{
					UserID:     "websec1",
					Email:      "websec1@company.com",
					FullName:   "Alex Web Security",
					Level:      "IC5",
					Role:       models.RoleMember,
					JoinedAt:   now.AddDate(0, -2, 0),
					AddedBy:    "websec.lead",
					IsOnCall:   false,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 4),
				},
			},
			OnCallSchedule: models.OnCallSchedule{
				Rotations: []models.OnCallRotation{
					{
						Name:           "Primary",
						OnCallMember:   "websec.lead",
						RotationStarts: now.AddDate(0, 0, -5),
						RotationEnds:   now.AddDate(0, 0, 9),
					},
				},
				Enabled:     true,
				LastUpdated: now.Add(-time.Hour * 120),
				UpdatedBy:   "websec.lead",
			},
			OwnedAssets: []string{
				"asset://web/security-scanner",
				"asset://web/firewall-config",
			},
		},
		{
			ID:          "CREW-9012",
			VanityName:  "Web Analytics Team",
			Description: "Manages web analytics, tracking, and user insights",
			Type:        models.CrewTypeStandard,
			CreatedAt:   now.AddDate(0, -2, 0),
			CreatedBy:   "analytics.lead",
			Status:      models.CrewStatusActive,
			Members: []models.Member{
				{
					UserID:     "analytics.lead",
					Email:      "analytics.lead@company.com",
					FullName:   "Analytics Lead",
					Level:      "IC5",
					Role:       models.RoleOwner,
					JoinedAt:   now.AddDate(0, -2, 0),
					AddedBy:    "system",
					IsOnCall:   false,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 1),
				},
				{
					UserID:     "webanalyst1",
					Email:      "webanalyst1@company.com",
					FullName:   "Jordan Web Analyst",
					Level:      "IC3",
					Role:       models.RoleMember,
					JoinedAt:   now.AddDate(0, -1, 0),
					AddedBy:    "analytics.lead",
					IsOnCall:   false,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 5),
				},
			},
			OnCallSchedule: models.OnCallSchedule{
				Rotations: []models.OnCallRotation{},
				Enabled:     false,
				LastUpdated: now.Add(-time.Hour * 48),
				UpdatedBy:   "analytics.lead",
			},
			OwnedAssets: []string{
				"asset://web/analytics-dashboard",
				"asset://web/tracking-system",
			},
		},
		{
			ID:          "CREW-0123",
			VanityName:  "Web Infrastructure Team",
			Description: "Maintains web servers, CDN, and deployment infrastructure",
			Type:        models.CrewTypeStandard,
			CreatedAt:   now.AddDate(0, -5, 0),
			CreatedBy:   "webinfra.lead",
			Status:      models.CrewStatusActive,
			Members: []models.Member{
				{
					UserID:     "webinfra.lead",
					Email:      "webinfra.lead@company.com",
					FullName:   "Web Infrastructure Lead",
					Level:      "IC6",
					Role:       models.RoleOwner,
					JoinedAt:   now.AddDate(0, -5, 0),
					AddedBy:    "system",
					IsOnCall:   true,
					OnCallType: models.OnCallPrimary,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 1),
				},
				{
					UserID:     "webops1",
					Email:      "webops1@company.com",
					FullName:   "Chris Web Ops",
					Level:      "IC4",
					Role:       models.RoleMember,
					JoinedAt:   now.AddDate(0, -4, 0),
					AddedBy:    "webinfra.lead",
					IsOnCall:   true,
					OnCallType: models.OnCallSecondary,
					Status:     models.StatusActive,
					LastActive: now.Add(-time.Hour * 2),
				},
			},
			OnCallSchedule: models.OnCallSchedule{
				Rotations: []models.OnCallRotation{
					{
						Name:           "Primary",
						OnCallMember:   "webinfra.lead",
						RotationStarts: now.AddDate(0, 0, -7),
						RotationEnds:   now.AddDate(0, 0, 7),
					},
					{
						Name:           "Secondary",
						OnCallMember:   "webops1",
						RotationStarts: now.AddDate(0, 0, -10),
						RotationEnds:   now.AddDate(0, 0, 4),
					},
				},
				Enabled:     true,
				LastUpdated: now.Add(-time.Hour * 168),
				UpdatedBy:   "webinfra.lead",
			},
			OwnedAssets: []string{
				"asset://web/load-balancer",
				"asset://web/cdn-config",
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

	// Create catalog assets
	s.populateCatalogAssets(now)
}

// populateCatalogAssets creates the asset catalog simulation data
func (s *SimulationDataStore) populateCatalogAssets(now time.Time) {
	catalogAssets := []*models.CatalogAsset{
		{
			ID:           "AC123456",
			VanityName:   "EngX Web Application",
			AssetType:    models.AssetTypeWebApp,
			Description:  "Core engineering productivity web application with multi-product capabilities",
			OwningCrewID: "CREW-1234", // Web Platform Team
			CreatedAt:    now.AddDate(0, -8, 0),
			LastModified: now.Add(-time.Hour * 24),
			ModifiedBy:   "bthompso",
			Status:       models.AssetActive,
			Version:      "2.1.4",
			Environment:  "production",
			URNs:         []string{"asset://web-app/dashboard", "asset://web-app/engx"},
			Tags:         []string{"web", "productivity", "engineering", "dashboard"},

			AccessGrants: []models.AssetAccessGrant{
				{
					CrewID:      "CREW-1234", // Web Platform Team
					AccessLevel: models.AccessLevelOwner,
					GrantedAt:   now.AddDate(0, -8, 0),
					GrantedBy:   "system",
					Status:      models.GrantActive,
				},
				{
					CrewID:      "CREW-0014", // Universal Ownership
					AccessLevel: models.AccessLevelPublish,
					GrantedAt:   now.AddDate(0, -6, 0),
					GrantedBy:   "bthompso",
					Status:      models.GrantActive,
				},
				{
					CrewID:      "CREW-6789", // DPX Platform UX (placeholder)
					AccessLevel: models.AccessLevelReadOnly,
					GrantedAt:   now.AddDate(0, -4, 0),
					GrantedBy:   "bthompso",
					Status:      models.GrantActive,
				},
				{
					CrewID:      "CREW-2345", // Design Foundations
					AccessLevel: models.AccessLevelReadOnly,
					GrantedAt:   now.AddDate(0, -3, 0),
					GrantedBy:   "design.lead",
					Status:      models.GrantActive,
				},
				{
					CrewID:      "CREW-5678", // Design Engineering
					AccessLevel: models.AccessLevelReadOnly,
					GrantedAt:   now.AddDate(0, -2, 0),
					GrantedBy:   "design.eng.lead",
					Status:      models.GrantActive,
				},
				{
					CrewID:      "CREW-4567", // Design Systems
					AccessLevel: models.AccessLevelReadOnly,
					GrantedAt:   now.AddDate(0, -1, 0),
					GrantedBy:   "design.systems.lead",
					Status:      models.GrantActive,
				},
			},

			Dependencies: []models.AssetDependency{
				{
					DependencyID: "AC654321",
					Name:         "ShadCN UI Design System (SUDS)",
					Version:      "1.0.00",
					OwningCrewID: "CREW-123456", // placeholder
					DependencyType: models.DependencyRuntime,
					// No Health field in CatalogAsset
					LastChecked:  now.Add(-time.Hour * 2),
					Required:     true,
				},
				{
					DependencyID: "AC2F6534",
					Name:         "TrustBridge SSO for Web",
					Version:      "12.1.50",
					OwningCrewID: "CREW-4567",
					DependencyType: models.DependencyRuntime,
					// No Health field in CatalogAsset // GREEN - all good
					LastChecked:  now.Add(-time.Hour * 1),
					Required:     true,
				},
				{
					DependencyID: "AC2ACB61",
					Name:         "gRPC for Web Apps",
					Version:      "0.8.00",
					OwningCrewID: "CREW-1234", // Web Platform Team
					DependencyType: models.DependencyRuntime,
					// No Health field in CatalogAsset // GREEN - all good
					LastChecked:  now.Add(-time.Hour * 6),
					Required:     true,
				},
				{
					DependencyID: "AC9D27A6",
					Name:         "React Router 7",
					Version:      "7.1.50",
					OwningCrewID: "CREW-1234", // Web Platform Team
					DependencyType: models.DependencyRuntime,
					// No Health field in CatalogAsset // GREEN - all good
					LastChecked:  now.Add(-time.Minute * 15),
					Required:     true,
				},
				{
					DependencyID: "AC512F3A",
					Name:         "Apollo GraphQL",
					Version:      "0.1.10",
					OwningCrewID: "CREW-2345", // Design Foundations
					DependencyType: models.DependencyRuntime,
					// No Health field in CatalogAsset
					LastChecked:  now.Add(-time.Hour * 3),
					Required:     true,
				},
				{
					DependencyID: "ACF4A62B",
					Name:         "CREWS API",
					Version:      "3.0.00",
					OwningCrewID: "CREW-12", // placeholder
					DependencyType: models.DependencyRuntime,
					// No Health field in CatalogAsset // GREEN - all good
					LastChecked:  now.Add(-time.Hour * 24),
					Required:     true,
				},
				{
					DependencyID: "AC1B2941",
					Name:         "CATALOG API",
					Version:      "4.1.12",
					OwningCrewID: "CREW-12", // placeholder
					DependencyType: models.DependencyRuntime,
					// No Health field in CatalogAsset // GREEN - all good
					LastChecked:  now.Add(-time.Hour * 12),
					Required:     true,
				},
				{
					DependencyID: "AC724AB5",
					Name:         "Some Other Dependency Long Name",
					Version:      "15.6.45",
					OwningCrewID: "CREW-35897", // placeholder
					DependencyType: models.DependencyBuild,
					// No Health field in CatalogAsset
					LastChecked:  now.Add(-time.Hour * 4),
					Required:     false,
				},
				{
					DependencyID: "AC6842D1",
					Name:         "Yet Another Long Dependency Name",
					Version:      "4.1.20",
					OwningCrewID: "CREW-5129", // placeholder
					DependencyType: models.DependencyDev,
					Health:       models.HealthAttention,
					LastChecked:  now.Add(-time.Hour * 8),
					Required:     false,
				},
				{
					DependencyID: "AC02A048",
					Name:         "Super Simple Dependency Name",
					Version:      "2.0.60",
					OwningCrewID: "CREW-250", // placeholder
					DependencyType: models.DependencyTest,
					Health:       models.HealthELR, // MAGENTA - external library request
					LastChecked:  now.Add(-time.Hour * 1),
					Required:     false,
				},
			},
		},
		{
			ID:           "AC789012",
			VanityName:   "Web Performance Monitor",
			AssetType:    models.AssetTypeWebApp,
			Description:  "Real-time web application performance monitoring and alerting system",
			OwningCrewID: "CREW-7890", // Web Performance Team
			CreatedAt:    now.AddDate(0, -4, 0),
			LastModified: now.Add(-time.Hour * 12),
			ModifiedBy:   "perf.lead",
			Status:       models.AssetActive,
			Version:      "3.2.1",
			Environment:  "production",
			// No Health field in CatalogAsset
			URNs: []string{
				"asset://web/performance-monitor",
				"asset://monitoring/web-perf",
			},
			AccessGrants: []models.AssetAccessGrant{
				{
					CrewID:      "CREW-1234", // Web Platform Team
					AccessLevel: models.AccessLevelReadOnly,
					GrantedAt:   now.AddDate(0, -2, 0),
					GrantedBy:   "perf.lead",
					Status:      models.GrantActive,
				},
			},
			Dependencies: []models.AssetDependency{
				{
					DependencyID: "AC123456",
					Name:         "EngX Web Application",
					Version:      "2.1.4",
					OwningCrewID: "CREW-1234",
					DependencyType: models.DependencyRuntime,
					// No Health field in CatalogAsset
					LastChecked:  now.Add(-time.Hour * 2),
					Required:     true,
				},
			},
		},
		{
			ID:           "AC890123",
			VanityName:   "Web Security Scanner",
			AssetType:    models.AssetTypeService,
			Description:  "Automated web application security vulnerability scanner",
			OwningCrewID: "CREW-8901", // Web Security Team
			CreatedAt:    now.AddDate(0, -3, 0),
			LastModified: now.Add(-time.Hour * 6),
			ModifiedBy:   "websec.lead",
			Status:       models.AssetActive,
			Version:      "1.8.5",
			Environment:  "production",
			// No Health field in CatalogAsset
			URNs: []string{
				"asset://web/security-scanner",
				"asset://security/web-scanner",
			},
			AccessGrants: []models.AssetAccessGrant{
				{
					CrewID:      "CREW-2345", // Security Response Team
					AccessLevel: "read-write",
					GrantedAt:   now.AddDate(0, -2, 0),
					GrantedBy:   "websec.lead",
					Status:      models.GrantActive,
				},
			},
			Dependencies: []models.AssetDependency{
				{
					DependencyID: "AC123456",
					Name:         "EngX Web Application",
					Version:      "2.1.4",
					OwningCrewID: "CREW-1234",
					DependencyType: models.DependencyRuntime,
					// No Health field in CatalogAsset
					LastChecked:  now.Add(-time.Hour * 1),
					Required:     true,
				},
			},
		},
		{
			ID:           "AC901234",
			VanityName:   "Web Analytics Dashboard",
			AssetType:    models.AssetTypeWebApp,
			Description:  "Comprehensive web analytics and user behavior tracking dashboard",
			OwningCrewID: "CREW-9012", // Web Analytics Team
			CreatedAt:    now.AddDate(0, -2, 0),
			LastModified: now.Add(-time.Hour * 8),
			ModifiedBy:   "analytics.lead",
			Status:       models.AssetActive,
			Version:      "4.0.2",
			Environment:  "production",
			// No Health field in CatalogAsset
			URNs: []string{
				"asset://web/analytics-dashboard",
				"asset://analytics/web-dashboard",
			},
			AccessGrants: []models.AssetAccessGrant{
				{
					CrewID:      "CREW-1234", // Web Platform Team
					AccessLevel: models.AccessLevelReadOnly,
					GrantedAt:   now.AddDate(0, -1, 0),
					GrantedBy:   "analytics.lead",
					Status:      models.GrantActive,
				},
			},
			Dependencies: []models.AssetDependency{
				{
					DependencyID: "AC123456",
					Name:         "EngX Web Application",
					Version:      "2.1.4",
					OwningCrewID: "CREW-1234",
					DependencyType: "integration",
					// No Health field in CatalogAsset
					LastChecked:  now.Add(-time.Hour * 3),
					Required:     false,
				},
			},
		},
		{
			ID:           "AC012345",
			VanityName:   "Web Load Balancer",
			AssetType:    "Infrastructure",
			Description:  "High-availability web traffic load balancer and routing system",
			OwningCrewID: "CREW-0123", // Web Infrastructure Team
			CreatedAt:    now.AddDate(0, -5, 0),
			LastModified: now.Add(-time.Hour * 4),
			ModifiedBy:   "webinfra.lead",
			Status:       models.AssetActive,
			Version:      "2.7.3",
			Environment:  "production",
			// No Health field in CatalogAsset
			URNs: []string{
				"asset://web/load-balancer",
				"asset://infrastructure/web-lb",
			},
			AccessGrants: []models.AssetAccessGrant{
				{
					CrewID:      "CREW-1234", // Web Platform Team
					AccessLevel: models.AccessLevelReadOnly,
					GrantedAt:   now.AddDate(0, -3, 0),
					GrantedBy:   "webinfra.lead",
					Status:      models.GrantActive,
				},
				{
					CrewID:      "CREW-7890", // Web Performance Team
					AccessLevel: models.AccessLevelReadOnly,
					GrantedAt:   now.AddDate(0, -1, 0),
					GrantedBy:   "webinfra.lead",
					Status:      models.GrantActive,
				},
			},
			Dependencies: []models.AssetDependency{
				{
					DependencyID: "AC123456",
					Name:         "EngX Web Application",
					Version:      "2.1.4",
					OwningCrewID: "CREW-1234",
					DependencyType: models.DependencyRuntime,
					// No Health field in CatalogAsset
					LastChecked:  now.Add(-time.Hour * 1),
					Required:     true,
				},
			},
		},
	}

	// Store catalog assets and build indices
	for _, asset := range catalogAssets {
		s.catalog[asset.ID] = asset

		// Build URN index
		for _, urn := range asset.URNs {
			s.catalogURNs[urn] = asset.ID
		}

		// Build name index (both full and partial matches)
		s.catalogNames[strings.ToLower(asset.VanityName)] = asset.ID
		// Add partial name matches
		words := strings.Fields(strings.ToLower(asset.VanityName))
		for _, word := range words {
			if len(word) > 2 { // Only index words longer than 2 characters
				s.catalogNames[word] = asset.ID
			}
		}
		// Also add "engx" as a shortcut for the main asset
		if strings.Contains(strings.ToLower(asset.VanityName), "engx") {
			s.catalogNames["engx"] = asset.ID
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

// GetCatalogAsset retrieves a catalog asset by ID
func (s *SimulationDataStore) GetCatalogAsset(id string) (*models.CatalogAsset, error) {
	asset, exists := s.catalog[id]
	if !exists {
		return nil, fmt.Errorf("catalog asset %s not found", id)
	}
	return asset, nil
}

// GetCatalogAssetByURN retrieves a catalog asset by URN
func (s *SimulationDataStore) GetCatalogAssetByURN(urn string) (*models.CatalogAsset, error) {
	assetID, exists := s.catalogURNs[urn]
	if !exists {
		return nil, fmt.Errorf("catalog asset with URN %s not found", urn)
	}
	return s.GetCatalogAsset(assetID)
}

// GetCatalogAssetByName retrieves a catalog asset by vanity name (or partial name)
func (s *SimulationDataStore) GetCatalogAssetByName(name string) (*models.CatalogAsset, error) {
	assetID, exists := s.catalogNames[strings.ToLower(name)]
	if !exists {
		return nil, fmt.Errorf("catalog asset with name %s not found", name)
	}
	return s.GetCatalogAsset(assetID)
}

// ResolveCatalogAssetParameter handles multiple input formats for catalog assets
func (s *SimulationDataStore) ResolveCatalogAssetParameter(param string) (*models.CatalogAsset, error) {
	// Try direct ID match first (AC123456)
	if strings.HasPrefix(strings.ToUpper(param), "AC") {
		return s.GetCatalogAsset(strings.ToUpper(param))
	}

	// Try URN match (asset://web-app/dashboard)
	if strings.HasPrefix(param, "asset://") {
		return s.GetCatalogAssetByURN(param)
	}

	// Try name match (exact or partial)
	return s.GetCatalogAssetByName(param)
}

// Helper function to normalize user identifiers
func (s *SimulationDataStore) normalizeUserID(userID string) string {
	// Handle email addresses by extracting username
	if strings.Contains(userID, "@") {
		return strings.Split(userID, "@")[0]
	}
	return userID
}

// SearchAssets searches catalog assets by name, description, or ID
func (s *SimulationDataStore) SearchAssets(query string) ([]*models.CatalogAsset, error) {
	var results []*models.CatalogAsset
	query = strings.ToLower(query)

	for _, asset := range s.catalog {
		if strings.Contains(strings.ToLower(asset.VanityName), query) ||
			strings.Contains(strings.ToLower(asset.Description), query) ||
			strings.Contains(strings.ToLower(asset.ID), query) ||
			strings.Contains(strings.ToLower(string(asset.AssetType)), query) {
			results = append(results, asset)
		}
	}

	return results, nil
}

// SearchUsers searches for users by name, LDAP ID, or email across all crew memberships
func (s *SimulationDataStore) SearchUsers(query string) ([]models.Member, error) {
	var results []models.Member
	seen := make(map[string]bool) // Prevent duplicates
	query = strings.ToLower(query)

	for _, crew := range s.crews {
		for _, member := range crew.Members {
			// Create unique key to prevent duplicates
			memberKey := member.UserID
			if seen[memberKey] {
				continue
			}

			if strings.Contains(strings.ToLower(member.FullName), query) ||
				strings.Contains(strings.ToLower(member.UserID), query) ||
				strings.Contains(strings.ToLower(member.Email), query) ||
				strings.Contains(strings.ToLower(member.Level), query) {
				results = append(results, member)
				seen[memberKey] = true
			}
		}
	}

	return results, nil
}