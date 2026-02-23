package app

import (
	"net/http"

	"jobfinder/internal/config"
	"jobfinder/internal/database"

	"github.com/go-chi/chi/v5"
)

type Application struct {
	Config *config.Config
}

func New() (*Application, http.Handler) {
	
	cfg := config.Load()

	_ = database.NewPostgresPool(cfg.DBURL)

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Api is RUNNING!"))
	})

	return &Application{Config: cfg}, r
}