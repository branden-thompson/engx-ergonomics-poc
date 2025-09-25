package models

import (
	"time"
)

// CatalogAsset represents a complete asset in the asset catalog
type CatalogAsset struct {
	ID             string              `json:"id"`              // AC123456
	VanityName     string              `json:"vanity_name"`     // "EngX Web Application"
	AssetType      AssetType           `json:"asset_type"`      // Web Application, CLI Tool, etc.
	Description    string              `json:"description"`     // Asset description
	OwningCrewID   string              `json:"owning_crew_id"`  // CREW-1234
	CreatedAt      time.Time           `json:"created_at"`
	LastModified   time.Time           `json:"last_modified"`
	ModifiedBy     string              `json:"modified_by"`
	Status         AssetStatus         `json:"status"`          // active, deprecated, etc.

	// Access control
	AccessGrants   []AssetAccessGrant  `json:"access_grants"`   // Crew access permissions

	// Dependencies
	Dependencies   []AssetDependency   `json:"dependencies"`    // What this asset depends on
	DependentAssets []string           `json:"dependent_assets"` // Asset IDs that depend on this asset

	// Metadata
	Version        string              `json:"version"`         // Current version
	Environment    string              `json:"environment"`     // production, staging, etc.
	URNs           []string            `json:"urns"`            // asset://web-app/dashboard, etc.
	Tags           []string            `json:"tags"`            // searchable tags
	Metadata       map[string]string   `json:"metadata"`        // Additional key-value pairs
}

// AssetType defines the category of asset
type AssetType string

const (
	AssetTypeWebApp       AssetType = "Web Application (Multi-Product)"
	AssetTypeCLI          AssetType = "CLI Tool (Standalone)"
	AssetTypeRepository   AssetType = "Repository"
	AssetTypeService      AssetType = "Service"
	AssetTypeKafkaTopic   AssetType = "Kafka Topic"
	AssetTypePEMWorkflow  AssetType = "PEM Workflow"
	AssetType3rdParty     AssetType = "3rd Party Product"
	AssetTypeLibrary      AssetType = "Library"
	AssetTypePackage      AssetType = "Package"
	AssetTypeDataSource   AssetType = "Data Source"
)

// AssetAccessGrant represents crew access to an asset
type AssetAccessGrant struct {
	CrewID       string              `json:"crew_id"`         // CREW-1234
	AccessLevel  AssetAccessLevel    `json:"access_level"`    // Admin/Owner, can-publish, read-only
	GrantedAt    time.Time           `json:"granted_at"`
	GrantedBy    string              `json:"granted_by"`      // User who granted access
	ExpiresAt    *time.Time          `json:"expires_at"`      // Optional expiration
	Status       GrantStatus         `json:"status"`          // active, expired, revoked
}

// AssetAccessLevel defines the level of access to an asset
type AssetAccessLevel string

const (
	AccessLevelOwner     AssetAccessLevel = "Admin/Owner"
	AccessLevelPublish   AssetAccessLevel = "can-publish"
	AccessLevelReadOnly  AssetAccessLevel = "read-only"
	AccessLevelDeploy    AssetAccessLevel = "can-deploy"
	AccessLevelContrib   AssetAccessLevel = "contributor"
)

// AssetDependency represents a dependency relationship
type AssetDependency struct {
	DependencyID   string              `json:"dependency_id"`   // AC123456
	Name           string              `json:"name"`            // "ShadCN UI Design System (SUDS)"
	Version        string              `json:"version"`         // "1.0.00"
	OwningCrewID   string              `json:"owning_crew_id"`  // CREW-123456
	DependencyType DependencyType      `json:"dependency_type"` // runtime, build, dev
	Health         AssetHealthStatus   `json:"health"`          // healthy, attention, issue, etc.
	LastChecked    time.Time           `json:"last_checked"`
	Required       bool                `json:"required"`        // Is this a required dependency?
}

// DependencyType defines the type of dependency relationship
type DependencyType string

const (
	DependencyRuntime DependencyType = "runtime"
	DependencyBuild   DependencyType = "build"
	DependencyDev     DependencyType = "development"
	DependencyTest    DependencyType = "test"
)

// AssetHealthStatus represents the current health/status of an asset
type AssetHealthStatus string

const (
	HealthHealthy    AssetHealthStatus = "healthy"     // GREEN - All good
	HealthAttention  AssetHealthStatus = "attention"   // YELLOW - Some warnings
	HealthIssue      AssetHealthStatus = "issue"       // RED - Known issues/failures
	HealthDeploying  AssetHealthStatus = "deploying"   // BLUE - Currently deploying/running
	HealthUnknown    AssetHealthStatus = "unknown"     // GRAY - No information available
	HealthELR        AssetHealthStatus = "elr"         // MAGENTA - External Library Request
)

// OnCallMember represents someone who is currently on-call for an asset
type OnCallMember struct {
	UserID    string        `json:"user_id"`    // hbacot
	FullName  string        `json:"full_name"`  // "Hunter Bacot"
	Level     string        `json:"level"`      // IC3
	Priority  OnCallType    `json:"priority"`   // Primary, Secondary
	CrewID    string        `json:"crew_id"`    // CREW-1234
}

// Helper methods for asset management
func (a *CatalogAsset) GetOwningCrew() string {
	return a.OwningCrewID
}

func (a *CatalogAsset) HasAccess(crewID string) bool {
	for _, grant := range a.AccessGrants {
		if grant.CrewID == crewID && grant.Status == GrantActive {
			return true
		}
	}
	return false
}

func (a *CatalogAsset) GetAccessLevel(crewID string) AssetAccessLevel {
	for _, grant := range a.AccessGrants {
		if grant.CrewID == crewID && grant.Status == GrantActive {
			return grant.AccessLevel
		}
	}
	return ""
}

func (a *CatalogAsset) IsHealthy() bool {
	return a.Status == AssetActive
}

func (a *CatalogAsset) GetHealthyDependencies() []AssetDependency {
	var healthy []AssetDependency
	for _, dep := range a.Dependencies {
		if dep.Health == HealthHealthy {
			healthy = append(healthy, dep)
		}
	}
	return healthy
}

func (a *CatalogAsset) GetUnhealthyDependencies() []AssetDependency {
	var unhealthy []AssetDependency
	for _, dep := range a.Dependencies {
		if dep.Health != HealthHealthy {
			unhealthy = append(unhealthy, dep)
		}
	}
	return unhealthy
}

func (h AssetHealthStatus) GetColor() string {
	switch h {
	case HealthHealthy:
		return "92"   // ANSI foreground green
	case HealthAttention:
		return "226"  // Pure Yellow (#FFFF00)
	case HealthIssue:
		return "196"  // Pure Red (#FF0000)
	case HealthDeploying:
		return "94"   // ANSI foreground blue
	case HealthUnknown:
		return "251"  // Pure Grey (#C6C6C6)
	case HealthELR:
		return "201"  // Pure Magenta (#FF00FF)
	default:
		return "251"  // Pure Grey fallback
	}
}