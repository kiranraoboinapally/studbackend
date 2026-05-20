package transportmod

import (
	"university-erp-backend/internal/domain"

	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

// Buses
func (r *Repository) ListBuses(activeOnly bool) ([]domain.Bus, error) {
	var list []domain.Bus
	q := r.db.Order("bus_number")
	if activeOnly {
		q = q.Where("is_active = true")
	}
	return list, q.Find(&list).Error
}
func (r *Repository) GetBus(id uint) (*domain.Bus, error) {
	var b domain.Bus
	return &b, r.db.First(&b, id).Error
}
func (r *Repository) CreateBus(b *domain.Bus) error {
	return r.db.Create(b).Error
}
func (r *Repository) UpdateBus(b *domain.Bus) error {
	return r.db.Save(b).Error
}

// Routes
func (r *Repository) ListRoutes(activeOnly bool) ([]domain.Route, error) {
	var list []domain.Route
	q := r.db.Order("route_name")
	if activeOnly {
		q = q.Where("is_active = true")
	}
	return list, q.Find(&list).Error
}
func (r *Repository) GetRoute(id uint) (*domain.Route, error) {
	var rt domain.Route
	return &rt, r.db.First(&rt, id).Error
}
func (r *Repository) CreateRoute(rt *domain.Route) error {
	return r.db.Create(rt).Error
}
func (r *Repository) UpdateRoute(rt *domain.Route) error {
	return r.db.Save(rt).Error
}

// Stops
func (r *Repository) ListStops(routeID uint) ([]domain.Stop, error) {
	var list []domain.Stop
	q := r.db.Order("stop_order")
	if routeID > 0 {
		q = q.Where("route_id = ?", routeID)
	}
	return list, q.Find(&list).Error
}
func (r *Repository) CreateStop(s *domain.Stop) error {
	return r.db.Create(s).Error
}
func (r *Repository) UpdateStop(s *domain.Stop) error {
	return r.db.Save(s).Error
}

// Bus Assignments
func (r *Repository) ListBusAssignments(busID, routeID uint) ([]domain.BusAssignment, error) {
	var list []domain.BusAssignment
	q := r.db.Where("is_active = true").Order("effective_from DESC")
	if busID > 0 {
		q = q.Where("bus_id = ?", busID)
	}
	if routeID > 0 {
		q = q.Where("route_id = ?", routeID)
	}
	return list, q.Find(&list).Error
}
func (r *Repository) CreateBusAssignment(ba *domain.BusAssignment) error {
	return r.db.Create(ba).Error
}
func (r *Repository) GetBusAssignment(id uint) (*domain.BusAssignment, error) {
	var ba domain.BusAssignment
	return &ba, r.db.First(&ba, id).Error
}
func (r *Repository) UpdateBusAssignment(ba *domain.BusAssignment) error {
	return r.db.Save(ba).Error
}
// Student Passes
func (r *Repository) ListStudentPasses(studentID, routeID uint) ([]domain.StudentPass, error) {
	var list []domain.StudentPass
	q := r.db.Order("valid_from DESC")
	if studentID > 0 {
		q = q.Where("student_id = ?", studentID)
	}
	if routeID > 0 {
		q = q.Where("route_id = ?", routeID)
	}
	return list, q.Find(&list).Error
}
func (r *Repository) GetStudentPass(id uint) (*domain.StudentPass, error) {
	var sp domain.StudentPass
	return &sp, r.db.First(&sp, id).Error
}
func (r *Repository) CreateStudentPass(sp *domain.StudentPass) error {
	return r.db.Create(sp).Error
}
func (r *Repository) UpdateStudentPass(sp *domain.StudentPass) error {
	return r.db.Save(sp).Error
}

// Vehicle Maintenance
func (r *Repository) ListVehicleMaintenance(busID uint) ([]domain.VehicleMaintenance, error) {
	var list []domain.VehicleMaintenance
	q := r.db.Order("maintenance_date DESC")
	if busID > 0 {
		q = q.Where("bus_id = ?", busID)
	}
	return list, q.Find(&list).Error
}
func (r *Repository) CreateVehicleMaintenance(vm *domain.VehicleMaintenance) error {
	return r.db.Create(vm).Error
}
func (r *Repository) UpdateVehicleMaintenance(vm *domain.VehicleMaintenance) error {
	return r.db.Save(vm).Error
}

// Stats
func (r *Repository) GetTransportStats() (map[string]interface{}, error) {
	var totalBuses, activeBuses, totalRoutes, totalPasses int64
	r.db.Model(&domain.Bus{}).Count(&totalBuses)
	r.db.Model(&domain.Bus{}).Where("is_active = true").Count(&activeBuses)
	r.db.Model(&domain.Route{}).Where("is_active = true").Count(&totalRoutes)
	r.db.Model(&domain.StudentPass{}).Where("valid_to >= CURRENT_DATE").Count(&totalPasses)
	return map[string]interface{}{
		"total_buses":   totalBuses,
		"active_buses":  activeBuses,
		"total_routes":  totalRoutes,
		"active_passes": totalPasses,
	}, nil
}
