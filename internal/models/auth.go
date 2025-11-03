package models



type RegisterRequest struct {

	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}


// in internal/models/user.go or internal/models/auth_helpers.go

// func NewUserFromRegister(req RegisterRequest, hashedPassword string) *User {
//     return &User{
//         BaseModel: BaseModel{}, // fill ID/timestamps as needed from repo/ORM
//         Email:     req.Email,
//         Password:  hashedPassword,
//         IsActive:  true,
//     }
// }