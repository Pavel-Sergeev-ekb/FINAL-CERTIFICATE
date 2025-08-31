package main

import (
	"log"
	"os"

	"github.com/Pavel-Sergeev-ekb/FINAL-CERTIFICATE/backend/internal/data"
	"github.com/Pavel-Sergeev-ekb/FINAL-CERTIFICATE/backend/internal/server"
)

func main() {

	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		log.Println("TODO_PASSWORD is not set, using default password")
		os.Setenv("TODO_PASSWORD", "12345") // Устанавливаем значение по умолчанию
	}
	log.Println("Current TODO_PASSWORD:", os.Getenv("TODO_PASSWORD"))

	database, err := data.Init(data.GetDBPath())
	if err != nil {
		log.Fatalf("database initialization error: %v", err)
	}
	defer database.Close()

	logger := log.New(
		os.Stdout,
		"[SERVER]",
		log.Ldate|log.Ltime|log.Lmicroseconds,
	)

	// Передаем db в сервер
	srv := server.NewServer(logger, database)

	addr := srv.Server.Addr
	logger.Printf("Server is starting on %s\n", addr)

	if err := srv.Server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
