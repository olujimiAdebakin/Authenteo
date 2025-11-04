package handler

import "authentio/internal/service"

// Handler collects all HTTP handlers for the application
type Handler struct {
	Auth  *AuthHandler
	TwoFA *TwoFAHandler  
	User  *UserHandler
}

// NewHandler creates a new Handler instance with all dependencies
func NewHandler(authService service.AuthService) *Handler {
	return &Handler{
		Auth:  NewAuthHandler(authService),
		TwoFA: NewTwoFAHandler(authService), // Pass authService since it has 2FA methods
		User:  NewUserHandler(authService),  // Pass authService since it has user methods
	}
}