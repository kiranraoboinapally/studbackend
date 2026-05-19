package authmod

import (
	"university-erp-backend/internal/domain"

	"gorm.io/gorm"
)

// Repository handles all auth-related database operations.
type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ? AND is_active = ?", username, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ? AND is_active = ?", email, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindUserByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) CreateUser(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *Repository) UpdateUser(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *Repository) GetUserRoles(userID uint) ([]string, error) {
	var roleNames []string
	err := r.db.Table("shared.user_roles").
		Select("shared.roles.role_name").
		Joins("JOIN shared.roles ON shared.roles.id = shared.user_roles.role_id").
		Where("shared.user_roles.user_id = ?", userID).
		Scan(&roleNames).Error
	return roleNames, err
}

func (r *Repository) AssignRole(userID, roleID uint) error {
	ur := domain.UserRole{UserID: userID, RoleID: roleID}
	return r.db.FirstOrCreate(&ur, domain.UserRole{UserID: userID, RoleID: roleID}).Error
}

func (r *Repository) FindRoleByName(name string) (*domain.Role, error) {
	var role domain.Role
	err := r.db.Where("role_name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *Repository) RecordLoginAttempt(attempt *domain.LoginAttempt) error {
	return r.db.Create(attempt).Error
}

func (r *Repository) CreateSession(session *domain.UserSession) error {
	return r.db.Create(session).Error
}

func (r *Repository) InvalidateSession(token string) error {
	return r.db.Model(&domain.UserSession{}).
		Where("session_token = ?", token).
		Update("is_active", false).Error
}
