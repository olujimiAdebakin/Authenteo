package models

type User struct {
	BaseModel
	FirstName string `json:"first_name" db:"first_name"`
      LastName  string `json:"last_name" db:"last_name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"-" db:"password"`
	IsActive bool   `json:"is_active" db:"is_active"`
}