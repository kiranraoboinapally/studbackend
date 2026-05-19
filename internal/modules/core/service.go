package coremod

import (
	"context"

	"university-erp-backend/internal/domain"
	"university-erp-backend/internal/platform/apperrors"
)

type Service struct{ repo *Repository }

func NewService(repo *Repository) *Service { return &Service{repo: repo} }

func (s *Service) ListUniversities(_ context.Context) ([]domain.University, error) {
	return s.repo.ListUniversities()
}
func (s *Service) GetUniversity(_ context.Context, id uint) (*domain.University, error) {
	return s.repo.GetUniversity(id)
}
func (s *Service) CreateUniversity(_ context.Context, u *domain.University) error {
	if u.Name == "" {
		return apperrors.BadRequest("university name is required")
	}
	return s.repo.CreateUniversity(u)
}
func (s *Service) UpdateUniversity(_ context.Context, id uint, u *domain.University) error {
	existing, err := s.repo.GetUniversity(id)
	if err != nil {
		return apperrors.NotFound("university not found")
	}
	u.ID = existing.ID
	return s.repo.UpdateUniversity(u)
}

func (s *Service) ListCampuses(_ context.Context, universityID uint) ([]domain.Campus, error) {
	return s.repo.ListCampuses(universityID)
}
func (s *Service) GetCampus(_ context.Context, id uint) (*domain.Campus, error) {
	return s.repo.GetCampus(id)
}
func (s *Service) CreateCampus(_ context.Context, c *domain.Campus) error {
	if c.Name == "" {
		return apperrors.BadRequest("campus name is required")
	}
	return s.repo.CreateCampus(c)
}
func (s *Service) UpdateCampus(_ context.Context, id uint, c *domain.Campus) error {
	existing, err := s.repo.GetCampus(id)
	if err != nil {
		return apperrors.NotFound("campus not found")
	}
	c.ID = existing.ID
	return s.repo.UpdateCampus(c)
}

func (s *Service) ListDepartments(_ context.Context, campusID uint) ([]domain.Department, error) {
	return s.repo.ListDepartments(campusID)
}
func (s *Service) GetDepartment(_ context.Context, id uint) (*domain.Department, error) {
	return s.repo.GetDepartment(id)
}
func (s *Service) CreateDepartment(_ context.Context, d *domain.Department) error {
	if d.Name == "" {
		return apperrors.BadRequest("department name is required")
	}
	return s.repo.CreateDepartment(d)
}
func (s *Service) UpdateDepartment(_ context.Context, id uint, d *domain.Department) error {
	existing, err := s.repo.GetDepartment(id)
	if err != nil {
		return apperrors.NotFound("department not found")
	}
	d.ID = existing.ID
	return s.repo.UpdateDepartment(d)
}

func (s *Service) ListRooms(_ context.Context, campusID uint) ([]domain.Room, error) {
	return s.repo.ListRooms(campusID)
}
func (s *Service) CreateRoom(_ context.Context, rm *domain.Room) error {
	return s.repo.CreateRoom(rm)
}

func (s *Service) ListGenders(_ context.Context) ([]domain.Gender, error) {
	return s.repo.ListGenders()
}
func (s *Service) ListCategories(_ context.Context) ([]domain.Category, error) {
	return s.repo.ListCategories()
}
func (s *Service) ListBloodGroups(_ context.Context) ([]domain.BloodGroup, error) {
	return s.repo.ListBloodGroups()
}
func (s *Service) ListStatusCodes(_ context.Context, module string) ([]domain.StatusCode, error) {
	return s.repo.ListStatusCodes(module)
}
