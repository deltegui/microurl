package web

import (
	"microurl/internal"
	"microurl/web/views"
	"net/http"

	"github.com/deltegui/phoenix"
	"github.com/go-chi/chi"
)

func CreatePanelRoutes(ctx Ctx) chi.Router {
	r := chi.NewRouter()
	r.Get("/", ctx.Auth.Authorize(views.RenderHandler(views.Panel, views.PanelModel{Shorten: ""})))
	r.Post("/shorten", ctx.Auth.Authorize(shortenHandler(ctx, views.Panel)))
	return r
}

func PanelPresenter(w http.ResponseWriter, req *http.Request, render views.Render) phoenix.Present {
	return func(data interface{}, err error) {
		if err != nil {
			return
		}
		payload := data.(internal.ShortenResponse)
		render(w, views.PanelModel{Shorten: payload.URL})
	}
}

func shortenHandler(ctx Ctx, render views.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		present := PanelPresenter(w, req, render)
		user, ok := req.Context().Value(UserCtxKey).(string)
		if !ok {
			http.Redirect(w, req, rootPath, http.StatusTemporaryRedirect)
			return
		}
		req.ParseForm()
		res, err := ctx.Shorten.Exec(internal.ShortenRequest{
			Name: user,
			URL:  req.Form.Get("url"),
		})
		present(res, err)
	}
}
