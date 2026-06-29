package main

import (
	"log"
	"os"

	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
)

func getEnvOr(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func main() {
	databasePath := getEnvOr("TODO_DBFILE", "scheduler.db")

	if err := db.Init(databasePath); err != nil {
		log.Fatalf("database initialization failed: %v", err)
	}

	server.Run()
}
