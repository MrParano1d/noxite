package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mrparano1d/getregd/pkg/app/handler"
	"github.com/mrparano1d/getregd/pkg/core"

	"log"
	"net/http"
)

func ServeApp(app *core.ApplicationCore) error {

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	handler.AuthHandler(r, app)
	handler.PackageHandler(r, app)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})

	log.Println("server started at http://localhost:3000")
	return http.ListenAndServe(":3000", r)
}
