package admissionsmod

import (
        "university-erp-backend/internal/domain"

        "gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

// Admission Cycles
func (r *Repository) ListCycles() ([]domain.AdmissionCycle, error) {
        var list []domain.AdmissionCycle
        return list, r.db.Order("created_at DESC").Find(&list).Error
}
func (r *Repository) GetCycle(id uint) (*domain.AdmissionCycle, error) {
        var c domain.AdmissionCycle
        return &c, r.db.First(&c, id).Error
}
func (r *Repository) GetOpenCycles() ([]domain.AdmissionCycle, error) {
        var list []domain.AdmissionCycle
        return list, r.db.Where("is_open = true").Find(&list).Error
}
func (r *Repository) CreateCycle(c *domain.AdmissionCycle) error {
        return r.db.Create(c).Error
}
func (r *Repository) UpdateCycle(c *domain.AdmissionCycle) error {
        return r.db.Save(c).Error
}

// Applicants
func (r *Repository) ListApplicants(cycleID uint, page, pageSize int) ([]domain.Applicant, int64, error) {
        var list []domain.Applicant
        var total int64
        q := r.db.Model(&domain.Applicant{})
        if cycleID > 0 {
                q = q.Where("cycle_id = ?", cycleID)
        }
        q.Count(&total)
        if page < 1 { page = 1 }
        if pageSize < 1 { pageSize = 20 }
        offset := (page - 1) * pageSize
        return list, total, q.Offset(offset).Limit(pageSize).Order("applied_at DESC").Find(&list).Error
}
func (r *Repository) GetApplicant(id uint) (*domain.Applicant, error) {
        var a domain.Applicant
        return &a, r.db.First(&a, id).Error
}
func (r *Repository) GetApplicantByNumber(appNum string) (*domain.Applicant, error) {
        var a domain.Applicant
        return &a, r.db.Where("application_number = ?", appNum).First(&a).Error
}
func (r *Repository) CreateApplicant(a *domain.Applicant) error {
        return r.db.Create(a).Error
}
func (r *Repository) UpdateApplicant(a *domain.Applicant) error {
        return r.db.Save(a).Error
}
func (r *Repository) UpdateApplicantStatus(id uint, statusID uint) error {
        return r.db.Model(&domain.Applicant{}).Where("id = ?", id).Update("status_id", statusID).Error
}

// Documents
func (r *Repository) GetApplicantDocuments(applicantID uint) ([]domain.Document, error) {
        var list []domain.Document
        return list, r.db.Where("applicant_id = ?", applicantID).Find(&list).Error
}
func (r *Repository) CreateDocument(d *domain.Document) error {
        return r.db.Create(d).Error
}
func (r *Repository) VerifyDocument(id, verifiedBy uint) error {
        return r.db.Exec(`UPDATE admissions.documents SET verified_by = ?, verified_at = CURRENT_TIMESTAMP WHERE id = ?`, verifiedBy, id).Error
}

// Status History
func (r *Repository) CreateStatusHistory(h *domain.ApplicationStatusHistory) error {
        return r.db.Create(h).Error
}
func (r *Repository) GetStatusHistory(applicantID uint) ([]domain.ApplicationStatusHistory, error) {
        var list []domain.ApplicationStatusHistory
        return list, r.db.Where("applicant_id = ?", applicantID).Order("effective_from DESC").Find(&list).Error
}

// Seat Allocation
func (r *Repository) CreateSeatAllocation(sa *domain.SeatAllocation) error {
        return r.db.Create(sa).Error
}
func (r *Repository) GetSeatAllocation(applicantID uint) (*domain.SeatAllocation, error) {
        var sa domain.SeatAllocation
        return &sa, r.db.Where("applicant_id = ?", applicantID).First(&sa).Error
}
func (r *Repository) ListSeatAllocations(cycleID uint) ([]domain.SeatAllocation, error) {
        var list []domain.SeatAllocation
        return list, r.db.Where("cycle_id = ?", cycleID).Order("allocation_rank").Find(&list).Error
}

// Waitlist
func (r *Repository) AddToWaitlist(w *domain.Waitlist) error {
        return r.db.Create(w).Error
}
func (r *Repository) GetWaitlist(cycleID uint) ([]domain.Waitlist, error) {
        var list []domain.Waitlist
        return list, r.db.Where("cycle_id = ?", cycleID).Order("rank").Find(&list).Error
}

// Applicant-Student Mapping
func (r *Repository) CreateApplicantStudentMap(m *domain.ApplicantStudentMap) error {
        return r.db.Create(m).Error
}
func (r *Repository) GetApplicantStudentMap(applicantID uint) (*domain.ApplicantStudentMap, error) {
        var m domain.ApplicantStudentMap
        return &m, r.db.Where("applicant_id = ?", applicantID).First(&m).Error
}

// Stats
func (r *Repository) GetCycleStats(cycleID uint) (map[string]int64, error) {
        stats := make(map[string]int64)
        var total, allocated, waitlisted int64
        r.db.Model(&domain.Applicant{}).Where("cycle_id = ?", cycleID).Count(&total)
        r.db.Raw(`SELECT COUNT(*) FROM admissions.seat_allocations WHERE cycle_id = ?`, cycleID).Scan(&allocated)
        r.db.Raw(`SELECT COUNT(*) FROM admissions.waitlist WHERE cycle_id = ?`, cycleID).Scan(&waitlisted)
        stats["total"] = total
        stats["allocated"] = allocated
        stats["waitlisted"] = waitlisted
        return stats, nil
}
func (r *Repository) CountApplicationNumber(cycleID uint) (int64, error) {
        var count int64
        r.db.Model(&domain.Applicant{}).Where("cycle_id = ?", cycleID).Count(&count)
        return count, nil
}
