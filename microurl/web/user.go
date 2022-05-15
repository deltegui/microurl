package web

import (
	"microurl/internal"
	"microurl/web/views"
	"net/http"

	"github.com/deltegui/phoenix"
	"github.com/go-chi/chi"
)

const (
	rootPath   = "/"
	loginPath  = "/login"
	logoutPath = "/logout"
)

func CreateUserRoutes(ctx Ctx) chi.Router {
	r := chi.NewRouter()
	r.Get(rootPath, views.RedirectHandler(loginPath, http.StatusMovedPermanently))
	r.Get(loginPath, views.RenderHandler(views.Login, views.LoginModel{HadError: false}))
	r.Post(loginPath, loginHandler(ctx, views.Login))
	r.Get(logoutPath, logoutHandler(ctx))
	r.Post(logoutPath, logoutHandler(ctx))
	return r
}

func LoginPresenter(w http.ResponseWriter, req *http.Request, ctx Ctx, render views.Render) phoenix.Present {
	return func(data interface{}, err error) {
		if err == nil {
			payload := data.(internal.LoginResponse)
			ctx.Session.Store(w, req, payload.Token.Value)
			http.Redirect(w, req, "/panel", http.StatusMovedPermanently)
			return
		}
		caseErr := err.(internal.UseCaseError)
		render(w, views.LoginModel{
			HadError: true,
			Error:    caseErr.Reason,
		})
	}
}

func loginHandler(ctx Ctx, render views.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		present := LoginPresenter(w, req, ctx, render)
		if _, err := ctx.Session.Get(w, req); err == nil {
			http.Redirect(w, req, "/panel", http.StatusTemporaryRedirect)
			return
		}
		req.ParseForm()
		res, err := ctx.Login.Exec(internal.LoginRequest{
			Name:     req.Form.Get("name"),
			Password: req.Form.Get("password"),
		})
		present(res, err)
	}
}

func logoutHandler(ctx Ctx) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx.Session.Reset(w, req)
		http.Redirect(w, req, rootPath, http.StatusTemporaryRedirect)
	}
}
