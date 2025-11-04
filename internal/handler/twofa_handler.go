package handler

import (
	"net/http"

	"authentio/internal/service"

	"github.com/gin-gonic/gin"
)

// TwoFAHandler handles 2FA-related HTTP requests
type TwoFAHandler struct {
	authService service.AuthService
}

// NewTwoFAHandler creates a new TwoFAHandler instance
func NewTwoFAHandler(authService service.AuthService) *TwoFAHandler {
	return &TwoFAHandler{
		authService: authService,
	}
}

// SendOTP initiates 2FA by sending an OTP
func (h *TwoFAHandler) SendOTP(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.Send2FAOTP(c.Request.Context(), req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

// VerifyOTP validates the provided OTP for 2FA
func (h *TwoFAHandler) VerifyOTP(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.Verify2FA(c.Request.Context(), req.Email, req.Code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
}