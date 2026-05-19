package domain

import "time"

// ─── Core Organization ───────────────────────────────────────────────────────

type University struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Name            string    `gorm:"not null" json:"name"`
	ShortName       string    `gorm:"unique;not null;index" json:"short_name"`
	EstablishedYear int       `json:"established_year"`
	Address         string    `json:"address"`
	City            string    `json:"city"`
	State           string    `json:"state"`
	PostalCode      string    `json:"postal_code"`
	Phone           string    `json:"phone"`
	Email           string    `json:"email"`
	Website         string    `json:"website"`
	Vision          string    `json:"vision"`
	Mission         string    `json:"mission"`
	IsActive        bool      `gorm:"default:true;index" json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (University) TableName() string { return "core.universities" }

type Campus struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UniversityID uint      `gorm:"not null;index" json:"university_id"`
	Name         string    `gorm:"not null" json:"name"`
	Code         string    `gorm:"unique;not null;index" json:"code"`
	Address      string    `json:"address"`
	City         string    `json:"city"`
	State        string    `json:"state"`
	PostalCode   string    `json:"postal_code"`
	Phone        string    `json:"phone"`
	IsMainCampus bool      `gorm:"default:false" json:"is_main_campus"`
	IsActive     bool      `gorm:"default:true;index" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (Campus) TableName() string { return "core.campuses" }

type Department struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	CampusID           *uint     `gorm:"index" json:"campus_id,omitempty"`
	Name               string    `gorm:"not null" json:"name"`
	Code               string    `gorm:"unique;not null;index" json:"code"`
	ParentDepartmentID *uint     `gorm:"index" json:"parent_department_id,omitempty"`
	EstablishedYear    int       `json:"established_year"`
	HodEmployeeID      *uint    `gorm:"index" json:"hod_employee_id,omitempty"`
	Description        string    `json:"description"`
	IsActive           bool      `gorm:"default:true;index" json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (Department) TableName() string { return "core.departments" }

type Room struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CampusID   uint      `gorm:"not null;index" json:"campus_id"`
	RoomNumber string    `gorm:"not null" json:"room_number"`
	RoomType   string    `gorm:"type:varchar(50);index" json:"room_type"`
	Capacity   int       `gorm:"not null" json:"capacity"`
	Building   string    `json:"building"`
	Floor      int       `json:"floor"`
	IsActive   bool      `gorm:"default:true;index" json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}

func (Room) TableName() string { return "core.rooms" }
