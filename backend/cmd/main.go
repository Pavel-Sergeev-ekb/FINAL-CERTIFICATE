package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/Pavel-Sergeev-ekb/FINAL-CERTIFICATE/backend/internal/database"
	"github.com/Pavel-Sergeev-ekb/FINAL-CERTIFICATE/backend/internal/server"
)

func main() {

	db, err := sql.Open("sqlite3", database.GetDBPath())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("failed to connect db:", err)
	}
	if err := database.InitDB(db); err != nil {
		log.Fatal(err)
	}

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
