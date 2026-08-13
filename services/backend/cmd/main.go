package main

import (
	"net/http"

	"github.com/Aneeshie/nara-code/internal/config"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	http.ListenAndServe(":"+cfg.PORT, r)
}
