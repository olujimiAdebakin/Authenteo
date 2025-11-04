package database

import (
	"context"
	"database/sql"
	"time"
	"authentio/internal/models"
	"authentio/internal/repository"
)

type otpRepository struct {
	db *sql.DB
}

func NewOTPRepository(db *sql.DB) repository.OTPRepository {
	return &otpRepository{db: db}
}

func (r *otpRepository) CreateOTP(ctx context.Context, otp *models.OTP) error {
	// Set expiration to 10 minutes
	expiredAt := time.Now().Add(10 * time.Minute)
	otp.ExpiredAt = &expiredAt

	query := `
		INSERT INTO otps (user_id, email, code, type, expires_at) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`
	
	err := r.db.QueryRowContext(ctx, query,
		otp.UserID,
		otp.Email,
		otp.Code,
		otp.Type,
		otp.ExpiredAt,
	).Scan(&otp.ID, &otp.CreatedAt)
	
	return err
}

func (r *otpRepository) VerifyOTP(ctx context.Context, email, code, otpType string) (bool, error) {
	query := `
		UPDATE otps 
		SET used = TRUE 
		WHERE email = $1 AND code = $2 AND type = $3 
		AND used = FALSE AND expires_at > $4
		RETURNING id`
	
	var id int64
	err := r.db.QueryRowContext(ctx, query, email, code, otpType, time.Now()).Scan(&id)
	
	if err == sql.ErrNoRows {
		return false, nil // Code not found or expired
	}
	if err != nil {
		return false, err
	}
	
	return true, nil
}

func (r *otpRepository) CleanupExpiredOTPs(ctx context.Context) error {
	query := `DELETE FROM otps WHERE expires_at < $1`
	_, err := r.db.ExecContext(ctx, query, time.Now())
	return err
}