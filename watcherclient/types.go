package watcherclient

import "time"

// Audit represents a Watcher audit
type Audit struct {
	UUID        string                   `json:"uuid,omitempty"`
	Name        string                   `json:"name,omitempty"`
	AuditType   string                   `json:"audit_type"` // ONESHOT, CONTINUOUS
	State       string                   `json:"state,omitempty"`
	Goal        string                   `json:"goal"`
	Strategy    string                   `json:"strategy,omitempty"`
	Interval    int                      `json:"interval,omitempty"`
	Scope       []map[string]interface{} `json:"scope,omitempty"`
	Parameters  map[string]interface{}   `json:"parameters,omitempty"`
	AutoTrigger bool                     `json:"auto_trigger"`
	NextRunTime *time.Time               `json:"next_run_time,omitempty"`
	Hostname    string                   `json:"hostname,omitempty"`
	CreatedAt   *time.Time               `json:"created_at,omitempty"`
	UpdatedAt   *time.Time               `json:"updated_at,omitempty"`
	DeletedAt   *time.Time               `json:"deleted_at,omitempty"`
	Links       []Link                   `json:"links,omitempty"`
}

// AuditTemplate represents a Watcher audit template
type AuditTemplate struct {
	UUID        string                   `json:"uuid,omitempty"`
	Name        string                   `json:"name"`
	Description string                   `json:"description,omitempty"`
	Goal        string                   `json:"goal"`
	Strategy    string                   `json:"strategy,omitempty"`
	Scope       []map[string]interface{} `json:"scope,omitempty"`
	CreatedAt   *time.Time               `json:"created_at,omitempty"`
	UpdatedAt   *time.Time               `json:"updated_at,omitempty"`
	DeletedAt   *time.Time               `json:"deleted_at,omitempty"`
	Links       []Link                   `json:"links,omitempty"`
}

// ActionPlan represents a Watcher action plan
type ActionPlan struct {
	UUID           string                 `json:"uuid,omitempty"`
	AuditUUID      string                 `json:"audit_uuid,omitempty"`
	State          string                 `json:"state,omitempty"`
	Strategy       string                 `json:"strategy,omitempty"`
	GlobalEfficacy map[string]interface{} `json:"global_efficacy,omitempty"`
	Hostname       string                 `json:"hostname,omitempty"`
	CreatedAt      *time.Time             `json:"created_at,omitempty"`
	UpdatedAt      *time.Time             `json:"updated_at,omitempty"`
	DeletedAt      *time.Time             `json:"deleted_at,omitempty"`
	Links          []Link                 `json:"links,omitempty"`
}

// Action represents a Watcher action
type Action struct {
	UUID           string                 `json:"uuid,omitempty"`
	ActionPlanUUID string                 `json:"action_plan_uuid,omitempty"`
	ActionType     string                 `json:"action_type,omitempty"`
	State          string                 `json:"state,omitempty"`
	Parameters     map[string]interface{} `json:"parameters,omitempty"`
	ParentsUUIDs   []string               `json:"parents,omitempty"`
	CreatedAt      *time.Time             `json:"created_at,omitempty"`
	UpdatedAt      *time.Time             `json:"updated_at,omitempty"`
	DeletedAt      *time.Time             `json:"deleted_at,omitempty"`
	Links          []Link                 `json:"links,omitempty"`
}

type EfficacyIndicatorSpec struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Unit        *string `json:"unit"` // nullable
	Schema      string  `json:"schema"`
}

// Goal represents an optimization goal
type Goal struct {
	UUID        string                  `json:"uuid,omitempty"`
	Name        string                  `json:"name"`
	DisplayName string                  `json:"display_name,omitempty"`
	Efficacy    []EfficacyIndicatorSpec `json:"efficacy_specification,omitempty"`
	CreatedAt   *time.Time              `json:"created_at,omitempty"`
	UpdatedAt   *time.Time              `json:"updated_at,omitempty"`
	DeletedAt   *time.Time              `json:"deleted_at,omitempty"`
	Links       []Link                  `json:"links,omitempty"`
}

// Strategy represents an optimization strategy
type Strategy struct {
	UUID        string              `json:"uuid,omitempty"`
	Name        string              `json:"name"`
	DisplayName string              `json:"display_name,omitempty"`
	GoalUUID    string              `json:"goal_uuid,omitempty"`
	Parameters  []StrategyParameter `json:"parameters_spec,omitempty"`
	CreatedAt   *time.Time          `json:"created_at,omitempty"`
	UpdatedAt   *time.Time          `json:"updated_at,omitempty"`
	DeletedAt   *time.Time          `json:"deleted_at,omitempty"`
	Links       []Link              `json:"links,omitempty"`
}

// StrategyParameter represents a strategy parameter
type StrategyParameter struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Default     interface{} `json:"default,omitempty"`
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required"`
}

// DataModel represents the infrastructure data model
type DataModel struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// Link represents a HATEOAS link
type Link struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
}

// ListOptions represents common list options
type ListOptions struct {
	Limit   int    `json:"limit,omitempty"`
	Marker  string `json:"marker,omitempty"`
	SortKey string `json:"sort_key,omitempty"`
	SortDir string `json:"sort_dir,omitempty"`
}

// Response wrapper types
type AuditsResponse struct {
	Audits []Audit `json:"audits"`
}

type AuditTemplatesResponse struct {
	AuditTemplates []AuditTemplate `json:"audit_templates"`
}

type ActionPlansResponse struct {
	ActionPlans []ActionPlan `json:"action_plans"`
}

type ActionsResponse struct {
	Actions []Action `json:"actions"`
}

type GoalsResponse struct {
	Goals []Goal `json:"goals"`
}

type StrategiesResponse struct {
	Strategies []Strategy `json:"strategies"`
}
