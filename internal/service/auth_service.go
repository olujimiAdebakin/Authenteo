package service



import (
	"context"
	"errors"
	"time"

	"authentio/internal/models"
	"authentio/internal/pkg/response"
	"authentio/internal/repository"
	"authentio/pkg/jwt"
	"authentio/pkg/logger"
	"authentio/pkg/password"
)

type AuthService struct {
	userRepo repository.UserRepository
	twoFARepo repository.TwoFARepository
	jwtManager *jwt.Manager
}

// NewAuthService constructs the AuthService with its dependencies.
func NewAuthService(
	userRepo repository.UserRepository,
	twoFARepo repository.TwoFARepository,
	jwtManager *jwt.Manager,
) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		twoFARepo:  twoFARepo,
		jwtManager: jwtManager,
	}
}

type RegisterResponse struct {
	User    response.UserResponse `json:"user"`
	Message string                `json:"message"`
}

// LoginResponse is returned after successful login
type LoginResponse struct {
	User  response.UserResponse `json:"user"`
	Token string               `json:"token"`
}

// Register handles user registration flow.
func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest) error {
	existingUser, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		return errors.New("email already exists")
	}

	hashed, err := password.Hash(req.Password)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:     req.Email,
		Password:  hashed,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}


	if err := s.userRepo.Create(ctx, user); err != nil {
		return err
	}

		// Convert to response DTO
	userResponse := response.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		IsActive: user.IsActive,
	}

	// logger.Info("user registered", "email", req.Email)

	logger.Info("user registered", "email", req.Email)
	return nil
}

// Login validates credentials and returns a signed JWT.
func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return "", errors.New("invalid email or password")
	}

	if !password.Check(req.Password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	token, err := s.jwtManager.Generate(user.ID, user.Email)
	if err != nil {
		return "", err
	}

	logger.Info("user logged in", "email", req.Email)
	return token, nil
}

// Verify2FA checks OTP validity and activates user 2FA.
func (s *AuthService) Verify2FA(ctx context.Context, email, code string) error {
	valid, err := s.twoFARepo.VerifyCode(ctx, email, code)
	if err != nil || !valid {
		return errors.New("invalid or expired code")
	}
	return nil
}
