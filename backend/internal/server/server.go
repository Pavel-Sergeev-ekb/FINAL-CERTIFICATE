package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Pavel-Sergeev-ekb/FINAL-CERTIFICATE/backend/internal/api"
	"github.com/Pavel-Sergeev-ekb/FINAL-CERTIFICATE/backend/internal/services"
)

type AppServer struct {
	Logger *log.Logger
	Server *http.Server
	DB     *sql.DB
}

const BasePort = 7540

func NewServer(logger *log.Logger, db *sql.DB) *AppServer {

	mux := http.NewServeMux()

	mux.HandleFunc("/", BaseHandle)

	mux.HandleFunc("/api/nextdate", api.NextDateHandler)

	mux.HandleFunc("/api/signin", api.SignInHandler)

	mux.HandleFunc("/api/task", services.AuthMiddleware(api.TaskHandler(db)))

	mux.HandleFunc("/api/tasks", services.AuthMiddleware(api.GetTasksHandler(db)))

	mux.HandleFunc("/api/task/done", services.AuthMiddleware(api.DoneHandler(db)))

	// Fetch tasks (protected)
	port := GetPort()
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &AppServer{
		Logger: logger,
		Server: server,
		DB:     db,
	}

}

func GetPort() int {

	envPort := os.Getenv("TODO_PORT")
	if envPort != "" {
		port, err := strconv.Atoi(envPort)
		if err == nil {
			return port
		}
	}

	return BasePort
}
