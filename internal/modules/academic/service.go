package academicmod

import (
        "context"
        "time"

        "university-erp-backend/internal/domain"
        "university-erp-backend/internal/platform/apperrors"
)

type Service struct{ repo *Repository }

func NewService(repo *Repository) *Service { return &Service{repo: repo} }

// Terms
func (s *Service) ListTerms(ctx context.Context, campusID uint) ([]domain.AcademicTerm, error) {
        return s.repo.ListTerms(campusID)
}
func (s *Service) GetCurrentTerm(ctx context.Context) (*domain.AcademicTerm, error) {
        return s.repo.GetCurrentTerm()
}
func (s *Service) GetTerm(ctx context.Context, id uint) (*domain.AcademicTerm, error) {
        return s.repo.GetTerm(id)
}
func (s *Service) CreateTerm(ctx context.Context, t *domain.AcademicTerm) error {
        if t.TermName == "" {
                return apperrors.BadRequest("term name is required")
        }
        if t.StartDate.IsZero() || t.EndDate.IsZero() {
                return apperrors.BadRequest("start and end dates are required")
        }
        return s.repo.CreateTerm(t)
}
func (s *Service) UpdateTerm(ctx context.Context, id uint, t *domain.AcademicTerm) error {
        existing, err := s.repo.GetTerm(id)
        if err != nil {
                return apperrors.NotFound("term not found")
        }
        t.ID = existing.ID
        return s.repo.UpdateTerm(t)
}
func (s *Service) SetCurrentTerm(ctx context.Context, id uint) error {
        if _, err := s.repo.GetTerm(id); err != nil {
                return apperrors.NotFound("term not found")
        }
        return s.repo.SetCurrentTerm(id)
}

// Programs
func (s *Service) ListPrograms(ctx context.Context, departmentID uint) ([]domain.Program, error) {
        return s.repo.ListPrograms(departmentID)
}
func (s *Service) GetProgram(ctx context.Context, id uint) (*domain.Program, error) {
        return s.repo.GetProgram(id)
}
func (s *Service) CreateProgram(ctx context.Context, p *domain.Program) error {
        if p.Name == "" || p.Code == "" {
                return apperrors.BadRequest("program name and code are required")
        }
        return s.repo.CreateProgram(p)
}
func (s *Service) UpdateProgram(ctx context.Context, id uint, p *domain.Program) error {
        existing, err := s.repo.GetProgram(id)
        if err != nil {
                return apperrors.NotFound("program not found")
        }
        p.ID = existing.ID
        return s.repo.UpdateProgram(p)
}
func (s *Service) GetProgramWithCurriculum(ctx context.Context, programID uint) (map[string]interface{}, error) {
        program, err := s.repo.GetProgram(programID)
        if err != nil {
                return nil, apperrors.NotFound("program not found")
        }
        semesters, _ := s.repo.ListProgramSemesters(programID)
        curriculum, _ := s.repo.GetProgramCurriculum(programID)
        return map[string]interface{}{
                "program":    program,
                "semesters":  semesters,
                "curriculum": curriculum,
        }, nil
}
func (s *Service) AddSubjectToProgram(ctx context.Context, ps *domain.ProgramSubject) error {
        return s.repo.AddSubjectToProgram(ps)
}
func (s *Service) RemoveSubjectFromProgram(ctx context.Context, programID, subjectID uint) error {
        return s.repo.RemoveSubjectFromProgram(programID, subjectID)
}

// Subjects
func (s *Service) ListSubjects(ctx context.Context, departmentID uint) ([]domain.Subject, error) {
        return s.repo.ListSubjects(departmentID)
}
func (s *Service) GetSubject(ctx context.Context, id uint) (*domain.Subject, error) {
        return s.repo.GetSubject(id)
}
func (s *Service) CreateSubject(ctx context.Context, sub *domain.Subject) error {
        if sub.SubjectName == "" || sub.SubjectCode == "" {
                return apperrors.BadRequest("subject name and code are required")
        }
        return s.repo.CreateSubject(sub)
}
func (s *Service) UpdateSubject(ctx context.Context, id uint, sub *domain.Subject) error {
        existing, err := s.repo.GetSubject(id)
        if err != nil {
                return apperrors.NotFound("subject not found")
        }
        sub.ID = existing.ID
        return s.repo.UpdateSubject(sub)
}

