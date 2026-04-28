package api

import (
	"net/http"
	"os"
	"strings"

	"scheduler/pkg/db"
)

type App struct {
	store    *db.Store
	password string
}

func New(store *db.Store) *App {
	return &App{
		store:    store,
		password: strings.TrimSpace(os.Getenv("TODO_PASSWORD")),
	}
}

func (a *App) Init() {
	http.HandleFunc("/api/signin", a.signinHandler)
	http.HandleFunc("/api/nextdate", a.nextDateHandler)
	http.HandleFunc("/api/task", a.auth(a.taskHandler))
	http.HandleFunc("/api/tasks", a.auth(a.tasksHandler))
	http.HandleFunc("/api/task/done", a.auth(a.doneTaskHandler))
}
