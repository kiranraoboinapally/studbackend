package hostelmod

import (
	"university-erp-backend/internal/domain"

	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

// Hostels
func (r *Repository) ListHostels(campusID uint) ([]domain.Hostel, error) {
	var list []domain.Hostel
	q := r.db.Where("is_active = true").Order("name")
	if campusID > 0 {
		q = q.Where("campus_id = ?", campusID)
	}
	return list, q.Find(&list).Error
}
func (r *Repository) GetHostel(id uint) (*domain.Hostel, error) {
	var h domain.Hostel
	return &h, r.db.First(&h, id).Error
}
func (r *Repository) CreateHostel(h *domain.Hostel) error {
	return r.db.Create(h).Error
}
func (r *Repository) UpdateHostel(h *domain.Hostel) error {
	return r.db.Save(h).Error
}

// Rooms
func (r *Repository) ListRooms(hostelID uint) ([]domain.HostelRoom, error) {
	var list []domain.HostelRoom
	q := r.db.Order("room_number")
	if hostelID > 0 {
		q = q.Where("hostel_id = ?", hostelID)
	}
	return list, q.Find(&list).Error
}
func (r *Repository) GetRoom(id uint) (*domain.HostelRoom, error) {
	var rm domain.HostelRoom
	return &rm, r.db.First(&rm, id).Error
}
func (r *Repository) CreateRoom(rm *domain.HostelRoom) error {
	return r.db.Create(rm).Error
}
func (r *Repository) UpdateRoom(rm *domain.HostelRoom) error {
	return r.db.Save(rm).Error
}

// Beds
func (r *Repository) ListBeds(roomID uint) ([]domain.HostelBed, error) {
	var list []domain.HostelBed
	q := r.db.Order("bed_number")
	if roomID > 0 {
		q = q.Where("room_id = ?", roomID)
	}
	return list, q.Find(&list).Error
}
func (r *Repository) CreateBed(b *domain.HostelBed) error {
	return r.db.Create(b).Error
}
func (r *Repository) UpdateBed(b *domain.HostelBed) error {
	return r.db.Save(b).Error
}

// Allocations
func (r *Repository) CreateAllocation(a *domain.HostelAllocation) error {
	return r.db.Create(a).Error
}
func (r *Repository) GetAllocation(id uint) (*domain.HostelAllocation, error) {
	var a domain.HostelAllocation
	return &a, r.db.First(&a, id).Error
}
func (r *Repository) GetActiveAllocation(studentID uint) (*domain.HostelAllocation, error) {
	var a domain.HostelAllocation
	return &a, r.db.Where("student_id = ? AND allocated_to IS NULL", studentID).First(&a).Error
}
func (r *Repository) ListAllocations(hostelID uint, studentID uint) ([]domain.HostelAllocation, error) {
	var list []domain.HostelAllocation
	q := r.db.Order("allocated_from DESC")
	if hostelID > 0 {
		q = q.Where("room_id IN (SELECT id FROM hostel.rooms WHERE hostel_id = ?)", hostelID)
	}
	if studentID > 0 {
		q = q.Where("student_id = ?", studentID)
	}
	return list, q.Find(&list).Error
}
func (r *Repository) UpdateAllocation(a *domain.HostelAllocation) error {
	return r.db.Save(a).Error
}
func (r *Repository) CreateAllocationHistory(h *domain.HostelAllocationHistory) error {
	return r.db.Create(h).Error
}

// Mess Bills
func (r *Repository) CreateMessBill(m *domain.MessBill) error {
	return r.db.Create(m).Error
}
func (r *Repository) ListMessBills(studentID uint) ([]domain.MessBill, error) {
	var list []domain.MessBill
	q := r.db.Order("month DESC")
	if studentID > 0 {
		q = q.Where("student_id = ?", studentID)
	}
	return list, q.Find(&list).Error
}
func (r *Repository) UpdateMessBill(m *domain.MessBill) error {
	return r.db.Save(m).Error
}

// Maintenance Requests
func (r *Repository) CreateMaintenanceRequest(m *domain.MaintenanceRequest) error {
	return r.db.Create(m).Error
}
func (r *Repository) ListMaintenanceRequests(roomID uint, studentID uint) ([]domain.MaintenanceRequest, error) {
	var list []domain.MaintenanceRequest
	q := r.db.Order("created_at DESC")
	if roomID > 0 {
		q = q.Where("room_id = ?", roomID)
	}
	if studentID > 0 {
		q = q.Where("student_id = ?", studentID)
	}
	return list, q.Find(&list).Error
}
func (r *Repository) UpdateMaintenanceRequest(m *domain.MaintenanceRequest) error {
	return r.db.Save(m).Error
}

// Visitor Logs
func (r *Repository) CreateVisitorLog(v *domain.VisitorLog) error {
	return r.db.Create(v).Error
}
func (r *Repository) ListVisitorLogs(hostelID uint) ([]domain.VisitorLog, error) {
	var list []domain.VisitorLog
	q := r.db.Order("entry_time DESC")
	if hostelID > 0 {
		q = q.Where("hostel_id = ?", hostelID)
	}
	return list, q.Find(&list).Error
}
func (r *Repository) UpdateVisitorLog(v *domain.VisitorLog) error {
	return r.db.Save(v).Error
}

// Stats
func (r *Repository) GetHostelStats() (map[string]interface{}, error) {
	var totalHostels, totalRooms, totalBeds, occupiedBeds int64
	r.db.Model(&domain.Hostel{}).Where("is_active = true").Count(&totalHostels)
	r.db.Model(&domain.HostelRoom{}).Count(&totalRooms)
	r.db.Model(&domain.HostelBed{}).Count(&totalBeds)
	r.db.Model(&domain.HostelBed{}).Where("is_occupied = true").Count(&occupiedBeds)
	return map[string]interface{}{
		"total_hostels":  totalHostels,
		"total_rooms":    totalRooms,
		"total_beds":     totalBeds,
		"occupied_beds":  occupiedBeds,
		"available_beds": totalBeds - occupiedBeds,
	}, nil
}
