package academicmod

import (
        "university-erp-backend/internal/domain"

        "gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

// Academic Terms
func (r *Repository) ListTerms(campusID uint) ([]domain.AcademicTerm, error) {
        var list []domain.AcademicTerm
        q := r.db.Order("start_date DESC")
        if campusID > 0 {
                q = q.Where("campus_id = ?", campusID)
        }
        return list, q.Find(&list).Error
}
func (r *Repository) GetCurrentTerm() (*domain.AcademicTerm, error) {
        var t domain.AcademicTerm
        return &t, r.db.Where("is_current = true").First(&t).Error
}
func (r *Repository) GetTerm(id uint) (*domain.AcademicTerm, error) {
        var t domain.AcademicTerm
        return &t, r.db.First(&t, id).Error
}
func (r *Repository) CreateTerm(t *domain.AcademicTerm) error {
        return r.db.Create(t).Error
}
func (r *Repository) UpdateTerm(t *domain.AcademicTerm) error {
        return r.db.Save(t).Error
}
func (r *Repository) SetCurrentTerm(id uint) error {
        if err := r.db.Model(&domain.AcademicTerm{}).Where("is_current = true").Update("is_current", false).Error; err != nil {
                return err
        }
        return r.db.Model(&domain.AcademicTerm{}).Where("id = ?", id).Update("is_current", true).Error
}

// Programs
func (r *Repository) ListPrograms(departmentID uint) ([]domain.Program, error) {
        var list []domain.Program
        q := r.db.Where("is_active = true").Order("name")
        if departmentID > 0 {
                q = q.Where("department_id = ?", departmentID)
        }
        return list, q.Find(&list).Error
}
func (r *Repository) GetProgram(id uint) (*domain.Program, error) {
        var p domain.Program
        return &p, r.db.First(&p, id).Error
}
func (r *Repository) CreateProgram(p *domain.Program) error {
        return r.db.Create(p).Error
}
func (r *Repository) UpdateProgram(p *domain.Program) error {
        return r.db.Save(p).Error
}

// Program Semesters
func (r *Repository) ListProgramSemesters(programID uint) ([]domain.ProgramSemester, error) {
        var list []domain.ProgramSemester
        return list, r.db.Where("program_id = ?", programID).Order("semester_number").Find(&list).Error
}
func (r *Repository) CreateProgramSemester(ps *domain.ProgramSemester) error {
        return r.db.Create(ps).Error
}

// Subjects
func (r *Repository) ListSubjects(departmentID uint) ([]domain.Subject, error) {
        var list []domain.Subject
        q := r.db.Where("is_active = true").Order("subject_name")
        if departmentID > 0 {
                q = q.Where("department_id = ?", departmentID)
        }
        return list, q.Find(&list).Error
}
func (r *Repository) GetSubject(id uint) (*domain.Subject, error) {
        var s domain.Subject
        return &s, r.db.First(&s, id).Error
}
func (r *Repository) CreateSubject(s *domain.Subject) error {
        return r.db.Create(s).Error
}
func (r *Repository) UpdateSubject(s *domain.Subject) error {
        return r.db.Save(s).Error
}

// Program Subjects (Curriculum)
func (r *Repository) GetProgramCurriculum(programID uint) ([]domain.ProgramSubject, error) {
        var list []domain.ProgramSubject
        return list, r.db.Where("program_id = ?", programID).Order("semester_number").Find(&list).Error
}
func (r *Repository) AddSubjectToProgram(ps *domain.ProgramSubject) error {
        return r.db.Save(ps).Error
}
func (r *Repository) RemoveSubjectFromProgram(programID, subjectID uint) error {
        return r.db.Where("program_id = ? AND subject_id = ?", programID, subjectID).Delete(&domain.ProgramSubject{}).Error
}

// Batches
func (r *Repository) ListBatches(programID uint) ([]domain.Batch, error) {
        var list []domain.Batch
        q := r.db.Order("batch_year DESC")
        if programID > 0 {
                q = q.Where("program_id = ?", programID)
        }
        return list, q.Find(&list).Error
}
func (r *Repository) GetBatch(id uint) (*domain.Batch, error) {
        var b domain.Batch
        return &b, r.db.First(&b, id).Error
}
func (r *Repository) CreateBatch(b *domain.Batch) error {
        return r.db.Create(b).Error
}

// Sections
func (r *Repository) ListSections(batchID uint) ([]domain.Section, error) {
        var list []domain.Section
        return list, r.db.Where("batch_id = ?", batchID).Order("section_name").Find(&list).Error
}
func (r *Repository) CreateSection(s *domain.Section) error {
        return r.db.Create(s).Error
}

// Course Offerings
func (r *Repository) ListOfferings(termID, programID uint) ([]domain.CourseOffering, error) {
        var list []domain.CourseOffering
        q := r.db.Where("status = 'Active'")
        if termID > 0 {
                q = q.Where("academic_term_id = ?", termID)
        }
        if programID > 0 {
                q = q.Where("program_id = ?", programID)
        }
        return list, q.Find(&list).Error
}
func (r *Repository) GetOffering(id uint) (*domain.CourseOffering, error) {
        var o domain.CourseOffering
        return &o, r.db.First(&o, id).Error
}
func (r *Repository) CreateOffering(o *domain.CourseOffering) error {
        return r.db.Create(o).Error
}
func (r *Repository) UpdateOffering(o *domain.CourseOffering) error {
        return r.db.Save(o).Error
}

// Term Registrations
func (r *Repository) GetStudentTermRegistration(studentID, termID uint) (*domain.TermRegistration, error) {
        var tr domain.TermRegistration
        return &tr, r.db.Where("student_id = ? AND academic_term_id = ?", studentID, termID).First(&tr).Error
}
func (r *Repository) CreateTermRegistration(tr *domain.TermRegistration) error {
        return r.db.Create(tr).Error
}
func (r *Repository) ListTermRegistrations(termID uint) ([]domain.TermRegistration, error) {
        var list []domain.TermRegistration
        return list, r.db.Where("academic_term_id = ?", termID).Find(&list).Error
}

// Course Registrations
func (r *Repository) ListStudentCourseRegistrations(studentID, termID uint) ([]domain.CourseRegistration, error) {
        var list []domain.CourseRegistration
        q := r.db.Where("student_id = ?", studentID)
        if termID > 0 {
                q = q.Joins("JOIN academic.course_offerings o ON o.id = academic.course_registrations.offering_id").
                        Where("o.academic_term_id = ?", termID)
        }
        return list, q.Find(&list).Error
}
func (r *Repository) CreateCourseRegistration(cr *domain.CourseRegistration) error {
        return r.db.Create(cr).Error
}
func (r *Repository) DropCourseRegistration(id uint) error {
        return r.db.Model(&domain.CourseRegistration{}).Where("id = ?", id).Update("registration_status", "Dropped").Error
}

// Timetable
func (r *Repository) GetTimetable(offeringID uint) ([]domain.Timetable, error) {
        var list []domain.Timetable
        return list, r.db.Where("offering_id = ?", offeringID).Order("day_of_week, start_time").Find(&list).Error
}
func (r *Repository) GetStudentTimetable(studentID, termID uint) ([]domain.Timetable, error) {
        var list []domain.Timetable
        return list, r.db.Raw(`
                SELECT t.* FROM academic.timetable t
                JOIN academic.course_offerings o ON o.id = t.offering_id
                JOIN academic.course_registrations cr ON cr.offering_id = o.id
                WHERE cr.student_id = ? AND o.academic_term_id = ? AND cr.registration_status = 'Enrolled'
                ORDER BY t.day_of_week, t.start_time
        `, studentID, termID).Scan(&list).Error
}
func (r *Repository) CreateTimetableEntry(t *domain.Timetable) error {
        return r.db.Create(t).Error
}

// Calendar
func (r *Repository) ListCalendar(campusID uint) ([]domain.AcademicCalendar, error) {
        var list []domain.AcademicCalendar
        q := r.db.Order("event_date")
        if campusID > 0 {
                q = q.Where("campus_id = ?", campusID)
        }
        return list, q.Find(&list).Error
}
func (r *Repository) CreateCalendarEvent(e *domain.AcademicCalendar) error {
        return r.db.Create(e).Error
}

// Student Attendance - works via ClassSession → StudentAttendance
func (r *Repository) CreateClassSession(cs *domain.ClassSession) error {
        return r.db.Create(cs).Error
}
func (r *Repository) GetClassSessions(offeringID uint) ([]domain.ClassSession, error) {
        var list []domain.ClassSession
        return list, r.db.Where("offering_id = ?", offeringID).Order("class_date DESC").Find(&list).Error
}
func (r *Repository) MarkAttendance(a *domain.StudentAttendance) error {
        return r.db.Save(a).Error
}
func (r *Repository) GetStudentAttendance(studentID, offeringID uint) ([]domain.StudentAttendance, error) {
        var list []domain.StudentAttendance
        q := r.db.Where("student_id = ?", studentID)
        if offeringID > 0 {
                q = q.Joins("JOIN student.class_sessions cs ON cs.id = student.attendance.session_id").
                        Where("cs.offering_id = ?", offeringID)
        }
        return list, q.Order("created_at DESC").Find(&list).Error
}
func (r *Repository) GetAttendanceSummary(studentID, termID uint) ([]AttendanceSummary, error) {
        var result []AttendanceSummary
        return result, r.db.Raw(`
                SELECT o.id as offering_id, sub.subject_name, sub.subject_code,
                        COUNT(*) as total_classes,
                        SUM(CASE WHEN sc.code = 'PRESENT' THEN 1 ELSE 0 END) as present_count,
                        ROUND(100.0 * SUM(CASE WHEN sc.code = 'PRESENT' THEN 1 ELSE 0 END) / NULLIF(COUNT(*),0), 2) as percentage
                FROM student.attendance sa
                JOIN student.class_sessions cs ON cs.id = sa.session_id
                JOIN academic.course_offerings o ON o.id = cs.offering_id
                JOIN academic.subjects sub ON sub.id = o.subject_id
                LEFT JOIN system.status_codes sc ON sc.id = sa.status_id
                WHERE sa.student_id = ? AND o.academic_term_id = ?
                GROUP BY o.id, sub.subject_name, sub.subject_code
        `, studentID, termID).Scan(&result).Error
}

type AttendanceSummary struct {
        OfferingID    uint    `json:"offering_id"`
        SubjectName   string  `json:"subject_name"`
        SubjectCode   string  `json:"subject_code"`
        TotalClasses  int     `json:"total_classes"`
        PresentCount  int     `json:"present_count"`
        Percentage    float64 `json:"percentage"`
}
