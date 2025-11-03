package models

type User struct {
	BaseModel
	Email    string `json:"email" db:"email"`
	Password string `json:"-" db:"password"`
	IsActive bool   `json:"is_active" db:"is_active"`
}