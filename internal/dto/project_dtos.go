package dto

type UpdateProjectRequest struct {
	NewTitle *string `json:"new_title" validate:"required"`
}

type CreateProjectRequest struct {
	Title string `json:"title" binding:"required" validate:"required"`
}
