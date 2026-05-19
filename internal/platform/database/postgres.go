package database

import (
	"fmt"
	"log"
	"strings"

	"university-erp-backend/internal/config"
	"university-erp-backend/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// AllSchemas defines every PostgreSQL schema the ERP uses.
var AllSchemas = []string{
	"shared", "system", "core", "academic", "hr",
	"student", "admissions", "finance", "exam",
	"hostel", "transport", "library", "security", "audit",
}

// Connect establishes a GORM connection, creates schemas, and runs AutoMigrate.
func Connect(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	logLevel := logger.Warn
	if cfg.AppEnv == "development" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)

	log.Println("✅ Database connected")

	createSchemas(db)
	autoMigrate(db)
	installTriggers(db)

	return db
}

func createSchemas(db *gorm.DB) {
	for _, s := range AllSchemas {
		if err := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", s)).Error; err != nil {
			log.Fatalf("❌ Failed to create schema %s: %v", s, err)
		}
	}
	log.Println("✅ All schemas ready")
}

func autoMigrate(db *gorm.DB) {
	models := []interface{}{
		// Shared / Identity
		&domain.User{}, &domain.Role{}, &domain.UserRole{}, &domain.AuditLog{},
		&domain.OutboxEvent{},
		// System lookups
		&domain.Gender{}, &domain.Category{}, &domain.BloodGroup{}, &domain.StatusCode{},
		&domain.Configuration{}, &domain.Notification{}, &domain.ScheduledJob{},
		// Core
		&domain.University{}, &domain.Campus{}, &domain.Department{}, &domain.Room{},
		// Academic
		&domain.AcademicTerm{}, &domain.Batch{}, &domain.Section{},
		&domain.Program{}, &domain.ProgramSemester{},
		&domain.Subject{}, &domain.ProgramSubject{}, &domain.SubjectPrerequisite{},
		&domain.CourseOffering{}, &domain.TermRegistration{}, &domain.CourseRegistration{},
		&domain.Timetable{}, &domain.AcademicCalendar{},
		// HR
		&domain.Designation{}, &domain.EmploymentType{}, &domain.LeaveType{},
		&domain.Employee{}, &domain.EmployeeDepartmentHistory{}, &domain.EmployeeDesignationHistory{},
		&domain.Faculty{}, &domain.Staff{},
		&domain.SalaryComponent{}, &domain.Salary{}, &domain.SalaryDetail{}, &domain.PayrollRun{},
		&domain.LeaveBalance{}, &domain.LeaveRequest{}, &domain.HRAttendance{},
		&domain.RecruitmentJob{}, &domain.JobApplication{},
		// Student
		&domain.Student{}, &domain.StudentStatusHistory{}, &domain.Guardian{}, &domain.MedicalRecord{},
		&domain.Grievance{}, &domain.ClassSession{}, &domain.StudentAttendance{},
		&domain.StudentEnrollment{}, &domain.Alumni{},
		// Admissions
		&domain.AdmissionCycle{}, &domain.Applicant{}, &domain.ApplicationStatusHistory{},
		&domain.Document{}, &domain.SeatAllocation{}, &domain.ApplicantStudentMap{}, &domain.Waitlist{},
		// Finance
		&domain.FeeHead{}, &domain.FeeStructure{},
		&domain.Invoice{}, &domain.InvoiceItem{}, &domain.Payment{}, &domain.PaymentAllocation{},
		&domain.Scholarship{}, &domain.StudentScholarship{}, &domain.StudentDiscount{},
		&domain.InstallmentPlan{}, &domain.Refund{},
		// Exam
		&domain.ExamComponent{}, &domain.ExamSchedule{},
		&domain.ComponentMarks{}, &domain.Result{},
		&domain.RevaluationRequest{}, &domain.SupplementaryExam{},
		// Hostel
		&domain.Hostel{}, &domain.HostelRoom{}, &domain.HostelBed{},
		&domain.HostelAllocation{}, &domain.HostelAllocationHistory{},
		&domain.MessBill{}, &domain.MaintenanceRequest{}, &domain.VisitorLog{},
		// Transport
		&domain.Bus{}, &domain.Route{}, &domain.Stop{},
		&domain.BusAssignment{}, &domain.StudentPass{}, &domain.VehicleMaintenance{},
		// Library
		&domain.Author{}, &domain.Book{}, &domain.BookCopy{}, &domain.BookAuthor{},
		&domain.DigitalResource{}, &domain.Circulation{}, &domain.Reservation{},
		&domain.LibraryFine{}, &domain.PurchaseRequest{},
		// Security
		&domain.Permission{}, &domain.RolePermission{},
		&domain.UserSession{}, &domain.LoginAttempt{}, &domain.PasswordReset{}, &domain.APIKey{},
		// Audit
		&domain.SystemEvent{},
	}

	for _, m := range models {
		if err := db.AutoMigrate(m); err != nil {
			log.Fatalf("❌ Migration failed for %T: %v", m, err)
		}
	}
	log.Println("✅ All tables migrated")
}

func installTriggers(db *gorm.DB) {
	// Auto-update updated_at trigger
	db.Exec(`
		CREATE OR REPLACE FUNCTION shared.update_updated_at()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = CURRENT_TIMESTAMP;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`)

	tables := []string{
		"shared.users", "core.universities", "core.campuses", "core.departments",
		"academic.programs", "academic.academic_terms", "academic.subjects",
		"hr.employees", "student.students", "finance.invoices",
	}
	for _, tbl := range tables {
		// Replace dot with underscore for trigger name (PostgreSQL doesn't allow dots in trigger names)
		triggerName := fmt.Sprintf("trg_%s_updated_at", strings.ReplaceAll(tbl, ".", "_"))
		db.Exec(fmt.Sprintf("DROP TRIGGER IF EXISTS %s ON %s", triggerName, tbl))
		db.Exec(fmt.Sprintf(
			"CREATE TRIGGER %s BEFORE UPDATE ON %s FOR EACH ROW EXECUTE FUNCTION shared.update_updated_at()",
			triggerName, tbl,
		))
	}

	// Row-Level Security on sensitive tables
	rlsTables := []string{
		"exam.results", "hr.salaries", "student.students", "finance.invoices",
	}
	for _, tbl := range rlsTables {
		db.Exec(fmt.Sprintf("ALTER TABLE %s ENABLE ROW LEVEL SECURITY", tbl))
	}

	log.Println("✅ Triggers & RLS installed")
}
