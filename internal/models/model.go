package models


import "time"

type BaseModel struct {
	ID        int64     `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
	ExpiredAt *time.Time `json:"expired_at,omitempty" db:"expired_at"`
}