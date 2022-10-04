package dto

type CreateUserDto struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"username" binding:"required" validate:"required, email"`
	Password string `json:"password" binding:"required" validate:"required, min=8"`
}
