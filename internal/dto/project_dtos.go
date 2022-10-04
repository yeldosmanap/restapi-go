package dto

type UpdateProjectDto struct {
	NewTitle *string `json:"new_title" validate:"required"`
}

type CreateProjectDto struct {
	Title  string `json:"title" binding:"required" validate:"required"`
}
