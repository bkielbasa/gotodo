package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/bkielbasa/gotodo/handlers"
	"github.com/bkielbasa/gotodo/repositories"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	db := getDB()
	repo := repositories.NewPostgres(db)
	r := mux.NewRouter()
	todoHandler := handlers.ToDo{Repo: repo}
	projectHandler := handlers.Project{Repo: repo}
	r.HandleFunc("/projects", projectHandler.List)
	r.HandleFunc("/project/create", projectHandler.Create)
	r.HandleFunc("/project/{id:[0-9a-z\\-]+}/archive", projectHandler.Archive)
	r.HandleFunc("/projects", projectHandler.List)
	r.HandleFunc("/todos", todoHandler.List)
	r.HandleFunc("/todo/create", todoHandler.Create)
	r.HandleFunc("/todo/{id:[0-9a-z\\-]+}", todoHandler.Get)
	r.HandleFunc("/todo/{id:[0-9a-z\\-]+}/done", todoHandler.MarkAsDone)
	r.HandleFunc("/todo/{id:[0-9a-z\\-]+}/undone", todoHandler.MarkAsUndone)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getDB() *sql.DB {
	url := os.Getenv("POSTGRES_URL")
	db, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}

	return db
}
