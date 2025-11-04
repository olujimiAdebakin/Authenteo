package response

import (
	
	"time"
)





type UserResponse struct {
    ID        int64    `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Email     string    `json:"email"`
    IsActive  bool      `json:"is_active"`
    CreatedAt time.Time `json:"created_at,omitempty"`
}

// I Added a helper method to get full name
func (u *UserResponse) GetFullName() string {
    return u.FirstName + " " + u.LastName
}