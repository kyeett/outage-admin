package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"

	"github.com/go-chi/chi"
)

type orderFailure struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Message   string    `json:"message"`
	Handled   bool      `json:"handled"`
}

var orderFailures = []orderFailure{
	orderFailure{1, time.Now().Add(-115 * time.Minute), "Major outtage, lost 100 SEK", false},
	orderFailure{2, time.Now().Add(-75 * time.Minute), "Minor outtage, lost 10 SEK", true},
	orderFailure{3, time.Now().Add(-60 * time.Minute), "Warning outtage, lost 100 SEK", false},
	orderFailure{4, time.Now().Add(-30 * time.Minute), "Major outtage, lost 100000 SEK", false},
	orderFailure{5, time.Now().Add(-10 * time.Minute), "Warning outtage, lost 100 SEK", false},
	orderFailure{6, time.Now().Add(-10 * time.Minute), "Warning outtage, lost 100 SEK", false},
	orderFailure{7, time.Now().Add(-1000 * time.Minute), "Warning outtage, lost 1 SEK", true},
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	r.Mount("/static", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.Route("/order_failures", func(r chi.Router) {
		r.Get("/", listOrderFailures)
		r.Post("/mark_as_handled", markAsHandled)
		r.Post("/mark_as_unhandled", markAsUnhandled)
	})

	http.ListenAndServe(":"+port, r)
}
