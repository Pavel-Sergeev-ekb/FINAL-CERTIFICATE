package server

import (
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Pavel-Sergeev-ekb/FINAL-CERTIFICATE/internal/handlers"
)

type AppServer struct {
	Logger *log.Logger
	Server *http.Server
}

func NewServer(logger *log.Logger) *AppServer {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.BaseHandle)

	server := &http.Server{
		Addr:         ":7540",
		Handler:      mux,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &AppServer{
		Logger: logger,
		Server: server,
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
	_, portStr, _ := net.SplitHostPort(server.addr)
	port, _ := strconv.Atoi(portStr)
	return port
}
