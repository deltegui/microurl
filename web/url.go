package web

import (
	"microurl/internal"
	"microurl/web/views"
	"net/http"

	"github.com/deltegui/phoenix"
	"github.com/go-chi/chi"
)

const (
	rootPath = "/"
)

func CreateURLRoutes(ctx Ctx) chi.Router {
	r := chi.NewRouter()
	r.Get(rootPath, views.RedirectHandler(loginPath, http.StatusMovedPermanently))
	r.Get("/{id}", urlHandler(ctx, views.URLError))
	return r
}

func URLPresenter(w http.ResponseWriter, req *http.Request, render views.Render) phoenix.Present {
	return func(data interface{}, err error) {
		if err != nil {
			render(w, struct{}{})
			return
		}
		payload := data.(internal.URLResponse)
		http.Redirect(w, req, payload.Original, http.StatusMovedPermanently)
	}
}

func urlHandler(ctx Ctx, render views.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		present := URLPresenter(w, req, render)
		res, err := ctx.Access.Exec(internal.AccessRequest{
			Hash: chi.URLParam(req, "id"),
		})
		present(res, err)
	}
}
