package service

import (
	"context"

	_ "github.com/golang/mock/gomock"

	"gorest-api/internal/dto"
	"gorest-api/internal/model"
	"gorest-api/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(ctx context.Context, user model.User) (string, error)
	GenerateToken(ctx context.Context, email, password string) (string, error)
	ParseToken(token string) (string, error)
}

type Projects interface {
	Create(ctx context.Context, userId string, dto dto.CreateProjectDto) (string, error)
	Update(ctx context.Context, userId string, projectId string, newTitle string) error
	Delete(ctx context.Context, userId string, projectId string) error
	GetByTitle(ctx context.Context, userId string, title string) (model.Project, error)
	GetAll(ctx context.Context) ([]model.Project, error)
}

type Service struct {
	Authorization
	Projects
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		Projects:      NewProjectsService(repo.Projects),
	}
}
