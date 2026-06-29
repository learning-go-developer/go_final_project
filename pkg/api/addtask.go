package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go_final_project/pkg/db"
)

func sendJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func normalizeTaskDate(t *db.Task) error {
	now := time.Now()
	todayStr := now.Format(dateLayout)

	if strings.TrimSpace(t.Date) == "" {
		t.Date = todayStr
	}

	taskTime, err := time.Parse(dateLayout, t.Date)
	if err != nil {
		return errors.New("некорректный формат даты, ожидается YYYYMMDD")
	}

	taskDay := time.Date(taskTime.Year(), taskTime.Month(), taskTime.Day(), 0, 0, 0, 0, time.UTC)
	todayDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	if taskDay.Before(todayDay) {
		if strings.TrimSpace(t.Repeat) == "" {
			t.Date = todayStr
		} else {
			next, err := NextDate(now, t.Date, t.Repeat)
			if err != nil {
				return err
			}
			t.Date = next
		}
	} else if strings.TrimSpace(t.Repeat) != "" {
		_, err := NextDate(now, t.Date, t.Repeat)
		if err != nil {
			return err
		}
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "ошибка десериализации JSON: " + err.Error()})
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

	id, err := db.AddTask(&task)
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, map[string]string{"error": "ошибка сохранения в БД: " + err.Error()})
		return
	}

	sendJSON(w, http.StatusOK, map[string]string{"id": strconv.FormatInt(id, 10)})
}
