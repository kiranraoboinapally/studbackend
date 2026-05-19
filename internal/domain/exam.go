package domain

import "time"

// ─── Exam Components & Schedule ──────────────────────────────────────────────

type ExamComponent struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	SubjectID     uint   `gorm:"not null;index" json:"subject_id"`
	ComponentName string `json:"component_name"`
	MaxMarks      int    `json:"max_marks"`
	DisplayOrder  int    `json:"display_order"`
}

func (ExamComponent) TableName() string { return "exam.exam_components" }

type ExamSchedule struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	SubjectID      uint      `gorm:"not null;index" json:"subject_id"`
	AcademicTermID uint      `gorm:"not null;index" json:"academic_term_id"`
	ExamDate       time.Time `gorm:"not null;index" json:"exam_date"`
	StartTime      string    `json:"start_time"`
	EndTime        string    `json:"end_time"`
	ExamType       string    `gorm:"type:varchar(30)" json:"exam_type"`
	Venue          string    `json:"venue"`
	TotalMarks     int       `gorm:"not null" json:"total_marks"`
	PassingMarks   int       `json:"passing_marks"`
	CreatedAt      time.Time `json:"created_at"`
}

func (ExamSchedule) TableName() string { return "exam.exam_schedules" }

// ─── Results & Grades ────────────────────────────────────────────────────────

type ComponentMarks struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	ResultID      uint    `gorm:"not null;index" json:"result_id"`
	ComponentID   uint    `gorm:"not null" json:"component_id"`
	MarksObtained float64 `json:"marks_obtained"`
}

func (ComponentMarks) TableName() string { return "exam.component_marks" }

type Result struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	StudentID      uint       `gorm:"not null;index" json:"student_id"`
	CourseRegID    uint       `gorm:"not null;index" json:"course_reg_id"`
	SubjectID      uint       `gorm:"not null;index" json:"subject_id"`
	AcademicTermID uint       `gorm:"not null;index" json:"academic_term_id"`
	MarksObtained  float64    `json:"marks_obtained"`
	MaxMarks       float64    `json:"max_marks"`
	Grade          string     `gorm:"index" json:"grade"`
	GradePoint     float64    `json:"grade_point"`
	IsPassed       bool       `json:"is_passed"`
	PublishedAt    *time.Time `json:"published_at,omitempty"`
	PublishedBy    *uint      `json:"published_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

func (Result) TableName() string { return "exam.results" }

type RevaluationRequest struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	ResultID      uint       `gorm:"not null;index" json:"result_id"`
	StudentID     uint       `gorm:"not null;index" json:"student_id"`
	SubjectID     uint       `gorm:"not null;index" json:"subject_id"`
	RequestDate   time.Time  `gorm:"default:CURRENT_TIMESTAMP;index" json:"request_date"`
	StatusID      *uint      `gorm:"index" json:"status_id,omitempty"`
	ReviewedMarks float64    `json:"reviewed_marks"`
	ReviewedGrade string     `json:"reviewed_grade"`
	Remarks       string     `json:"remarks"`
	ProcessedAt   *time.Time `json:"processed_at,omitempty"`
	FeePaid       float64    `json:"fee_paid"`
}

func (RevaluationRequest) TableName() string { return "exam.revaluation_requests" }

type SupplementaryExam struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	SubjectID      uint       `gorm:"not null;index" json:"subject_id"`
	AcademicTermID uint       `gorm:"not null;index" json:"academic_term_id"`
	ExamDate       *time.Time `json:"exam_date,omitempty"`
	ResultDeclared bool       `gorm:"default:false" json:"result_declared"`
	CreatedAt      time.Time  `json:"created_at"`
}

func (SupplementaryExam) TableName() string { return "exam.supplementary_exams" }
