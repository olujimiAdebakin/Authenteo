package models

type OTP struct {
	BaseModel
	UserID *int64 `db:"user_id" json:"-"`
	Email     string    `json:"email"`
	Code      string    `json:"code"`
	Type string `db:"type" json:"type"`
}