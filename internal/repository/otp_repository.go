package repository

import (
	"context"
	"authentio/internal/models"
)

type OTPRepository interface {
	// CreateOTP creates a new OTP code
	CreateOTP(ctx context.Context, otp *models.OTP) error
	
	// VerifyOTP verifies an OTP code and marks it as used
	VerifyOTP(ctx context.Context, email, code, otpType string) (bool, error)
	
	// CleanupExpiredOTPs removes expired OTP codes
	CleanupExpiredOTPs(ctx context.Context) error
}