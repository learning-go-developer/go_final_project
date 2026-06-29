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

	api.Init()

	http.Handle("/", http.FileServer(http.Dir("web")))

	log.Printf("server started on :%s", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
