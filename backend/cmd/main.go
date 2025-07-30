package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Pavel-Sergeev-ekb/FINAL-CERTIFICATE/internal/server"
)

func main() {
	port := server.GetPort
	logger := log.New(
		os.Stdout,
		"[%s]",
		log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC,
	)
	srv := server.NewServer(logger)

	addr := srv.Server.Addr

	logger.Printf("Server is starting on %s\n", &addr)

	fmt.Println("Get port %s", port)

	if err := srv.Server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
