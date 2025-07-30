package main

import (
	"log"
	"os"

	"github.com/Pavel-Sergeev-ekb/FINAL-CERTIFICATE/backend/internal/server"
)

func main() {

	logger := log.New(
		os.Stdout,
		"",
		log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC,
	)
	srv := server.NewServer(logger)

	addr := srv.Server.Addr

	logger.Printf("Server is starting on %s\n", addr)

	if err := srv.Server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
