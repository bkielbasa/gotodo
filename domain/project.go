package domain

import "errors"

type Project struct {
	id string
	name string
}

func (p Project) Name() string {
	return p.name
}

func (p Project) ID() string {
	return p.id
}

func NewProject(id, name string) (Project, error)  {
	if id == "" {
		return Project{}, errors.New("the ID cannot be empty")
	}

	if name == "" {
		return Project{}, errors.New("the name cannot be empty")
	}

	return Project{id: id, name: name}, nil
}
