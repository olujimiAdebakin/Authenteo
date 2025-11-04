package service

import (
	"context"
	"errors"
	"time"

	"authentio/internal/constants"
	"authentio/internal/models"
	"authentio/internal/repository"
	"authentio/pkg/jwt"
	"authentio/pkg/logger"
	"authentio/pkg/password"
	"authentio/pkg/response"
)

type AuthService struct {
	userRepo repository.UserRepository
	twoFARepo repository.TwoFARepository
	otpRepo repository.OTPRepository
	jwtManager *jwt.Manager
}

// NewAuthService constructs the AuthService with its dependencies.
func NewAuthService(
	userRepo repository.UserRepository,
	twoFARepo repository.TwoFARepository,
	otpRepo repository.OTPRepository,
	jwtManager *jwt.Manager,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		otpRepo: otpRepo,
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
func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest)  (*RegisterResponse, error) {
	existingUser, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	hashed, err := password.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		FirstName: req.FirstName,  // Add firstName
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashed,
		IsActive:  true,
		BaseModel: models.BaseModel{
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
    },
	}


	if err := s.userRepo.Create(ctx, user); err != nil {
		 return nil, err
	}

		// Convert to response DTO
	userResponse := response.UserResponse{
		ID:       user.ID,
		FirstName: user.FirstName,  // Add firstName
		LastName:  user.LastName,  
		Email:    user.Email,
		IsActive: user.IsActive,
	}

	// logger.Info("user registered", "email", req.Email)

	logger.Info("user registered", "email", req.Email)
      return &RegisterResponse{
	User:  userResponse,
	Message: "Registration successful",
	}, nil
}

// Login validates credentials and returns a signed JWT.
func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return nil, errors.New("invalid email or password")
	}

	if !password.Check(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.FirstName, user.LastName)
	if err != nil {
		return nil, err
	}

	
	// Convert to response DTO
	userResponse := response.UserResponse{
		ID:       user.ID, 
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:    user.Email,
		IsActive: user.IsActive,
	}

	logger.Info("user logged in", "email", req.Email)
	

	return &LoginResponse{
		User:  userResponse,
		Token: token,
	}, nil
}

// GetUserProfile returns user profile without sensitive data
func (s *AuthService) GetUserProfile(ctx context.Context, userID int64) (*response.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	userResponse := &response.UserResponse{
		ID:       user.ID,
		FirstName: user.FirstName,  // Add firstName
		LastName:  user.LastName, 
		Email:    user.Email,
		IsActive: user.IsActive,
	}

	return userResponse, nil
}

// Verify2FA checks OTP validity and activates user 2FA.
func (s *AuthService) Verify2FA(ctx context.Context, email, code string) error {
	valid, err := s.otpRepo.VerifyOTP(ctx, email, code, string(constants.Type2FA))
	if err != nil || !valid {
		return errors.New("invalid or expired code")
	}
	return nil
}

// RequestPasswordReset initiates password reset flow
func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) error {
	// Check if user exists (but don't reveal if they don't)
	user, _ := s.userRepo.FindByEmail(ctx, email)
	if user == nil {
		// Still return success to prevent email enumeration
		logger.Info("Password reset requested for non-existent email", "email", email)
		return nil
	}

	// Generate OTP code
	code := generateRandomCode(6)

	// Store OTP with password_reset type
	otp := &models.OTP{
		UserID: &user.ID,
		Email:  email,
		Code:   code,
		Type:   string(constants.TypePasswordReset),
	}

	if err := s.otpRepo.CreateOTP(ctx, otp); err != nil {
		return err
	}

	// TODO: Send email with the code
	logger.Info("Password reset code generated", "email", email, "code", code)

	return nil
}

// ResetPassword resets user password after OTP verification
func (s *AuthService) ResetPassword(ctx context.Context, email, code, newPassword string) error {
	// Verify the reset code
	valid, err := s.otpRepo.VerifyOTP(ctx, email, code, string(constants.TypePasswordReset))
	if err != nil || !valid {
		return errors.New("invalid or expired reset code")
	}

	// Find the user
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	// Hash new password
	hashedPassword, err := password.Hash(newPassword)
	if err != nil {
		return err
	}

	// Update user's password
	user.Password = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	logger.Info("Password reset successful", "email", email)
	return nil
}

// Send2FAOTP generates and sends a 2FA OTP code
func (s *AuthService) Send2FAOTP(ctx context.Context, email string) error {
	// Check if user exists
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	// Generate OTP code
	code := generateRandomCode(6)

	// Store OTP with 2fa type
	otp := &models.OTP{
		UserID: &user.ID,
		Email:  email,
		Code:   code,
		Type:   string(constants.Type2FA),
	}

	if err := s.otpRepo.CreateOTP(ctx, otp); err != nil {
		return err
	}

	// TODO: Send email with the code
	logger.Info("2FA code sent", "email", email, "code", code)

	return nil
}

// UpdateProfile updates user profile information
func (s *AuthService) UpdateProfile(ctx context.Context, userID int64, firstName, lastName, email string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	// If email is being changed, check it's not already taken
	if email != "" && email != user.Email {
		existingUser, _ := s.userRepo.FindByEmail(ctx, email)
		if existingUser != nil {
			return errors.New("email already exists")
		}
		user.Email = email
	}

	// Update other fields if provided
	if firstName != "" {
		user.FirstName = firstName
	}
	if lastName != "" {
		user.LastName = lastName
	}

	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	logger.Info("profile updated", "userID", userID)
	return nil
}

// Helper function to generate random code
func generateRandomCode(length int) string {
	const digits = "0123456789"
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = digits[time.Now().UnixNano()%int64(len(digits))]
	}
	return string(bytes)
}