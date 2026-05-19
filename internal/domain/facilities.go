package domain

import "time"

// ─── Hostel ──────────────────────────────────────────────────────────────────

type Hostel struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"not null" json:"name"`
	Code          string    `gorm:"unique;not null;index" json:"code"`
	CampusID      *uint     `gorm:"index" json:"campus_id,omitempty"`
	GenderID      *uint     `gorm:"index" json:"gender_id,omitempty"`
	TotalRooms    int       `json:"total_rooms"`
	WardenID      *uint     `gorm:"index" json:"warden_id,omitempty"`
	ContactNumber string    `json:"contact_number"`
	Address       string    `json:"address"`
	IsActive      bool      `gorm:"default:true;index" json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
}

func (Hostel) TableName() string { return "hostel.hostels" }

type HostelRoom struct {
	ID               uint    `gorm:"primaryKey" json:"id"`
	HostelID         uint    `gorm:"not null;index" json:"hostel_id"`
	RoomNumber       string  `gorm:"not null" json:"room_number"`
	RoomType         string  `gorm:"type:varchar(20);index" json:"room_type"`
	Capacity         int     `gorm:"not null" json:"capacity"`
	CurrentOccupancy int     `gorm:"default:0" json:"current_occupancy"`
	MonthlyRent      float64 `json:"monthly_rent"`
	IsAvailable      bool    `gorm:"default:true;index" json:"is_available"`
	CreatedAt        time.Time `json:"created_at"`
}

func (HostelRoom) TableName() string { return "hostel.rooms" }

type HostelBed struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	RoomID     uint   `gorm:"not null;index" json:"room_id"`
	BedNumber  string `json:"bed_number"`
	IsOccupied bool   `gorm:"default:false" json:"is_occupied"`
}

func (HostelBed) TableName() string { return "hostel.beds" }

type HostelAllocation struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	StudentID     uint       `gorm:"not null;index" json:"student_id"`
	RoomID        uint       `gorm:"not null;index" json:"room_id"`
	BedID         *uint      `gorm:"index" json:"bed_id,omitempty"`
	AllocatedFrom time.Time  `gorm:"not null;index" json:"allocated_from"`
	AllocatedTo   *time.Time `json:"allocated_to,omitempty"`
	StatusID      *uint      `gorm:"index" json:"status_id,omitempty"`
	CreatedBy     *uint      `json:"created_by,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (HostelAllocation) TableName() string { return "hostel.allocations" }

type HostelAllocationHistory struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	StudentID     uint      `gorm:"not null;index" json:"student_id"`
	RoomID        uint      `gorm:"not null;index" json:"room_id"`
	AllocatedFrom time.Time `gorm:"not null" json:"allocated_from"`
	AllocatedTo   time.Time `gorm:"not null" json:"allocated_to"`
	Reason        string    `json:"reason"`
	CreatedAt     time.Time `json:"created_at"`
}

func (HostelAllocationHistory) TableName() string { return "hostel.allocation_history" }

type MessBill struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	StudentID uint       `gorm:"not null;index" json:"student_id"`
	Month     time.Time  `gorm:"not null;index" json:"month"`
	Amount    float64    `gorm:"not null" json:"amount"`
	Paid      bool       `gorm:"default:false;index" json:"paid"`
	PaidAt    *time.Time `json:"paid_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

func (MessBill) TableName() string { return "hostel.mess_bills" }

type MaintenanceRequest struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	StudentID   uint       `gorm:"not null;index" json:"student_id"`
	RoomID      uint       `gorm:"not null;index" json:"room_id"`
	Category    string     `gorm:"index" json:"category"`
	Description string     `gorm:"not null" json:"description"`
	StatusID    *uint      `gorm:"index" json:"status_id,omitempty"`
	AssignedTo  *uint      `json:"assigned_to,omitempty"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (MaintenanceRequest) TableName() string { return "hostel.maintenance_requests" }

type VisitorLog struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	HostelID    uint       `gorm:"not null;index" json:"hostel_id"`
	VisitorName string     `gorm:"not null" json:"visitor_name"`
	StudentID   *uint      `gorm:"index" json:"student_id,omitempty"`
	EntryTime   time.Time  `gorm:"default:CURRENT_TIMESTAMP;index" json:"entry_time"`
	ExitTime    *time.Time `json:"exit_time,omitempty"`
	Purpose     string     `json:"purpose"`
	IDProof     string     `json:"id_proof"`
}

func (VisitorLog) TableName() string { return "hostel.visitor_logs" }

// ─── Transport ───────────────────────────────────────────────────────────────

type Bus struct {
	ID                  uint       `gorm:"primaryKey" json:"id"`
	BusNumber           string     `gorm:"unique;not null;index" json:"bus_number"`
	RegistrationNo      string     `gorm:"unique;not null;index" json:"registration_no"`
	Capacity            int        `gorm:"not null" json:"capacity"`
	DriverEmployeeID    *uint      `gorm:"index" json:"driver_employee_id,omitempty"`
	DriverLicenseExpiry *time.Time `json:"driver_license_expiry,omitempty"`
	IsActive            bool       `gorm:"default:true;index" json:"is_active"`
	CreatedAt           time.Time  `json:"created_at"`
}

func (Bus) TableName() string { return "transport.buses" }

type Route struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	RouteName     string    `gorm:"not null;index" json:"route_name"`
	Description   string    `json:"description"`
	DistanceKm    float64   `json:"distance_km"`
	EstimatedTime string    `json:"estimated_time"`
	IsActive      bool      `gorm:"default:true;index" json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
}

