package server

import (
	"log"
	"net/http"
	"os"

	"go_final_project/pkg/api"
)

func Run() {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("web")))

	api.RegisterRoutes(mux)

	log.Printf("Server is swimming on port :%s", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server stopped with error: %v", err)
	}
}
