package models

// TwoFA represents a user's persistent Two-Factor Authentication configuration.
// It stores the long-term secret and the enabled status.
// This model is kept separate from OTP (transient codes) for architectural clarity.
type TwoFA struct {
	BaseModel
	
	// UserID links this 2FA configuration to the user.
	UserID int64 `db:"user_id" json:"user_id"`
	
	// Method is the type of 2FA (e.g., "totp", "email", "sms").
	Method string `db:"method" json:"method"` 
	
	// Secret is the cryptographic key used to generate OTPs. 
	// The json:"-" tag ensures this highly sensitive data is never sent in API responses.
	Secret string `db:"secret" json:"-"`
	
	// Enabled indicates if the user has completed 2FA setup and it is active.
	Enabled bool `db:"enabled" json:"enabled"`
}
