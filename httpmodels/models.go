package httpmodels

type CreateToDoRequest struct {
	Name string `json:"name"`
	ProjectID string `json:"project_id"`
}

type CreateProjectRequest struct {
	Name string `json:"name"`
}

type CreateProjectResponse struct {
	ID string `json:"id"`
}

type CreateToDoResponse struct {
	ID string `json:"id"`
}

type GetToDoResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Done bool `json:"done"`
}

type ListToDoResponse struct {
	ToDos []GetToDoResponse `json:"todos"`
}

type GetProjectResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Archived bool `json:"archived"`
}

type ListProjectsResponse struct {
	Projects []GetProjectResponse `json:"projects"`
}
