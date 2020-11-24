package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bkielbasa/gotodo/handlers"
	"github.com/bkielbasa/gotodo/repositories"
	_ "github.com/lib/pq"
)

func main() {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	ctx := context.Background()
	run, shutdown := App(ctx, 8090)

	go func() {
		_ = <-gracefulStop
		fmt.Println("shutting down...")
		err := shutdown()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	err := run()
	if !errors.Is(err, http.ErrServerClosed) {
		fmt.Println(err)
		os.Exit(1)
	}
}

func App(ctx context.Context, port int) (func() error, func() error) {
	m := http.NewServeMux()
	s := http.Server{Addr: fmt.Sprintf(":%d", port), Handler: m}
	log.Printf("starting on port %d", port)

	db := getDB()
	repo := repositories.NewPostgres(db)
	todoHandler := handlers.ToDo{Repo: repo}
	projectHandler := handlers.Project{Repo: repo}
	m.HandleFunc("/projects", projectHandler.List)
	m.HandleFunc("/project/create", projectHandler.Create)
	m.HandleFunc("/project/{id:[0-9a-z\\-]+}/archive", projectHandler.Archive)
	m.HandleFunc("/todos", todoHandler.List)
	m.HandleFunc("/todo/create", todoHandler.Create)
	m.HandleFunc("/todo/{id:[0-9a-z\\-]+}", todoHandler.Get)
	m.HandleFunc("/todo/{id:[0-9a-z\\-]+}/done", todoHandler.MarkAsDone)
	m.HandleFunc("/todo/{id:[0-9a-z\\-]+}/undone", todoHandler.MarkAsUndone)
	return s.ListenAndServe, func() error {
		return s.Shutdown(ctx)
	}
}

func getDB() *sql.DB {
	url := os.Getenv("POSTGRES_URL")
	db, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}

	return db
}
