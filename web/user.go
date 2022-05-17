package web

import (
	"microurl/internal"
	"microurl/web/views"
	"net/http"

	"github.com/deltegui/phoenix"
	"github.com/go-chi/chi/v5"
)

const (
	loginPath  = "/user/login"
	logoutPath = "/user/logout"
)

func CreateUserRoutes(ctx Ctx) chi.Router {
	r := chi.NewRouter()
	r.Get("/login", showLoginHander(ctx, views.Login))
	r.Post("/login", loginHandler(ctx, views.Login))
	r.Get("/logout", logoutHandler(ctx))
	r.Post("/logout", logoutHandler(ctx))
	return r
}

func LoginPresenter(w http.ResponseWriter, req *http.Request, ctx Ctx, render views.Render) phoenix.Present {
	return func(data interface{}, err error) {
		if err == nil {
			payload := data.(internal.LoginResponse)
			ctx.Session.Store(w, req, payload.Token.Value)
			http.Redirect(w, req, panelPath, http.StatusMovedPermanently)
			return
		}
		caseErr := err.(internal.UseCaseError)
		render(w, views.LoginModel{
			HadError: true,
			Error:    caseErr.Reason,
		})
	}
}

func showLoginHander(ctx Ctx, render views.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		redirectIfNotLogged(w, req, ctx)
		render(w, views.LoginModel{
			HadError: false,
		})
	}
}

func loginHandler(ctx Ctx, render views.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		present := LoginPresenter(w, req, ctx, render)
		redirectIfNotLogged(w, req, ctx)
		req.ParseForm()
		res, err := ctx.Login.Exec(internal.LoginRequest{
			Name:     req.Form.Get("name"),
			Password: req.Form.Get("password"),
		})
		present(res, err)
	}
}

func redirectIfNotLogged(w http.ResponseWriter, req *http.Request, ctx Ctx) {
	if _, err := ctx.Session.Get(w, req); err == nil {
		http.Redirect(w, req, panelPath, http.StatusTemporaryRedirect)
		return
	}
}

func logoutHandler(ctx Ctx) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx.Session.Reset(w, req)
		http.Redirect(w, req, rootPath, http.StatusTemporaryRedirect)
	}
}
