package models

import (
	"time"
)

// Crew represents a team or group of people with shared responsibilities
type Crew struct {
	ID              string              `json:"id"`              // CREW-1234
	VanityName      string              `json:"vanity_name"`     // "Web Platform Team"
	Description     string              `json:"description"`     // Team purpose
	Type            CrewType            `json:"type"`            // STANDARD, VIRTUAL, etc.
	CreatedAt       time.Time           `json:"created_at"`
	CreatedBy       string              `json:"created_by"`      // Creator LDAP
	Members         []Member            `json:"members"`
	OnCallSchedule  OnCallSchedule      `json:"oncall_schedule"`
	OwnedAssets     []string            `json:"owned_assets"`    // Asset URNs
	AccessGrants    []AccessGrant       `json:"access_grants"`   // External permissions
	Status          CrewStatus          `json:"status"`          // active, archived, etc.
	LastModified    time.Time           `json:"last_modified"`
	ModifiedBy      string              `json:"modified_by"`
}

// Member represents a user's membership in a crew with role-based access
type Member struct {
	UserID          string              `json:"user_id"`         // LDAP username
	Email           string              `json:"email"`           // Full email address
	FullName        string              `json:"full_name"`       // Display name
	Level           string              `json:"level"`           // IC5, MR3, IC3, etc.
	Role            MemberRole          `json:"role"`            // Access level
	JoinedAt        time.Time           `json:"joined_at"`
	AddedBy         string              `json:"added_by"`        // Who added them
	IsOnCall        bool                `json:"is_oncall"`       // Current on-call status
	OnCallType      OnCallType          `json:"oncall_type"`     // Primary, backup, temp
	OnCallSchedule  []OnCallPeriod      `json:"oncall_periods"`  // Scheduled rotations
	LastActive      time.Time           `json:"last_active"`     // Activity tracking
	Status          MemberStatus        `json:"status"`          // active, inactive, etc.
}

// MemberRole defines the access level within a crew
type MemberRole string

const (
	RoleOwner   MemberRole = "owner"   // Full control, can transfer ownership
	RoleAdmin   MemberRole = "admin"   // Management privileges
	RoleMember  MemberRole = "member"  // Standard access
	RoleTemp    MemberRole = "temp"    // Temporary access with expiration
	RoleAuto    MemberRole = "auto"    // Inherited from asset ownership
	RoleRemoved MemberRole = "removed" // Explicitly revoked
)

// MemberStatus defines the current status of a member
type MemberStatus string

const (
	StatusActive   MemberStatus = "active"
	StatusInactive MemberStatus = "inactive"
	StatusPending  MemberStatus = "pending"
	StatusRemoved  MemberStatus = "removed"
)

// CrewStatus defines the current status of a crew
type CrewStatus string

const (
	CrewStatusActive   CrewStatus = "active"
	CrewStatusArchived CrewStatus = "archived"
	CrewStatusPending  CrewStatus = "pending"
	CrewStatusDisabled CrewStatus = "disabled"
)

// OnCallSchedule manages on-call rotation for a crew
type OnCallSchedule struct {
	Rotations       []OnCallRotation    `json:"rotations"`       // Array of active rotations
	Enabled         bool                `json:"enabled"`         // Whether on-call is active
	LastUpdated     time.Time           `json:"last_updated"`
	UpdatedBy       string              `json:"updated_by"`
}

// OnCallRotation represents a time-based on-call assignment for a crew member
type OnCallRotation struct {
	Name            string              `json:"name"`            // "Primary", "Secondary", "Manager"
	OnCallMember    string              `json:"oncall_member"`   // LDAP username
	RotationStarts  time.Time           `json:"rotation_starts"` // Start date/time
	RotationEnds    time.Time           `json:"rotation_ends"`   // End date/time
}

// OnCallPeriod represents a specific on-call assignment
type OnCallPeriod struct {
	UserID          string              `json:"user_id"`
	StartTime       time.Time           `json:"start_time"`
	EndTime         time.Time           `json:"end_time"`
	Type            OnCallType          `json:"type"`            // primary, backup
	Notes           string              `json:"notes"`           // Optional notes
}

// OnCallType defines the type of on-call assignment
type OnCallType string

const (
	OnCallPrimary   OnCallType = "primary"
	OnCallSecondary OnCallType = "secondary"
	OnCallTemp      OnCallType = "temp"
)

// AccessGrant represents external access delegation
type AccessGrant struct {
	GrantedTo       string              `json:"granted_to"`      // User or crew
	GrantedBy       string              `json:"granted_by"`      // Granting user
	GrantedAt       time.Time           `json:"granted_at"`
	AccessLevel     AccessLevel         `json:"access_level"`    // read, write, admin
	AssetScope      []string            `json:"asset_scope"`     // Specific assets
	ExpiresAt       *time.Time          `json:"expires_at"`      // Optional expiration
	Reason          string              `json:"reason"`          // Grant justification
	Status          GrantStatus         `json:"status"`          // active, expired, revoked
}

