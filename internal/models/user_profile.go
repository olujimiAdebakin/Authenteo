package models



// UserProfile stores non-critical, mutable user details, separating them
// from the core authentication fields in the User model.
type UserProfile struct {
	BaseModel

	// UserID is the foreign key linking the profile back to the core User model.
	UserID int64 `json:"user_id" db:"user_id"`

	// FirstName and LastName are the user's real names.
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`

	// DisplayName is a user-selected name, often used in public views.
	DisplayName string `json:"display_name" db:"display_name"`

	// AvatarURL stores the link to the user's profile picture.
	AvatarURL *string `json:"avatar_url" db:"avatar_url"`

	// Bio is a short description about the user.
	Bio *string `json:"bio" db:"bio"`
}

// UpdateProfileRequest defines the input structure for updating a user's profile.
type UpdateProfileRequest struct {
	FirstName   string `json:"first_name" validate:"omitempty"`
	LastName    string `json:"last_name" validate:"omitempty"`
	DisplayName string `json:"display_name" validate:"omitempty"`
	Bio         string `json:"bio" validate:"omitempty,max=500"`
	AvatarURL   *string `json:"avatar_url" validate:"omitempty,url"`
}