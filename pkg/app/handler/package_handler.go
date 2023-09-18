package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mrparano1d/getregd/pkg/core"

	json "github.com/bytedance/sonic"
)

type publishRes struct {
	OK string `json:"ok"`
}

func PackageHandler(r chi.Router, app *core.ApplicationCore) {
	r.Get("/{packageName}", func(w http.ResponseWriter, r *http.Request) {
		// TODO implement
		log.Println("package get")
		w.WriteHeader(http.StatusNotImplemented)
	})

	r.Put("/{packageName}", func(w http.ResponseWriter, r *http.Request) {

		manifest, err := app.PackageService().ParseManifest(r.Body)
		if err != nil {
			// TODO replace log with proper logging
			log.Println("manifest parse failed", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := app.PackageService().PublishPackage(manifest); err != nil {
			// TODO replace log with proper logging
			log.Println("package publish failed", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.ConfigDefault.NewEncoder(w).Encode(publishRes{
			OK: "package " + manifest.Name + " published",
		})
	})
}
