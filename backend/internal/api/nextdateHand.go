package api

import (
	"log"
	"net/http"
	"time"

	"github.com/Pavel-Sergeev-ekb/FINAL-CERTIFICATE/backend/internal/services"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	dStartStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	if dStartStr == "" || repeatStr == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse("20060102", nowStr)
		if err != nil {
			http.Error(w, "Invalid date format for now", http.StatusBadRequest)
			return
		}
	}

	next, err := services.NextDate(now, dStartStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Successful calculation: nextDate=%s", next)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.Write([]byte(next))

	log.Printf("Response sent to client: %s", next)
}
