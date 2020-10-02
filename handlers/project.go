package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bkielbasa/gotodo/httpmodels"
	"github.com/bkielbasa/gotodo/repositories"
	"github.com/gorilla/mux"
)

type Project struct {
	Repo repositories.PostgresRepository
}

func (p Project) Create(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("cannot read the body: %s", err)
		http.Error(w, "cannot read the body", http.StatusBadRequest)
		return
	}

	req := httpmodels.CreateProjectRequest{}
	err = json.Unmarshal(b, &req)
	if err != nil {
		log.Printf("cannot read the body: %s", err)
		http.Error(w, "invalid JSON provided", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		log.Printf("the name is required")
		http.Error(w, "the name is required", http.StatusBadRequest)
		return
	}

	id, err := p.Repo.CreateProject(req.Name)
	if err != nil {
		log.Printf("internal server error: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := httpmodels.CreateProjectResponse{id}
	b, _ = json.Marshal(resp)

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_,_ = w.Write(b)
}

func (p Project) List(w http.ResponseWriter, r *http.Request) {
	resp, err := p.Repo.ListProjects()
	if err != nil {
		log.Printf("internal server error: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(resp)

	w.Header().Add("content-type", "application/json")
	_,_ = w.Write(b)
}

func (p Project) Archive(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := p.Repo.ArchiveProject(id)
	if err != nil {
		if err.Error() == "the project does not exist" {
			http.NotFound(w, r)
			return
		}

		log.Printf("internal server error: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}