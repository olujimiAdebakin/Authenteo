package handler



import (
	"regexp"

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
		re := regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[\W_]).{8,}$`)
		return re.MatchString(fl.Field().String())
	})

}

	func formatValidationError(err error) map[string]string {
	errs := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		errs[e.Field()] = e.Tag()
	}
	return errs
}

