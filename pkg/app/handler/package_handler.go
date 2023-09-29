package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mrparano1d/noxite/pkg/app/auth"
	"github.com/mrparano1d/noxite/pkg/core"
	"github.com/mrparano1d/noxite/pkg/core/coreerrors"
	"github.com/mrparano1d/noxite/pkg/core/services"

	json "github.com/bytedance/sonic"
)

type publishRes struct {
	OK string `json:"ok"`
}

func PackageHandler(r chi.Router, app *core.ApplicationCore) {
	r.Get("/{packageName}", func(w http.ResponseWriter, r *http.Request) {

		user := auth.GetUserFromContext(r.Context())

		packageName := chi.URLParam(r, "packageName")

		pkg, err := app.PackageService().GetPackage(r.Context(), user, packageName, "latest")
		if err != nil {
			if _, ok := err.(*services.PackageServicePackageNotFoundError); ok {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			} else if _, ok := err.(*services.PackageServiceGetPackageError); ok {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else if _, ok := err.(*coreerrors.NotAllowedToGetPackageError); ok {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}
			// TODO replace log with proper logging
			log.Println("package get failed: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(pkg)

		fmt.Println("manifest: ", string(pkg))
	})

	r.Put("/{packageName}", func(w http.ResponseWriter, r *http.Request) {

		user := auth.GetUserFromContext(r.Context())

		manifest, err := app.PackageService().ParseManifest(r.Context(), user, r.Body)
		if err != nil {
			// TODO replace log with proper logging
			log.Println("manifest parse failed", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := app.PackageService().PublishPackage(r.Context(), user, manifest); err != nil {
			// TODO replace log with proper logging
			log.Println("package publish failed: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.ConfigDefault.NewEncoder(w).Encode(publishRes{
			OK: "package " + manifest.Name.String() + " published",
		})
	})
}
