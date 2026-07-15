package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"go_final_project/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан id"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sendJSON(w, http.StatusNotFound, map[string]string{
				"error": "task not found",
			})
			return
		}

		log.Printf("failed to get task %q: %v", id, err)

		sendJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
		return
	}

	sendJSON(w, http.StatusOK, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "ошибка десериализации JSON"})
		return
	}

	if task.ID == "" {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан id задачи"})
		return
	}

	if strings.TrimSpace(task.Title) == "" {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	if err := normalizeTaskDate(&task); err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		sendJSON(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена или не обновлена"})
		return
	}

	sendJSON(w, http.StatusOK, map[string]string{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан id"})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		sendJSON(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена или не может быть удалена"})
		return
	}

	sendJSON(w, http.StatusOK, map[string]string{})
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")

	tasks, err := db.GetTasks(50, search)
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, map[string]string{"error": "ошибка получения задач из БД: " + err.Error()})
		return
	}

	sendJSON(w, http.StatusOK, map[string][]db.Task{"tasks": tasks})
}

func finishTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан id"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		sendJSON(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
		return
	}

	if strings.TrimSpace(task.Repeat) == "" {
		if err := db.DeleteTask(id); err != nil {
			sendJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	} else {
		nextDate, err := NextDate(time.Now().UTC(), task.Date, task.Repeat)
		if err != nil {
			sendJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		task.Date = nextDate
		if err := db.UpdateTask(task); err != nil {
			sendJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}

	sendJSON(w, http.StatusOK, map[string]string{})
}
