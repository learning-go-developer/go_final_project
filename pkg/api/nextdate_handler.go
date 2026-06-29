package api

import (
	"net/http"
	"time"
)

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	var (
		now time.Time
		err error
	)

	if nowParam := r.FormValue("now"); nowParam != "" {
		if now, err = time.Parse(dateLayout, nowParam); err != nil {
			http.Error(w, "bad now parameter", http.StatusBadRequest)
			return
		}
	} else {
		now = time.Now().UTC()
	}

	res, err := NextDate(now, r.FormValue("date"), r.FormValue("repeat"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte(res))
}
