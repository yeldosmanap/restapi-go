package dto

type UpdateProject struct {
	NewTitle *string `json:"new_title" validate:"required"`
}

type CreateProject struct {
	Title string `json:"title" binding:"required" validate:"required"`
}
