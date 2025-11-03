package response

import (
	
	"time"
)





type UserResponse struct {
     ID        string    `json:"id"`
    Email     string    `json:"email"`
    IsActive  bool      `json:"is_active"`
    CreatedAt time.Time `json:"created_at,omitempty"`
}