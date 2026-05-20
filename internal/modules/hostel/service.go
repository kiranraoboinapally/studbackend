package hostelmod

import (
	"context"
	"fmt"
	"time"

	"university-erp-backend/internal/domain"
	"university-erp-backend/internal/platform/apperrors"
	"university-erp-backend/internal/platform/eventbus"
	"university-erp-backend/internal/platform/outbox"

	"gorm.io/gorm"
)

type Service struct {
	repo   *Repository
	bus    *eventbus.Bus
	outbox *outbox.Writer
	db     *gorm.DB
}

func NewService(repo *Repository, bus *eventbus.Bus, ob *outbox.Writer, db *gorm.DB) *Service {
	return &Service{repo: repo, bus: bus, outbox: ob, db: db}
}

// Hostels
func (s *Service) ListHostels(ctx context.Context, campusID uint) ([]domain.Hostel, error) {
	return s.repo.ListHostels(campusID)
}
func (s *Service) GetHostel(ctx context.Context, id uint) (*domain.Hostel, error) {
	h, err := s.repo.GetHostel(id)
	if err != nil {
		return nil, apperrors.NotFound("hostel not found")
	}
	return h, nil
}
func (s *Service) CreateHostel(ctx context.Context, h *domain.Hostel) error {
	if h.Name == "" || h.Code == "" {
		return apperrors.BadRequest("hostel name and code are required")
	}
	return s.repo.CreateHostel(h)
}
func (s *Service) UpdateHostel(ctx context.Context, id uint, h *domain.Hostel) error {
	existing, err := s.repo.GetHostel(id)
	if err != nil {
		return apperrors.NotFound("hostel not found")
	}
	h.ID = existing.ID
	return s.repo.UpdateHostel(h)
}

// Rooms
func (s *Service) ListRooms(ctx context.Context, hostelID uint) ([]domain.HostelRoom, error) {
	return s.repo.ListRooms(hostelID)
}
func (s *Service) GetRoom(ctx context.Context, id uint) (*domain.HostelRoom, error) {
	rm, err := s.repo.GetRoom(id)
	if err != nil {
		return nil, apperrors.NotFound("room not found")
	}
	return rm, nil
}
func (s *Service) CreateRoom(ctx context.Context, rm *domain.HostelRoom) error {
	if rm.HostelID == 0 {
		return apperrors.BadRequest("hostel_id is required")
	}
	return s.repo.CreateRoom(rm)
}
func (s *Service) UpdateRoom(ctx context.Context, id uint, rm *domain.HostelRoom) error {
	existing, err := s.repo.GetRoom(id)
	if err != nil {
		return apperrors.NotFound("room not found")
	}
	rm.ID = existing.ID
	return s.repo.UpdateRoom(rm)
}

// Beds
func (s *Service) ListBeds(ctx context.Context, roomID uint) ([]domain.HostelBed, error) {
	return s.repo.ListBeds(roomID)
}
func (s *Service) CreateBed(ctx context.Context, b *domain.HostelBed) error {
	if b.RoomID == 0 {
		return apperrors.BadRequest("room_id is required")
	}
	return s.repo.CreateBed(b)
}

// Allocations - with outbox event
func (s *Service) AllocateRoom(ctx context.Context, studentID, roomID uint, bedID *uint, createdBy uint) (*domain.HostelAllocation, error) {
	if studentID == 0 || roomID == 0 {
		return nil, apperrors.BadRequest("student_id and room_id are required")
	}

	// Check room capacity
	room, err := s.repo.GetRoom(roomID)
	if err != nil {
		return nil, apperrors.NotFound("room not found")
	}
	if room.CurrentOccupancy >= room.Capacity {
		return nil, apperrors.Conflict("room is at full capacity")
	}

	var activeStatus domain.StatusCode
	s.db.Where("module = ? AND code = ?", "hostel", "ACTIVE").First(&activeStatus)

	allocation := &domain.HostelAllocation{
		StudentID:     studentID,
		RoomID:        roomID,
		BedID:         bedID,
		AllocatedFrom: time.Now(),
		StatusID:      &activeStatus.ID,
		CreatedBy:     &createdBy,
	}

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(allocation).Error; err != nil {
			return err
		}
		// Update room occupancy
		tx.Model(&domain.HostelRoom{}).Where("id = ?", roomID).
			Updates(map[string]interface{}{
				"current_occupancy": gorm.Expr("current_occupancy + 1"),
				"is_available":      false,
			})
		// Update bed occupancy
		if bedID != nil {
			tx.Model(&domain.HostelBed{}).Where("id = ?", *bedID).Update("is_occupied", true)
		}
		return s.outbox.WriteEvent(tx, "HostelAllocation", fmt.Sprintf("%d", allocation.ID),
			eventbus.EventHostelAllocated,
			map[string]interface{}{
				"allocation_id": allocation.ID,
				"student_id":   studentID,
				"room_id":      roomID,
				"hostel_id":     room.HostelID,
			},
		)
	})
	if txErr != nil {
		return nil, txErr
	}
	return allocation, nil
}

