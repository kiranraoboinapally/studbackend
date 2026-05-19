package domain

import "time"

// ─── Identity & Access Management ────────────────────────────────────────────

type User struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Username     string     `gorm:"unique;not null;index" json:"username"`
	Email        string     `gorm:"unique;not null;index" json:"email"`
	PasswordHash string     `gorm:"not null" json:"-"`
	IsActive     bool       `gorm:"default:true;index" json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (User) TableName() string { return "shared.users" }

type Role struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	RoleName    string `gorm:"unique;not null" json:"role_name"`
	Description string `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Role) TableName() string { return "shared.roles" }

type UserRole struct {
	UserID     uint      `gorm:"primaryKey" json:"user_id"`
	RoleID     uint      `gorm:"primaryKey" json:"role_id"`
	AssignedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"assigned_at"`
	AssignedBy *uint     `json:"assigned_by,omitempty"`
}

func (UserRole) TableName() string { return "shared.user_roles" }

type AuditLog struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        *uint     `gorm:"index" json:"user_id,omitempty"`
	Action        string    `gorm:"type:varchar(50);not null;index" json:"action"`
	SchemaName    string    `gorm:"type:varchar(50)" json:"schema_name"`
	AffectedTable string   `gorm:"type:varchar(100)" json:"affected_table"`
	RecordID      string    `gorm:"index" json:"record_id"`
	OldValue      string    `gorm:"type:jsonb" json:"old_value,omitempty"`
	NewValue      string    `gorm:"type:jsonb" json:"new_value,omitempty"`
	IPAddress     string    `json:"ip_address"`
	UserAgent     string    `json:"user_agent"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP;index" json:"created_at"`
}

func (AuditLog) TableName() string { return "shared.audit_logs" }

// ─── Outbox Event (Transactional Outbox Pattern) ─────────────────────────────

type OutboxEvent struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	AggregateType string    `gorm:"type:varchar(100);not null;index" json:"aggregate_type"`
	AggregateID   string    `gorm:"type:varchar(100);not null;index" json:"aggregate_id"`
	EventType     string    `gorm:"type:varchar(100);not null;index" json:"event_type"`
	Payload       string    `gorm:"type:jsonb;not null" json:"payload"`
	Published     bool      `gorm:"default:false;index" json:"published"`
	PublishedAt   *time.Time `json:"published_at,omitempty"`
	RetryCount    int       `gorm:"default:0" json:"retry_count"`
	LastError     string    `json:"last_error,omitempty"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP;index" json:"created_at"`
}

func (OutboxEvent) TableName() string { return "shared.outbox_events" }
