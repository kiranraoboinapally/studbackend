package authmod

import (
	"context"
	"fmt"
	"time"

	"university-erp-backend/internal/domain"
	"university-erp-backend/internal/platform/apperrors"
	"university-erp-backend/internal/platform/auth"
	"university-erp-backend/internal/platform/eventbus"
	"university-erp-backend/internal/platform/outbox"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Service contains auth business logic.
type Service struct {
	repo   *Repository
	jwt    *auth.JWTManager
	bus    *eventbus.Bus
	outbox *outbox.Writer
	db     *gorm.DB
}

func NewService(repo *Repository, jwt *auth.JWTManager, bus *eventbus.Bus, ob *outbox.Writer, db *gorm.DB) *Service {
	return &Service{repo: repo, jwt: jwt, bus: bus, outbox: ob, db: db}
}

// ─── Login ───────────────────────────────────────────────────────────────────

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token    string   `json:"token"`
	UserID   uint     `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

func (s *Service) Login(ctx context.Context, req LoginRequest, ip, ua string) (*LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, apperrors.BadRequest("username and password are required")
	}

	user, err := s.repo.FindUserByUsername(req.Username)
	if err != nil {
		// Record failed attempt
		s.repo.RecordLoginAttempt(&domain.LoginAttempt{
			Username:      req.Username,
			Success:       false,
			IPAddress:     ip,
			UserAgent:     ua,
			FailureReason: "user not found",
			AttemptedAt:   time.Now(),
		})
		return nil, apperrors.Unauthorized("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		s.repo.RecordLoginAttempt(&domain.LoginAttempt{
			UserID:        &user.ID,
			Username:      req.Username,
			Success:       false,
			IPAddress:     ip,
			UserAgent:     ua,
			FailureReason: "wrong password",
			AttemptedAt:   time.Now(),
		})
		return nil, apperrors.Unauthorized("invalid credentials")
	}

	roles, err := s.repo.GetUserRoles(user.ID)
	if err != nil {
		return nil, apperrors.Internal("failed to get roles", err)
	}

	token, err := s.jwt.GenerateToken(user.ID, user.Username, roles)
	if err != nil {
		return nil, apperrors.Internal("failed to generate token", err)
	}

	// Record successful login
	now := time.Now()
	s.repo.RecordLoginAttempt(&domain.LoginAttempt{
		UserID:      &user.ID,
		Username:    req.Username,
		Success:     true,
		IPAddress:   ip,
		UserAgent:   ua,
		AttemptedAt: now,
	})
	s.repo.CreateSession(&domain.UserSession{
		UserID:       user.ID,
		SessionToken: token,
		IPAddress:    ip,
		UserAgent:    ua,
		LoginAt:      now,
		LastActivity: now,
		IsActive:     true,
	})
	s.db.Model(user).Update("last_login_at", &now)

	// Publish event
	s.bus.PublishAsync(ctx, eventbus.Event{
		Type:          eventbus.EventUserLoggedIn,
		AggregateType: "User",
		AggregateID:   fmt.Sprintf("%d", user.ID),
		Payload:       map[string]interface{}{"user_id": user.ID, "ip": ip},
	})

	return &LoginResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
		Roles:    roles,
	}, nil
}

// ─── Register ────────────────────────────────────────────────────────────────

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleName string `json:"role_name"` // e.g. "student", "faculty"
}

type RegisterResponse struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, apperrors.BadRequest("username, email, and password are required")
	}
	if len(req.Password) < 6 {
		return nil, apperrors.Validation("password must be at least 6 characters")
	}

	// Check duplicates
	if existing, _ := s.repo.FindUserByUsername(req.Username); existing != nil {
		return nil, apperrors.Conflict("username already exists")
	}
	if existing, _ := s.repo.FindUserByEmail(req.Email); existing != nil {
		return nil, apperrors.Conflict("email already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, apperrors.Internal("failed to hash password", err)
	}

	var resp *RegisterResponse

	// Transactional: create user + assign role + write outbox event
	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		user := domain.User{
			Username:     req.Username,
			Email:        req.Email,
			PasswordHash: string(hash),
			IsActive:     true,
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// Assign role
		roleName := req.RoleName
		if roleName == "" {
			roleName = "student"
		}
		var role domain.Role
		if err := tx.Where("role_name = ?", roleName).First(&role).Error; err != nil {
			return fmt.Errorf("role '%s' not found", roleName)
		}
		if err := tx.Create(&domain.UserRole{UserID: user.ID, RoleID: role.ID, AssignedAt: time.Now()}).Error; err != nil {
			return err
		}

		// Write outbox event (same transaction!)
		if err := s.outbox.WriteEvent(tx, "User", fmt.Sprintf("%d", user.ID),
			eventbus.EventUserRegistered,
			map[string]interface{}{
				"user_id":  user.ID,
				"username": user.Username,
				"email":    user.Email,
				"role":     roleName,
			},
		); err != nil {
			return err
		}

		resp = &RegisterResponse{
			UserID:   user.ID,
			Username: user.Username,
			Email:    user.Email,
		}
		return nil
	})

	if txErr != nil {
		return nil, apperrors.Internal("registration failed", txErr)
	}
	return resp, nil
}

// ─── Profile ─────────────────────────────────────────────────────────────────

type ProfileResponse struct {
	UserID   uint     `json:"user_id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	IsActive bool     `json:"is_active"`
}

func (s *Service) GetProfile(ctx context.Context, userID uint) (*ProfileResponse, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, apperrors.NotFound("user")
	}
	roles, _ := s.repo.GetUserRoles(userID)
	return &ProfileResponse{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Roles:    roles,
		IsActive: user.IsActive,
	}, nil
}