func (s *Service) DeallocateRoom(ctx context.Context, allocationID uint) error {
	allocation, err := s.repo.GetAllocation(allocationID)
	if err != nil {
		return apperrors.NotFound("allocation not found")
	}
	if allocation.AllocatedTo != nil {
		return apperrors.Conflict("allocation already ended")
	}

	now := time.Now()
	allocation.AllocatedTo = &now

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(allocation).Error; err != nil {
			return err
		}
		// Update room occupancy
		room, _ := s.repo.GetRoom(allocation.RoomID)
		newOccupancy := room.CurrentOccupancy - 1
		tx.Model(&domain.HostelRoom{}).Where("id = ?", allocation.RoomID).
			Updates(map[string]interface{}{
				"current_occupancy": newOccupancy,
				"is_available":      newOccupancy < room.Capacity,
			})
		// Update bed
		if allocation.BedID != nil {
			tx.Model(&domain.HostelBed{}).Where("id = ?", *allocation.BedID).Update("is_occupied", false)
		}
		// Create history
		hist := &domain.HostelAllocationHistory{
			StudentID:     allocation.StudentID,
			RoomID:        allocation.RoomID,
			AllocatedFrom: allocation.AllocatedFrom,
			AllocatedTo:   now,
			Reason:        "Deallocation",
		}
		return tx.Create(hist).Error
	})
}

func (s *Service) ListAllocations(ctx context.Context, hostelID, studentID uint) ([]domain.HostelAllocation, error) {
	return s.repo.ListAllocations(hostelID, studentID)
}

// Mess Bills
func (s *Service) CreateMessBill(ctx context.Context, m *domain.MessBill) error {
	if m.StudentID == 0 || m.Amount <= 0 {
		return apperrors.BadRequest("student_id and amount are required")
	}
	return s.repo.CreateMessBill(m)
}
func (s *Service) ListMessBills(ctx context.Context, studentID uint) ([]domain.MessBill, error) {
	return s.repo.ListMessBills(studentID)
}
func (s *Service) PayMessBill(ctx context.Context, id uint) error {
	var m domain.MessBill
	if err := s.db.First(&m, id).Error; err != nil {
		return apperrors.NotFound("mess bill not found")
	}
	now := time.Now()
	m.Paid = true
	m.PaidAt = &now
	return s.repo.UpdateMessBill(&m)
}

// Maintenance Requests - with outbox event
func (s *Service) CreateMaintenanceRequest(ctx context.Context, m *domain.MaintenanceRequest) error {
	if m.StudentID == 0 || m.RoomID == 0 {
		return apperrors.BadRequest("student_id and room_id are required")
	}

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(m).Error; err != nil {
			return err
		}
		return s.outbox.WriteEvent(tx, "MaintenanceRequest", fmt.Sprintf("%d", m.ID),
			eventbus.EventMaintenanceRequested,
			map[string]interface{}{
				"request_id":  m.ID,
				"student_id":  m.StudentID,
				"room_id":     m.RoomID,
				"category":    m.Category,
				"description": m.Description,
			},
		)
	})
	return txErr
}
func (s *Service) ListMaintenanceRequests(ctx context.Context, roomID, studentID uint) ([]domain.MaintenanceRequest, error) {
	return s.repo.ListMaintenanceRequests(roomID, studentID)
}
func (s *Service) ResolveMaintenanceRequest(ctx context.Context, id, assignedTo uint) error {
	var m domain.MaintenanceRequest
	if err := s.db.First(&m, id).Error; err != nil {
		return apperrors.NotFound("maintenance request not found")
	}
	now := time.Now()
	m.AssignedTo = &assignedTo
	m.ResolvedAt = &now
	return s.repo.UpdateMaintenanceRequest(&m)
}

// Visitor Logs
func (s *Service) CreateVisitorLog(ctx context.Context, v *domain.VisitorLog) error {
	if v.HostelID == 0 || v.VisitorName == "" {
		return apperrors.BadRequest("hostel_id and visitor_name are required")
	}
	return s.repo.CreateVisitorLog(v)
}
func (s *Service) ListVisitorLogs(ctx context.Context, hostelID uint) ([]domain.VisitorLog, error) {
	return s.repo.ListVisitorLogs(hostelID)
}
func (s *Service) ExitVisitor(ctx context.Context, id uint) error {
	var v domain.VisitorLog
	if err := s.db.First(&v, id).Error; err != nil {
		return apperrors.NotFound("visitor log not found")
	}
	now := time.Now()
	v.ExitTime = &now
	return s.repo.UpdateVisitorLog(&v)
}

// Stats
func (s *Service) GetStats(ctx context.Context) (map[string]interface{}, error) {
	return s.repo.GetHostelStats()
}
