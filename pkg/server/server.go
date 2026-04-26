package server

import (
	"net/http"
	"os"
	"strings"
)

const (
	defaultPort = "7540"
	webDir      = "web"
)

func getPort() string {
	port := strings.TrimSpace(os.Getenv("TODO_PORT"))
	if port == "" {
		return defaultPort
	}
	return port
}

func Run() error {
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	return http.ListenAndServe(":"+getPort(), nil)
}
