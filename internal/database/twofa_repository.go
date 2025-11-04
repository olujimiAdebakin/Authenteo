package database


import (
	"context"
	"database/sql"
	"errors"
	_"authentio/internal/models"
	"authentio/internal/repository"
)

type twoFARepository struct {
	db *sql.DB
}

func NewTwoFARepository(db *sql.DB) repository.TwoFARepository {
	return &twoFARepository{db: db}
}

func (r *twoFARepository) Enable2FA(ctx context.Context, userID int64, secret string) error {
	// For email OTP, secret is not used
	query := `
		INSERT INTO two_fa_configs (user_id, method, secret, enabled) 
		VALUES ($1, 'email', '', TRUE)
		ON CONFLICT (user_id) 
		DO UPDATE SET method = 'email', enabled = TRUE, updated_at = CURRENT_TIMESTAMP`
	
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *twoFARepository) Disable2FA(ctx context.Context, userID int64) error {
	query := `UPDATE two_fa_configs SET enabled = FALSE WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *twoFARepository) Is2FAEnabled(ctx context.Context, userID int64) (bool, error) {
	query := `SELECT enabled FROM two_fa_configs WHERE user_id = $1`
	
	var enabled bool
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&enabled)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	
	return enabled, nil
}

func (r *twoFARepository) Get2FASecret(ctx context.Context, userID int64) (string, error) {
	// Not used for email OTP
	return "", errors.New("not applicable for email OTP")
}

func (r *twoFARepository) VerifyCode(ctx context.Context, userID int64, code string) (bool, error) {
	// This method is not used for email OTP
	// Email OTP verification is handled by OTPRepository
	return false, errors.New("use OTPRepository for email OTP verification")
}