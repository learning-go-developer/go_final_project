package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", nextDateHandler)

	mux.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			addTaskHandler(w, r)
		case http.MethodGet:
			getTaskHandler(w, r)
		case http.MethodPut:
			updateTaskHandler(w, r)
		case http.MethodDelete:
			deleteTaskHandler(w, r)
		default:
			sendJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не поддерживается"})
		}
	})

	mux.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			sendJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не поддерживается"})
			return
		}
		getTasksHandler(w, r)
	})

	mux.HandleFunc("/api/task/done", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			sendJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не поддерживается"})
			return
		}
		finishTaskHandler(w, r)
	})
}
