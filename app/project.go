package app

import (
	"context"
	"errors"
	"github.com/bkielbasa/gotodo/domain"
	"github.com/google/uuid"
)

var ErrProjectNotFound = errors.New("the project is not found")

type ProjectService struct {
	repo Repository
}

func NewProjectService(storage Repository) ProjectService {
	return ProjectService{repo: storage}
}

type Repository interface {
	Store(ctx context.Context, p domain.Project) error
	Get(ctx context.Context, id string) (domain.Project, error)
}

func (serv ProjectService) Add(ctx context.Context, name string) (domain.Project, error) {
	id := uuid.New().String()
	p, err := domain.NewProject(id, name)
	if err != nil {
		return domain.Project{}, err
	}

	err = serv.repo.Store(ctx, p)
	if err != nil {
		return domain.Project{}, err
	}

	return p, err
}

func (serv ProjectService) Get(ctx context.Context, id string) (domain.Project, error) {
	return serv.repo.Get(ctx, id)
}
