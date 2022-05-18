package api

import (
	"encoding/json"
	"microurl/internal"
	"microurl/web"
	"net/http"

	"github.com/deltegui/phoenix"
	"github.com/go-chi/chi/v5"
)

func CreateUserRoutes(ctx web.Ctx) chi.Router {
	r := chi.NewRouter()
	r.Post("/login", loginHandler(ctx))
	return r
}

func loginHandler(ctx web.Ctx) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		present := phoenix.JSONPresenter(w, req)
		var loginReq internal.LoginRequest
		if err := json.NewDecoder(req.Body).Decode(&loginReq); err != nil {
			present(nil, internal.MalformedRequestErr)
			return
		}
		present(ctx.Login.Exec(loginReq))
	}
}
