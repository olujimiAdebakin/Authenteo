package handler

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func InitValidator() {
	Validate = validator.New()

	// Allow only alphabets and spaces for names
	Validate.RegisterValidation("alphaSpace", func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile(`^[a-zA-Z\s]+$`)
		return re.MatchString(fl.Field().String())
	})

	// Enforce strong password policy:
	// - at least one lowercase
	// - at least one uppercase
	// - at least one number
	// - at least one special character
	// - at least 8 chars total
	Validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		
		// Check minimum length
		if len(password) < 8 {
			return false
		}
		
		// Check for at least one lowercase letter
		hasLower := false
		for _, r := range password {
			if r >= 'a' && r <= 'z' {
				hasLower = true
				break
			}
		}
		if !hasLower {
			return false
		}
		
		// Check for at least one uppercase letter
		hasUpper := false
		for _, r := range password {
			if r >= 'A' && r <= 'Z' {
				hasUpper = true
				break
			}
		}
		if !hasUpper {
			return false
		}
		
		// Check for at least one digit
		hasDigit := false
		for _, r := range password {
			if r >= '0' && r <= '9' {
				hasDigit = true
				break
			}
		}
		if !hasDigit {
			return false
		}
		
		// Check for at least one special character
		hasSpecial := false
		specialChars := "!@#$%^&*()-_=+[]{}|;:',.<>?/`~"
		for _, r := range password {
			for _, s := range specialChars {
				if r == s {
					hasSpecial = true
					break
				}
			}
			if hasSpecial {
				break
			}
		}
		
		return hasSpecial
	})

}

func FormatValidationError(err error) map[string]string {
	errs := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		// Here is a more user-friendly error messages
		switch e.Tag() {
		case "required":
			errs[strings.ToLower(e.Field())] = "This field is required"
		case "email":
			errs[strings.ToLower(e.Field())] = "Invalid email format"
		case "min":
			errs[strings.ToLower(e.Field())] = "Value is too short"
		case "password":
			errs[strings.ToLower(e.Field())] = "Password must contain uppercase, lowercase, number, and special character"
		case "alphaSpace":
			errs[strings.ToLower(e.Field())] = "Only letters and spaces are allowed"
		default:
			errs[strings.ToLower(e.Field())] = "Invalid value"
		}
	}
	return errs
}

