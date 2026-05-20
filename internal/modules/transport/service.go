package transportmod

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

// Buses
func (s *Service) ListBuses(ctx context.Context, activeOnly bool) ([]domain.Bus, error) {
	return s.repo.ListBuses(activeOnly)
}
func (s *Service) GetBus(ctx context.Context, id uint) (*domain.Bus, error) {
	b, err := s.repo.GetBus(id)
	if err != nil {
		return nil, apperrors.NotFound("bus not found")
	}
	return b, nil
}
func (s *Service) CreateBus(ctx context.Context, b *domain.Bus) error {
	if b.BusNumber == "" || b.RegistrationNo == "" {
		return apperrors.BadRequest("bus_number and registration_no are required")
	}
	return s.repo.CreateBus(b)
}
func (s *Service) UpdateBus(ctx context.Context, id uint, b *domain.Bus) error {
	existing, err := s.repo.GetBus(id)
	if err != nil {
		return apperrors.NotFound("bus not found")
	}
	b.ID = existing.ID
	return s.repo.UpdateBus(b)
}

// Routes
func (s *Service) ListRoutes(ctx context.Context, activeOnly bool) ([]domain.Route, error) {
	return s.repo.ListRoutes(activeOnly)
}
func (s *Service) GetRoute(ctx context.Context, id uint) (*domain.Route, error) {
	rt, err := s.repo.GetRoute(id)
	if err != nil {
		return nil, apperrors.NotFound("route not found")
	}
	return rt, nil
}
func (s *Service) CreateRoute(ctx context.Context, rt *domain.Route) error {
	if rt.RouteName == "" {
		return apperrors.BadRequest("route_name is required")
	}
	return s.repo.CreateRoute(rt)
}
func (s *Service) UpdateRoute(ctx context.Context, id uint, rt *domain.Route) error {
	existing, err := s.repo.GetRoute(id)
	if err != nil {
		return apperrors.NotFound("route not found")
	}
	rt.ID = existing.ID
	return s.repo.UpdateRoute(rt)
}

// Stops
func (s *Service) ListStops(ctx context.Context, routeID uint) ([]domain.Stop, error) {
	return s.repo.ListStops(routeID)
}
func (s *Service) CreateStop(ctx context.Context, st *domain.Stop) error {
	if st.RouteID == 0 || st.StopName == "" {
		return apperrors.BadRequest("route_id and stop_name are required")
	}
	return s.repo.CreateStop(st)
}
func (s *Service) UpdateStop(ctx context.Context, st *domain.Stop) error {
	return s.repo.UpdateStop(st)
}

// Bus Assignments
func (s *Service) ListBusAssignments(ctx context.Context, busID, routeID uint) ([]domain.BusAssignment, error) {
	return s.repo.ListBusAssignments(busID, routeID)
}
func (s *Service) AssignBusToRoute(ctx context.Context, ba *domain.BusAssignment) error {
	if ba.BusID == 0 || ba.RouteID == 0 {
		return apperrors.BadRequest("bus_id and route_id are required")
	}
	ba.EffectiveFrom = time.Now()
	ba.IsActive = true
	return s.repo.CreateBusAssignment(ba)
}
func (s *Service) EndBusAssignment(ctx context.Context, id uint) error {
	ba, err := s.repo.GetBusAssignment(id)
	if err != nil {
		return apperrors.NotFound("bus assignment not found")
	}
	now := time.Now()
	ba.EffectiveTo = &now
	ba.IsActive = false
	return s.repo.UpdateBusAssignment(ba)
}

// Student Passes - with outbox event
func (s *Service) IssuePass(ctx context.Context, sp *domain.StudentPass, createdBy uint) (*domain.StudentPass, error) {
	if sp.StudentID == 0 || sp.RouteID == 0 {
		return nil, apperrors.BadRequest("student_id and route_id are required")
	}

	var activeStatus domain.StatusCode
	s.db.Where("module = ? AND code = ?", "transport", "ACTIVE").First(&activeStatus)
	sp.StatusID = &activeStatus.ID

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(sp).Error; err != nil {
			return err
		}
		return s.outbox.WriteEvent(tx, "StudentPass", fmt.Sprintf("%d", sp.ID),
			eventbus.EventTransportPassIssued,
			map[string]interface{}{
				"pass_id":     sp.ID,
				"student_id":  sp.StudentID,
				"route_id":    sp.RouteID,
				"valid_from":  sp.ValidFrom,
				"valid_to":    sp.ValidTo,
				"fee_paid":    sp.FeePaid,
			},
		)
	})
	if txErr != nil {
		return nil, txErr
	}
	return sp, nil
}
func (s *Service) ListStudentPasses(ctx context.Context, studentID, routeID uint) ([]domain.StudentPass, error) {
	return s.repo.ListStudentPasses(studentID, routeID)
}
func (s *Service) RenewPass(ctx context.Context, id uint, validTo time.Time, feePaid float64) error {
	sp, err := s.repo.GetStudentPass(id)
	if err != nil {
		return apperrors.NotFound("student pass not found")
	}
	sp.ValidTo = validTo
	sp.FeePaid = sp.FeePaid + feePaid
	return s.repo.UpdateStudentPass(sp)
}

// Vehicle Maintenance
func (s *Service) ListVehicleMaintenance(ctx context.Context, busID uint) ([]domain.VehicleMaintenance, error) {
	return s.repo.ListVehicleMaintenance(busID)
}
func (s *Service) CreateVehicleMaintenance(ctx context.Context, vm *domain.VehicleMaintenance) error {
	if vm.BusID == 0 {
		return apperrors.BadRequest("bus_id is required")
	}
	vm.MaintenanceDate = time.Now()
	return s.repo.CreateVehicleMaintenance(vm)
}

// Stats
func (s *Service) GetStats(ctx context.Context) (map[string]interface{}, error) {
	return s.repo.GetTransportStats()
}