func (Route) TableName() string { return "transport.routes" }

type Stop struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	RouteID       uint      `gorm:"not null;index" json:"route_id"`
	StopName      string    `gorm:"not null" json:"stop_name"`
	StopOrder     int       `gorm:"not null" json:"stop_order"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	ArrivalTime   string    `json:"arrival_time"`
	DepartureTime string    `json:"departure_time"`
	CreatedAt     time.Time `json:"created_at"`
}

func (Stop) TableName() string { return "transport.stops" }

type BusAssignment struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	BusID         uint       `gorm:"not null;index" json:"bus_id"`
	RouteID       uint       `gorm:"not null;index" json:"route_id"`
	EffectiveFrom time.Time  `gorm:"not null;index" json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to,omitempty"`
	IsActive      bool       `gorm:"default:true;index" json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
}

func (BusAssignment) TableName() string { return "transport.bus_assignments" }

type StudentPass struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	StudentID    uint      `gorm:"not null;index" json:"student_id"`
	RouteID      uint      `gorm:"not null;index" json:"route_id"`
	PickupStopID uint      `gorm:"not null" json:"pickup_stop_id"`
	DropStopID   uint      `gorm:"not null" json:"drop_stop_id"`
	ValidFrom    time.Time `gorm:"not null;index" json:"valid_from"`
	ValidTo      time.Time `gorm:"not null;index" json:"valid_to"`
	FeePaid      float64   `json:"fee_paid"`
	StatusID     *uint     `gorm:"index" json:"status_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

func (StudentPass) TableName() string { return "transport.student_passes" }

type VehicleMaintenance struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	BusID           uint       `gorm:"not null;index" json:"bus_id"`
	MaintenanceDate time.Time  `gorm:"not null;index" json:"maintenance_date"`
	Description     string     `json:"description"`
	Cost            float64    `json:"cost"`
	NextDueDate     *time.Time `json:"next_due_date,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

func (VehicleMaintenance) TableName() string { return "transport.vehicle_maintenance" }

// ─── Library ─────────────────────────────────────────────────────────────────

type Author struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null;index" json:"name"`
	Biography string    `json:"biography"`
	CreatedAt time.Time `json:"created_at"`
}

func (Author) TableName() string { return "library.authors" }

type Book struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Title           string    `gorm:"not null;index" json:"title"`
	ISBN            string    `gorm:"unique;index" json:"isbn"`
	Publisher       string    `json:"publisher"`
	PublicationYear int       `json:"publication_year"`
	Edition         string    `json:"edition"`
	TotalCopies     int       `gorm:"default:1" json:"total_copies"`
	AvailableCopies int       `gorm:"default:1" json:"available_copies"`
	Location        string    `json:"location"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (Book) TableName() string { return "library.books" }

type BookCopy struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	BookID        uint      `gorm:"not null;index" json:"book_id"`
	Barcode       string    `gorm:"unique;index" json:"barcode"`
	CopyNumber    int       `json:"copy_number"`
	Condition     string    `gorm:"type:varchar(20)" json:"condition"`
	ShelfLocation string    `json:"shelf_location"`
	StatusID      *uint     `gorm:"index" json:"status_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

func (BookCopy) TableName() string { return "library.book_copies" }

type BookAuthor struct {
	BookID   uint `gorm:"primaryKey" json:"book_id"`
	AuthorID uint `gorm:"primaryKey" json:"author_id"`
}

func (BookAuthor) TableName() string { return "library.book_authors" }

type DigitalResource struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	Title             string     `gorm:"not null;index" json:"title"`
	ResourceType      string     `gorm:"type:varchar(50);index" json:"resource_type"`
	URL               string     `json:"url"`
	AccessLink        string     `json:"access_link"`
	Publisher         string     `json:"publisher"`
	LicenseValidUntil *time.Time `json:"license_valid_until,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
}

func (DigitalResource) TableName() string { return "library.digital_resources" }

type Circulation struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	BookCopyID   uint       `gorm:"not null;index" json:"book_copy_id"`
	StudentID    uint       `gorm:"not null;index" json:"student_id"`
	IssuedDate   time.Time  `gorm:"default:CURRENT_DATE;index" json:"issued_date"`
	DueDate      time.Time  `gorm:"not null;index" json:"due_date"`
	ReturnedDate *time.Time `gorm:"index" json:"returned_date,omitempty"`
	StatusID     *uint      `gorm:"index" json:"status_id,omitempty"`
	FineAmount   float64    `gorm:"default:0" json:"fine_amount"`
	FinePaid     bool       `gorm:"default:false" json:"fine_paid"`
	IssuedBy     *uint      `json:"issued_by,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