// AccessLevel defines the level of access granted
type AccessLevel string

const (
	AccessRead  AccessLevel = "read"
	AccessWrite AccessLevel = "write"
	AccessAdmin AccessLevel = "admin"
)

// GrantStatus defines the current status of an access grant
type GrantStatus string

const (
	GrantActive  GrantStatus = "active"
	GrantExpired GrantStatus = "expired"
	GrantRevoked GrantStatus = "revoked"
	GrantPending GrantStatus = "pending"
)

// AssetOwnership represents the relationship between crews and assets
type AssetOwnership struct {
	AssetURN        string              `json:"asset_urn"`       // Unique asset identifier
	AssetName       string              `json:"asset_name"`      // Human-readable name
	AssetType       string              `json:"asset_type"`      // web-app, service, etc.
	OwnerCrewID     string              `json:"owner_crew_id"`   // Primary owner
	Delegates       []string            `json:"delegates"`       // Delegated crews
	CreatedAt       time.Time           `json:"created_at"`
	LastModified    time.Time           `json:"last_modified"`
	ModifiedBy      string              `json:"modified_by"`
	Metadata        map[string]string   `json:"metadata"`        // Additional properties
	Status          AssetStatus         `json:"status"`          // active, deprecated, etc.
}

// AssetStatus defines the current status of an asset
type AssetStatus string

const (
	AssetActive     AssetStatus = "active"
	AssetDeprecated AssetStatus = "deprecated"
	AssetArchived   AssetStatus = "archived"
	AssetMaintenance AssetStatus = "maintenance"
)

// Helper methods for role checking
func (r MemberRole) IsAdmin() bool {
	return r == RoleOwner || r == RoleAdmin
}

func (r MemberRole) CanManage() bool {
	return r == RoleOwner || r == RoleAdmin
}

func (r MemberRole) IsOwner() bool {
	return r == RoleOwner
}

// Helper methods for member status
func (m *Member) IsActive() bool {
	return m.Status == StatusActive
}

func (m *Member) CanManageCrew() bool {
	return m.IsActive() && m.Role.CanManage()
}

func (m *Member) IsCurrentlyOnCall() bool {
	return m.IsActive() && m.IsOnCall
}

// Helper methods for crew
func (c *Crew) GetActiveMembers() []Member {
	var active []Member
	for _, member := range c.Members {
		if member.IsActive() {
			active = append(active, member)
		}
	}
	return active
}

func (c *Crew) GetAdmins() []Member {
	var admins []Member
	for _, member := range c.Members {
		if member.IsActive() && member.Role.IsAdmin() {
			admins = append(admins, member)
		}
	}
	return admins
}

func (c *Crew) GetOnCallMembers() []Member {
	var onCall []Member
	for _, member := range c.Members {
		if member.IsCurrentlyOnCall() {
			onCall = append(onCall, member)
		}
	}
	return onCall
}

func (c *Crew) GetOwner() *Member {
	for _, member := range c.Members {
		if member.IsActive() && member.Role.IsOwner() {
			return &member
		}
	}
	return nil
}

func (c *Crew) HasMember(userID string) bool {
	for _, member := range c.Members {
		if member.UserID == userID && member.IsActive() {
			return true
		}
	}
	return false
}

func (c *Crew) GetMember(userID string) *Member {
	for _, member := range c.Members {
		if member.UserID == userID {
			return &member
		}
	}
	return nil
}

func (c *Crew) IsActive() bool {
	return c.Status == CrewStatusActive
}

// GetCurrentOnCallRotations returns currently active on-call rotations
func (c *Crew) GetCurrentOnCallRotations() []OnCallRotation {
	if !c.OnCallSchedule.Enabled {
		return []OnCallRotation{}
	}

	now := time.Now()
	var currentRotations []OnCallRotation

	for _, rotation := range c.OnCallSchedule.Rotations {
		if now.After(rotation.RotationStarts) && now.Before(rotation.RotationEnds) {
			currentRotations = append(currentRotations, rotation)
		}
	}

	return currentRotations
}

// IsOnCallNow checks if a specific user is currently on-call
func (c *Crew) IsOnCallNow(userID string) bool {
	currentRotations := c.GetCurrentOnCallRotations()
	for _, rotation := range currentRotations {
		if rotation.OnCallMember == userID {
			return true
		}
	}
	return false
}

// GetOnCallRole returns the current on-call role name for a user (e.g., "Primary", "Secondary")
func (c *Crew) GetOnCallRole(userID string) string {
	currentRotations := c.GetCurrentOnCallRotations()
	for _, rotation := range currentRotations {
		if rotation.OnCallMember == userID {
			return rotation.Name
		}
	}
	return ""
}

// CrewType defines different types of crews
type CrewType string

const (
	CrewTypeStandard CrewType = "STANDARD"
	CrewTypeVirtual  CrewType = "VIRTUAL"
)