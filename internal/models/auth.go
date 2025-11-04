package models



type RegisterRequest struct {
	FirstName string `json:"first_name" db:"first_name"`
      LastName  string `json:"last_name" db:"last_name"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}


