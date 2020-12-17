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

type ToDo struct {
	Repo repositories.PostgresRepository
}

func (t ToDo) Create(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("cannot read the body: %s", err)
		http.Error(w, "cannot read the body", http.StatusBadRequest)
		return
	}

	req := httpmodels.CreateToDoRequest{}
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

	if req.ProjectID == "" {
		log.Printf("the ProjectID is required")
		http.Error(w, "the ProjectID is required", http.StatusBadRequest)
		return
	}

	id, err := t.Repo.AddToDo(req.ProjectID, req.Name)
	if err != nil {
		log.Printf("internal server error: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := httpmodels.CreateToDoResponse{ID: id}
	b, _ = json.Marshal(resp)

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

func (t ToDo) List(w http.ResponseWriter, r *http.Request) {
	resp, err := t.Repo.ListToDos()
	if err != nil {
		log.Printf("internal server error: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(resp)

	w.Header().Add("content-type", "application/json")
	_, _ = w.Write(b)
}

func (t ToDo) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	name, done, err := t.Repo.GetToDo(id)
	if err != nil {
		log.Printf("internal server error: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := httpmodels.GetToDoResponse{
		ID:   id,
		Name: name,
		Done: done,
	}

	b, _ := json.Marshal(resp)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

func (t ToDo) MarkAsDone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := t.Repo.MarkToDoAsDone(id)
	if err != nil {
		if err.Error() == "the to do does not exist" {
			http.NotFound(w, r)
			return
		}

		log.Printf("internal server error: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (t ToDo) MarkAsUndone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := t.Repo.MarkToDoAsUndone(id)
	if err != nil {
		if err.Error() == "the to do does not exist" {
			http.NotFound(w, r)
			return
		}

		log.Printf("internal server error: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
