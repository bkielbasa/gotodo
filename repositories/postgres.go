package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/bkielbasa/gotodo/httpmodels"
)

type Postgres struct {
	db *sql.DB
}

type PostgresRepository interface {
	GetToDo(id string) (string, bool, error)
	AddToDo(name string) (string, error)
	MarkToDoAsDone(id string) error
	MarkToDoAsUndone(id string) error
	ListToDos() (httpmodels.ListToDoResponse, error)
	ListProjects() (httpmodels.ListProjectsResponse, error)
	CreateProject(name string) (string, error)
	ArchiveProject(id string) error
}

func NewPostgres(db *sql.DB) Postgres {
	return Postgres{db}
}

func (p Postgres) AddToDo(name string) (string, error) {
	q := fmt.Sprintf(`INSERT INTO todos ("name") VALUES ('%s') returning id`, name)
	resp, err := p.db.Query(q)

	if err != nil {
		return "", err
	}

	resp.Next()

	var id string
	err = resp.Scan(&id)
	if err != nil {
		return "", err
	}

	return id, err
}

func (p Postgres) GetToDo(id string) (string, bool, error) {
	q := fmt.Sprintf(`SELECT "name", done FROM todos where id = '%s'`, id)
	row := p.db.QueryRow(q)

	var name string
	var done bool
	err := row.Scan(&name, &done)

	return name, done, err
}

func (p Postgres) MarkToDoAsDone(id string) error {
	q := fmt.Sprintf(`UPDATE todos SET done = true where id = '%s'`, id)
	res, err := p.db.Exec(q, id)

	if err != nil {
		return err
	}

	if c, _ := res.RowsAffected(); c != 1 {
		return errors.New("the to do does not exist")
	}

	return err
}

func (p Postgres) MarkToDoAsUndone(id string) error {
	q := fmt.Sprintf(`UPDATE todos SET done = false where id = '%s'`, id)
	res, err := p.db.Exec(q, id)

	if err != nil {
		return err
	}

	if c, _ := res.RowsAffected(); c != 1 {
		return errors.New("the to do does not exist")
	}

	return err
}

func (p Postgres) ListToDos() (httpmodels.ListToDoResponse, error) {
	q := `SELECT * FROM todos`
	res, err := p.db.Query(q)

	if err != nil {
		return httpmodels.ListToDoResponse{}, err
	}

	resp := httpmodels.ListToDoResponse{}

	for res.Next() {
		var todo httpmodels.GetToDoResponse
		err = res.Scan(&todo.ID, &todo.Name, &todo.Done)
		if err != nil {
			return httpmodels.ListToDoResponse{}, err
		}

		resp.ToDos = append(resp.ToDos, todo)
	}

	return resp, err
}

func (p Postgres) CreateProject(name string) (string, error) {
	q := `INSERT INTO project ("name") VALUES ($1) returning id`
	resp, err := p.db.Query(q, name)

	if err != nil {
		return "", err
	}

	resp.Next()

	var id string
	err = resp.Scan(&id)
	if err != nil {
		return "", err
	}

	return id, err
}

func (p Postgres) ListProjects() (httpmodels.ListProjectsResponse, error) {
	q := `SELECT * FROM project`
	res, err := p.db.Query(q)

	if err != nil {
		return httpmodels.ListProjectsResponse{}, err
	}

	resp := httpmodels.ListProjectsResponse{}

	for res.Next() {
		var project httpmodels.GetProjectResponse
		err = res.Scan(&project.ID, &project.Name, &project.Archived)
		if err != nil {
			return httpmodels.ListProjectsResponse{}, err
		}

		resp.Projects = append(resp.Projects, project)
	}

	return resp, err
}

func (p Postgres) ArchiveProject(id string) error {
	q := fmt.Sprintf(`UPDATE project SET archived = true where id = '%s'`, id)
	res, err := p.db.Exec(q)

	if err != nil {
		return err
	}

	if c, _ := res.RowsAffected(); c != 1 {
		return errors.New("the project does not exist")
	}

	return err
}
