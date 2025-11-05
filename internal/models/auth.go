package models

type RegisterRequest struct {
	FirstName string `json:"first_name" db:"first_name" validate:"required,alphaSpace,min=2,max=50"`
	LastName  string `json:"last_name" db:"last_name" validate:"required,alphaSpace,min=2,max=50"`
	Email     string `json:"email" db:"email" validate:"required,email,max=50"`
	Password  string `json:"password" db:"password" validate:"required,password"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required"`
}