// Batches & Sections
func (s *Service) ListBatches(ctx context.Context, programID uint) ([]domain.Batch, error) {
        return s.repo.ListBatches(programID)
}
func (s *Service) GetBatch(ctx context.Context, id uint) (*domain.Batch, error) {
        return s.repo.GetBatch(id)
}
func (s *Service) CreateBatch(ctx context.Context, b *domain.Batch) error {
        if b.ProgramID == 0 {
                return apperrors.BadRequest("program_id is required")
        }
        return s.repo.CreateBatch(b)
}
func (s *Service) ListSections(ctx context.Context, batchID uint) ([]domain.Section, error) {
        return s.repo.ListSections(batchID)
}
func (s *Service) CreateSection(ctx context.Context, sec *domain.Section) error {
        return s.repo.CreateSection(sec)
}

// Course Offerings
func (s *Service) ListOfferings(ctx context.Context, termID, programID uint) ([]domain.CourseOffering, error) {
        return s.repo.ListOfferings(termID, programID)
}
func (s *Service) GetOffering(ctx context.Context, id uint) (*domain.CourseOffering, error) {
        return s.repo.GetOffering(id)
}
func (s *Service) CreateOffering(ctx context.Context, o *domain.CourseOffering) error {
        if o.SubjectID == 0 || o.AcademicTermID == 0 {
                return apperrors.BadRequest("subject_id and academic_term_id are required")
        }
        return s.repo.CreateOffering(o)
}
func (s *Service) UpdateOffering(ctx context.Context, id uint, o *domain.CourseOffering) error {
        existing, err := s.repo.GetOffering(id)
        if err != nil {
                return apperrors.NotFound("offering not found")
        }
        o.ID = existing.ID
        return s.repo.UpdateOffering(o)
}

// Term Registrations
func (s *Service) RegisterStudentForTerm(ctx context.Context, tr *domain.TermRegistration) error {
        if tr.StudentID == 0 || tr.AcademicTermID == 0 {
                return apperrors.BadRequest("student_id and term_id are required")
        }
        tr.RegistrationDate = time.Now()
        tr.Status = "Active"
        return s.repo.CreateTermRegistration(tr)
}
func (s *Service) ListTermRegistrations(ctx context.Context, termID uint) ([]domain.TermRegistration, error) {
        return s.repo.ListTermRegistrations(termID)
}

// Course Registrations
func (s *Service) RegisterStudentForCourse(ctx context.Context, cr *domain.CourseRegistration) error {
        if cr.StudentID == 0 || cr.OfferingID == 0 {
                return apperrors.BadRequest("student_id and offering_id are required")
        }
        cr.RegistrationStatus = "Enrolled"
        return s.repo.CreateCourseRegistration(cr)
}
func (s *Service) ListStudentCourses(ctx context.Context, studentID, termID uint) ([]domain.CourseRegistration, error) {
        return s.repo.ListStudentCourseRegistrations(studentID, termID)
}
func (s *Service) DropCourse(ctx context.Context, id uint) error {
        return s.repo.DropCourseRegistration(id)
}

// Timetable
func (s *Service) GetOfferingTimetable(ctx context.Context, offeringID uint) ([]domain.Timetable, error) {
        return s.repo.GetTimetable(offeringID)
}
func (s *Service) GetStudentTimetable(ctx context.Context, studentID, termID uint) ([]domain.Timetable, error) {
        return s.repo.GetStudentTimetable(studentID, termID)
}
func (s *Service) CreateTimetableEntry(ctx context.Context, t *domain.Timetable) error {
        return s.repo.CreateTimetableEntry(t)
}

// Calendar
func (s *Service) ListCalendar(ctx context.Context, campusID uint) ([]domain.AcademicCalendar, error) {
        return s.repo.ListCalendar(campusID)
}
func (s *Service) CreateCalendarEvent(ctx context.Context, e *domain.AcademicCalendar) error {
        return s.repo.CreateCalendarEvent(e)
}

// Attendance
func (s *Service) CreateClassSession(ctx context.Context, cs *domain.ClassSession) error {
        return s.repo.CreateClassSession(cs)
}
func (s *Service) GetClassSessions(ctx context.Context, offeringID uint) ([]domain.ClassSession, error) {
        return s.repo.GetClassSessions(offeringID)
}
func (s *Service) MarkAttendance(ctx context.Context, a *domain.StudentAttendance) error {
        return s.repo.MarkAttendance(a)
}
func (s *Service) GetStudentAttendance(ctx context.Context, studentID, offeringID uint) ([]domain.StudentAttendance, error) {
        return s.repo.GetStudentAttendance(studentID, offeringID)
}
func (s *Service) GetAttendanceSummary(ctx context.Context, studentID, termID uint) ([]AttendanceSummary, error) {
        return s.repo.GetAttendanceSummary(studentID, termID)
}
