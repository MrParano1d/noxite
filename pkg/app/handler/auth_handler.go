package handler

import (
	"log"
	"net/http"

	json "github.com/bytedance/sonic"
	"github.com/go-chi/chi/v5"
	"github.com/mrparano1d/noxite/pkg/core"
)

type loginReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Type     string `json:"type"`
	Email    string `json:"email"`
}

type loginRes struct {
	OK    string `json:"ok"`
	Token string `json:"token"`
}

func AuthHandler(r chi.Router, app *core.ApplicationCore) {
	r.Put("/-/user/{orgCouchDBUser}", func(w http.ResponseWriter, r *http.Request) {
		var loginReq loginReq
		err := json.ConfigDefault.NewDecoder(r.Body).Decode(&loginReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		session, err := app.AuthService().Login(r.Context(), loginReq.Name, loginReq.Email, loginReq.Password)
		if err != nil {
			// TODO replace log with proper logging
			log.Println("login failed", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)

		json.ConfigDefault.NewEncoder(w).Encode(loginRes{
			OK:    "user " + loginReq.Name + " created",
			Token: session.Token.String(),
		})

		// TODO replace log with proper logging
		log.Println("login success")
	})
}
