package service

import (
	"context"

	"gorest-api/internal/dto"
	"gorest-api/internal/model"
	"gorest-api/internal/repository"
)

type ProjectsService struct {
	repo repository.Projects
}

func NewProjectsService(repo repository.Projects) *ProjectsService {
	return &ProjectsService{repo: repo}
}

func (s *ProjectsService) Create(ctx context.Context, userId string, dto dto.CreateProjectDto) (string, error) {
	project := model.Project{
		Title:    dto.Title,
		UserID:   userId,
		Archived: false,
	}

	return s.repo.Create(ctx, project)
}

func (s *ProjectsService) Update(ctx context.Context, userId string, projectId string, newTitle string) error {
	return s.repo.Update(ctx, userId, projectId, newTitle)
}

func (s *ProjectsService) GetAll(ctx context.Context) ([]model.Project, error) {
	return s.repo.GetAll(ctx)
}

func (s *ProjectsService) GetByTitle(ctx context.Context, userId string, title string) (model.Project, error) {
	return s.repo.GetByTitle(ctx, userId, title)
}

func (s *ProjectsService) Delete(ctx context.Context, userId string, projectId string) error {
	return s.repo.Delete(ctx, userId, projectId)
}
