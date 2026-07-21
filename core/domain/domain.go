// Package domain holds the core domain types and interfaces shared across the
// DevOS kernel and tenant services. It is intentionally free of business logic
// (see Sprint 0 Component 2 scope): it defines bounded-context identifiers
// and the minimal Aggregate contract used by event and persistence layers.
package domain

// Bounded-context identifiers are strong string types to prevent accidental
// confusion (e.g., an OrgID passed where a ProjectID is expected).
type (
	// OrgID identifies a tenant organization.
	OrgID string
	// ProjectID identifies a project within an organization.
	ProjectID string
	// IntentID identifies a user intent.
	IntentID string
	// PlanID identifies a generated plan (DAG).
	PlanID string
	// TaskID identifies a scheduled task.
	TaskID string
	// AgentID identifies a registered agent.
	AgentID string
	// WorkspaceID identifies an isolated workspace.
	WorkspaceID string
)

// IsValid reports whether the identifier is non-empty.
func (id OrgID) IsValid() bool { return id != "" }

// String returns the underlying value.
func (id OrgID) String() string { return string(id) }

// IsValid reports whether the identifier is non-empty.
func (id ProjectID) IsValid() bool { return id != "" }

// String returns the underlying value.
func (id ProjectID) String() string { return string(id) }

// IsValid reports whether the identifier is non-empty.
func (id IntentID) IsValid() bool { return id != "" }

// String returns the underlying value.
func (id IntentID) String() string { return string(id) }

// IsValid reports whether the identifier is non-empty.
func (id PlanID) IsValid() bool { return id != "" }

// String returns the underlying value.
func (id PlanID) String() string { return string(id) }

// IsValid reports whether the identifier is non-empty.
func (id TaskID) IsValid() bool { return id != "" }

// String returns the underlying value.
func (id TaskID) String() string { return string(id) }

// IsValid reports whether the identifier is non-empty.
func (id AgentID) IsValid() bool { return id != "" }

// String returns the underlying value.
func (id AgentID) String() string { return string(id) }

// IsValid reports whether the identifier is non-empty.
func (id WorkspaceID) IsValid() bool { return id != "" }

// String returns the underlying value.
func (id WorkspaceID) String() string { return string(id) }

// Aggregate is the minimal contract for a domain aggregate root.
type Aggregate interface {
	AggregateID() string
	AggregateType() string
}
