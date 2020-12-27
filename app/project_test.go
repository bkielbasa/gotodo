package app

import (
	"context"
	"errors"
	"github.com/bkielbasa/gotodo/domain"
	"github.com/google/uuid"
	"testing"
)

func TestAddNewProject(t *testing.T) {
	name := "my name:" + uuid.New().String()
	ctx := context.Background()

	projectServ := NewProjectService(newStoreMock())
	p, err := projectServ.Add(ctx, name)
	if err != nil {
		t.Errorf("expected no error but got: %s", err)
	}

	checkProjectName(t, p, name)

	p2, err := projectServ.Get(ctx, p.ID())
	if err != nil {
		t.Errorf("expected no error but got: %s", err)
	}

	checkProjectName(t, p2, p.Name())
	checkProjectID(t, p2, p.ID())
}

func checkProjectID(t *testing.T, p domain.Project, expectedID string) {
	if p.ID() !=  expectedID {
		t.Errorf("expected ID %s but %s given", expectedID,  p.ID())
	}
}

func checkProjectName(t *testing.T, p domain.Project, expectedName string) {
	if p.Name() !=  expectedName {
		t.Errorf("expected name %s but %s given", expectedName,  p.Name())
	}
}

func TestEveryProjectShouldHaveUniqueID(t *testing.T) {
	name := "a name"

	projectServ := NewProjectService(newStoreMock())
	p1, err := projectServ.Add(context.Background(), name)
	if err != nil {
		t.Errorf("expected no error but got: %s", err)
	}

	p2, err := projectServ.Add(context.Background(), name)
	if err != nil {
		t.Errorf("expected no error but got: %s", err)
	}

	if p1.ID() == p2.ID() {
		t.Error("every project should have a unique ID")
	}
}

func TestAddNewProjectWithEmptyName(t *testing.T) {
	name := ""

	projectServ := NewProjectService(newStoreMock())
	_, err := projectServ.Add(context.Background(), name)
	if err == nil {
		t.Errorf("expected error but got nil")
	}
}

func TestAGetNotExistingProject(t *testing.T) {
	id := "not exists"
	ctx := context.Background()
	storage := newStoreMock().withError(ErrProjectNotFound)

	projectServ := NewProjectService(storage)

	_, err := projectServ.Get(ctx, id)
	if !errors.Is(err, ErrProjectNotFound) {
		t.Errorf("expected error ErrProjectNotFound but got %v", err)
	}
}

type storeMock struct {
	data map[string]domain.Project
	err error
}

func newStoreMock() *storeMock {
	return &storeMock{
		data: make(map[string]domain.Project),
	}
}

func (s *storeMock) Store(ctx context.Context, p domain.Project) error {
	s.data[p.ID()] = p
	return nil
}

func (s *storeMock) Get(ctx context.Context, id string) (domain.Project, error) {
	return s.data[id], s.err
}

func (s *storeMock) withError(err error) *storeMock {
	s.err = err
	return s
}

