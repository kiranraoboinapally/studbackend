package coremod

import (
	"university-erp-backend/internal/domain"

	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

// University
func (r *Repository) ListUniversities() ([]domain.University, error) {
	var list []domain.University
	return list, r.db.Where("is_active = true").Order("name").Find(&list).Error
}
func (r *Repository) GetUniversity(id uint) (*domain.University, error) {
	var u domain.University
	return &u, r.db.First(&u, id).Error
}
func (r *Repository) CreateUniversity(u *domain.University) error {
	return r.db.Create(u).Error
}
func (r *Repository) UpdateUniversity(u *domain.University) error {
	return r.db.Save(u).Error
}

// Campus
func (r *Repository) ListCampuses(universityID uint) ([]domain.Campus, error) {
	var list []domain.Campus
	q := r.db.Where("is_active = true").Order("name")
	if universityID > 0 {
		q = q.Where("university_id = ?", universityID)
	}
	return list, q.Find(&list).Error
}
func (r *Repository) GetCampus(id uint) (*domain.Campus, error) {
	var c domain.Campus
	return &c, r.db.First(&c, id).Error
}
func (r *Repository) CreateCampus(c *domain.Campus) error {
	return r.db.Create(c).Error
}
func (r *Repository) UpdateCampus(c *domain.Campus) error {
	return r.db.Save(c).Error
}

// Department
func (r *Repository) ListDepartments(campusID uint) ([]domain.Department, error) {
	var list []domain.Department
	q := r.db.Where("is_active = true").Order("name")
	if campusID > 0 {
		q = q.Where("campus_id = ?", campusID)
	}
	return list, q.Find(&list).Error
}
func (r *Repository) GetDepartment(id uint) (*domain.Department, error) {
	var d domain.Department
	return &d, r.db.First(&d, id).Error
}
func (r *Repository) CreateDepartment(d *domain.Department) error {
	return r.db.Create(d).Error
}
func (r *Repository) UpdateDepartment(d *domain.Department) error {
	return r.db.Save(d).Error
}

// Rooms
func (r *Repository) ListRooms(campusID uint) ([]domain.Room, error) {
	var list []domain.Room
	q := r.db.Where("is_active = true").Order("room_number")
	if campusID > 0 {
		q = q.Where("campus_id = ?", campusID)
	}
	return list, q.Find(&list).Error
}
func (r *Repository) CreateRoom(rm *domain.Room) error {
	return r.db.Create(rm).Error
}

// System lookups
func (r *Repository) ListGenders() ([]domain.Gender, error) {
	var list []domain.Gender
	return list, r.db.Find(&list).Error
}
func (r *Repository) ListCategories() ([]domain.Category, error) {
	var list []domain.Category
	return list, r.db.Find(&list).Error
}
func (r *Repository) ListBloodGroups() ([]domain.BloodGroup, error) {
	var list []domain.BloodGroup
	return list, r.db.Find(&list).Error
}
func (r *Repository) ListStatusCodes(module string) ([]domain.StatusCode, error) {
	var list []domain.StatusCode
	q := r.db.Where("is_active = true")
	if module != "" {
		q = q.Where("module = ?", module)
	}
	return list, q.Find(&list).Error
}
