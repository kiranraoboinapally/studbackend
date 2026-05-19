package domain

import "time"

// ─── System Lookups & Reference Data ─────────────────────────────────────────

type Gender struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Code string `gorm:"unique;not null" json:"code"`
	Name string `gorm:"not null" json:"name"`
}

func (Gender) TableName() string { return "system.genders" }

type Category struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Code string `gorm:"unique;not null" json:"code"`
	Name string `gorm:"not null" json:"name"`
}

func (Category) TableName() string { return "system.categories" }

type BloodGroup struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Code string `gorm:"unique;not null" json:"code"`
	Name string `gorm:"not null" json:"name"`
}

func (BloodGroup) TableName() string { return "system.blood_groups" }

type StatusCode struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Module   string `gorm:"not null;index" json:"module"`
	Code     string `gorm:"not null;index" json:"code"`
	Name     string `gorm:"not null" json:"name"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
}

func (StatusCode) TableName() string { return "system.status_codes" }

// ─── System Configuration & Notifications ────────────────────────────────────

type Configuration struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	ConfigKey   string `gorm:"unique;not null;index" json:"config_key"`
	ConfigValue string `json:"config_value"`
	DataType    string `gorm:"default:'string'" json:"data_type"`
	Description string `json:"description"`
}

func (Configuration) TableName() string { return "system.configurations" }

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Title     string    `gorm:"not null" json:"title"`
	Message   string    `gorm:"not null" json:"message"`
	Type      string    `gorm:"default:'info';index" json:"type"`
	IsRead    bool      `gorm:"default:false;index" json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

func (Notification) TableName() string { return "system.notifications" }

type ScheduledJob struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	JobName      string     `gorm:"unique;not null;index" json:"job_name"`
	LastRun      *time.Time `json:"last_run,omitempty"`
	NextRun      *time.Time `json:"next_run,omitempty"`
	Status       string     `gorm:"type:varchar(20);default:'pending';index" json:"status"`
	ErrorMessage string     `json:"error_message,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

func (ScheduledJob) TableName() string { return "system.scheduled_jobs" }
