package domain

import "time"

// ─── Finance: Fee Structure ──────────────────────────────────────────────────

type FeeHead struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Code        string    `gorm:"unique;not null;index" json:"code"`
	Description string    `json:"description"`
	IsMandatory bool      `gorm:"default:true" json:"is_mandatory"`
	CreatedAt   time.Time `json:"created_at"`
}

func (FeeHead) TableName() string { return "finance.fee_heads" }

type FeeStructure struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ProgramID      uint      `gorm:"not null;index" json:"program_id"`
	SemesterNumber int       `gorm:"not null" json:"semester_number"`
	FeeHeadID      uint      `gorm:"not null;index" json:"fee_head_id"`
	Amount         float64   `gorm:"not null" json:"amount"`
	AcademicYear   string    `gorm:"not null;index" json:"academic_year"`
	IsActive       bool      `gorm:"default:true;index" json:"is_active"`
	CreatedBy      *uint     `json:"created_by,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

func (FeeStructure) TableName() string { return "finance.fee_structures" }

// ─── Finance: Invoices & Payments ────────────────────────────────────────────

type Invoice struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	StudentID      uint      `gorm:"not null;index" json:"student_id"`
	InvoiceNumber  string    `gorm:"unique;not null;index" json:"invoice_number"`
	AcademicTermID uint      `gorm:"not null;index" json:"academic_term_id"`
	GeneratedDate  time.Time `gorm:"default:CURRENT_DATE;index" json:"generated_date"`
	DueDate        time.Time `gorm:"not null;index" json:"due_date"`
	TotalAmount    float64   `gorm:"not null" json:"total_amount"`
	PaidAmount     float64   `gorm:"default:0" json:"paid_amount"`
	LateFeeApplied float64   `gorm:"default:0" json:"late_fee_applied"`
	StatusID       *uint     `gorm:"index" json:"status_id,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (Invoice) TableName() string { return "finance.invoices" }

type InvoiceItem struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	InvoiceID   uint    `gorm:"not null;index" json:"invoice_id"`
	FeeHeadID   uint    `gorm:"not null" json:"fee_head_id"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitAmount  float64 `json:"unit_amount"`
	Amount      float64 `json:"amount"`
}

func (InvoiceItem) TableName() string { return "finance.invoice_items" }

type Payment struct {
	ID                       uint      `gorm:"primaryKey" json:"id"`
	InvoiceID                uint      `gorm:"not null;index" json:"invoice_id"`
	StudentID                uint      `gorm:"not null;index" json:"student_id"`
	Amount                   float64   `gorm:"not null" json:"amount"`
	PaymentDate              time.Time `gorm:"default:CURRENT_TIMESTAMP;index" json:"payment_date"`
	PaymentModeID            *uint     `gorm:"index" json:"payment_mode_id,omitempty"`
	TransactionID            string    `gorm:"index" json:"transaction_id"`
	ReferenceNo              string    `json:"reference_no"`
	StatusID                 *uint     `gorm:"index" json:"status_id,omitempty"`
	ReceiptURL               string    `json:"receipt_url"`
	BankReconciliationStatus string    `gorm:"type:varchar(20);default:'Pending'" json:"bank_reconciliation_status"`
	CreatedAt                time.Time `json:"created_at"`
}

func (Payment) TableName() string { return "finance.payments" }

type PaymentAllocation struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	PaymentID       uint      `gorm:"not null;index" json:"payment_id"`
	InvoiceID       uint      `gorm:"not null;index" json:"invoice_id"`
	AllocatedAmount float64   `json:"allocated_amount"`
	AllocatedAt     time.Time `json:"allocated_at"`
}

func (PaymentAllocation) TableName() string { return "finance.payment_allocations" }

// ─── Finance: Scholarships & Discounts ───────────────────────────────────────

type Scholarship struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	Name                string    `gorm:"not null" json:"name"`
	Description         string    `json:"description"`
	EligibilityCriteria string    `gorm:"type:jsonb" json:"eligibility_criteria"`
	Amount              float64   `json:"amount"`
	Renewable           bool      `gorm:"default:false" json:"renewable"`
	CreatedAt           time.Time `json:"created_at"`
}

func (Scholarship) TableName() string { return "finance.scholarships" }

type StudentScholarship struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	StudentID     uint       `gorm:"not null;index" json:"student_id"`
	ScholarshipID uint       `gorm:"not null;index" json:"scholarship_id"`
	AcademicYear  string     `gorm:"not null;index" json:"academic_year"`
	AmountAwarded float64    `json:"amount_awarded"`
	Disbursed     bool       `gorm:"default:false;index" json:"disbursed"`
	DisbursedAt   *time.Time `json:"disbursed_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

func (StudentScholarship) TableName() string { return "finance.student_scholarships" }

type StudentDiscount struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	StudentID    uint       `gorm:"not null;index" json:"student_id"`
	FeeHeadID    uint       `gorm:"not null;index" json:"fee_head_id"`
	AcademicYear string     `gorm:"not null;index" json:"academic_year"`
	Amount       float64    `gorm:"not null" json:"amount"`
	Reason       string     `json:"reason"`
	ApprovedBy   *uint      `json:"approved_by,omitempty"`
	ApprovedAt   *time.Time `json:"approved_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

func (StudentDiscount) TableName() string { return "finance.student_discounts" }

type InstallmentPlan struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	StudentID      uint      `gorm:"not null;index" json:"student_id"`
	AcademicTermID uint      `gorm:"not null;index" json:"academic_term_id"`
	DueDate        time.Time `gorm:"not null;index" json:"due_date"`
	Amount         float64   `gorm:"not null" json:"amount"`
	PaidAmount     float64   `gorm:"default:0" json:"paid_amount"`
	StatusID       *uint     `gorm:"index" json:"status_id,omitempty"`
	LateFee        float64   `gorm:"default:0" json:"late_fee"`
	CreatedAt      time.Time `json:"created_at"`
}

func (InstallmentPlan) TableName() string { return "finance.installment_plans" }

type Refund struct {
	ID                  uint       `gorm:"primaryKey" json:"id"`
	PaymentID           uint       `gorm:"not null;index" json:"payment_id"`
	StudentID           uint       `gorm:"not null;index" json:"student_id"`
	Amount              float64    `gorm:"not null" json:"amount"`
	Reason              string     `json:"reason"`
	StatusID            *uint      `gorm:"index" json:"status_id,omitempty"`
	ApprovedBy          *uint      `json:"approved_by,omitempty"`
	ProcessedAt         *time.Time `json:"processed_at,omitempty"`
	RefundTransactionID string     `gorm:"index" json:"refund_transaction_id"`
	CreatedAt           time.Time  `json:"created_at"`
}

func (Refund) TableName() string { return "finance.refunds" }