func (Circulation) TableName() string { return "library.circulations" }

type Reservation struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	BookID        uint       `gorm:"not null;index" json:"book_id"`
	StudentID     uint       `gorm:"not null;index" json:"student_id"`
	ReservedFrom  time.Time  `gorm:"default:CURRENT_TIMESTAMP;index" json:"reserved_from"`
	ReservedUntil *time.Time `json:"reserved_until,omitempty"`
	StatusID      *uint      `gorm:"index" json:"status_id,omitempty"`
	NotifiedAt    *time.Time `json:"notified_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

func (Reservation) TableName() string { return "library.reservations" }

type LibraryFine struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	CirculationID uint       `gorm:"not null;index" json:"circulation_id"`
	Amount        float64    `gorm:"not null" json:"amount"`
	Reason        string     `json:"reason"`
	PaidDate      *time.Time `json:"paid_date,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

func (LibraryFine) TableName() string { return "library.fines" }

type PurchaseRequest struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	RequestedBy uint      `gorm:"not null;index" json:"requested_by"`
	Title       string     `gorm:"not null" json:"title"`
	Author      string     `json:"author"`
	ISBN        string     `json:"isbn"`
	Reason      string     `json:"reason"`
	StatusID    *uint      `gorm:"index" json:"status_id,omitempty"`
	ApprovedBy  *uint      `json:"approved_by,omitempty"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (PurchaseRequest) TableName() string { return "library.purchase_requests" }

// ─── Security & Access Control ───────────────────────────────────────────────

type Permission struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Resource    string `gorm:"not null;index" json:"resource"`
	Action      string `gorm:"not null" json:"action"`
	Description string `json:"description"`
}

func (Permission) TableName() string { return "security.permissions" }

type RolePermission struct {
	RoleID       uint      `gorm:"primaryKey" json:"role_id"`
	PermissionID uint      `gorm:"primaryKey" json:"permission_id"`
	GrantedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"granted_at"`
	GrantedBy    *uint     `json:"granted_by,omitempty"`
}

func (RolePermission) TableName() string { return "security.role_permissions" }

type UserSession struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	UserID       uint       `gorm:"not null;index" json:"user_id"`
	SessionToken string     `gorm:"unique;not null;index" json:"-"`
	IPAddress    string     `json:"ip_address"`
	UserAgent    string     `json:"user_agent"`
	LoginAt      time.Time  `gorm:"default:CURRENT_TIMESTAMP;index" json:"login_at"`
	LastActivity time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"last_activity"`
	LogoutAt     *time.Time `json:"logout_at,omitempty"`
	IsActive     bool       `gorm:"default:true;index" json:"is_active"`
}

func (UserSession) TableName() string { return "security.user_sessions" }

type LoginAttempt struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        *uint     `gorm:"index" json:"user_id,omitempty"`
	Username      string    `gorm:"not null;index" json:"username"`
	Success       bool      `gorm:"default:false;index" json:"success"`
	IPAddress     string    `json:"ip_address"`
	UserAgent     string    `json:"user_agent"`
	FailureReason string    `json:"failure_reason"`
	AttemptedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP;index" json:"attempted_at"`
}

func (LoginAttempt) TableName() string { return "security.login_attempts" }

type PasswordReset struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"not null;index" json:"user_id"`
	Token     string     `gorm:"unique;not null;index" json:"-"`
	ExpiresAt time.Time  `gorm:"not null;index" json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

func (PasswordReset) TableName() string { return "security.password_resets" }

type APIKey struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	UserID     uint       `gorm:"not null;index" json:"user_id"`
	KeyHash    string     `gorm:"unique;not null;index" json:"-"`
	Name       string     `json:"name"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	IsActive   bool       `gorm:"default:true;index" json:"is_active"`
	CreatedAt  time.Time  `json:"created_at"`
}

func (APIKey) TableName() string { return "security.api_keys" }

// ─── Audit ───────────────────────────────────────────────────────────────────

type SystemEvent struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	EventType    string    `gorm:"not null;index" json:"event_type"`
	Severity     string    `gorm:"type:varchar(20);index" json:"severity"`
	SourceModule string    `gorm:"index" json:"source_module"`
	Message      string    `json:"message"`
	Details      string    `gorm:"type:jsonb" json:"details"`
	IPAddress    string    `json:"ip_address"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP;index" json:"created_at"`
}

func (SystemEvent) TableName() string { return "audit.system_events" }
