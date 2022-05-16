package web

import (
	"microurl/internal"
	"microurl/web/views"
	"net/http"
	"strconv"

	"github.com/deltegui/phoenix"
	"github.com/go-chi/chi"
)

const panelPath = "/panel"

func CreatePanelRoutes(ctx Ctx) chi.Router {
	r := chi.NewRouter()
	r.Get("/", ctx.Auth.Authorize(panelHandler(ctx, views.Panel)))
	r.Post("/shorten", ctx.Auth.Authorize(shortenHandler(ctx, views.Panel)))
	r.Post("/delete/{id}", ctx.Auth.Authorize(deleteHandler(ctx, views.Panel)))
	return r
}

func PanelPresenter(w http.ResponseWriter, req *http.Request, render views.Render) phoenix.Present {
	return func(data interface{}, err error) {
		if err != nil {
			return
		}
		payload := data.(internal.AllURLsResponse)
		render(w, payload)
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
		_, err := ctx.Shorten.Exec(internal.ShortenRequest{
			Name: user,
			URL:  req.Form.Get("url"),
		})
		if err != nil {
			present(nil, internal.MalformedRequestErr)
		}
		panelHandler(ctx, render)(w, req)
	}
}

func panelHandler(ctx Ctx, render views.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		present := PanelPresenter(w, req, render)
		user, ok := req.Context().Value(UserCtxKey).(string)
		if !ok {
			http.Redirect(w, req, rootPath, http.StatusTemporaryRedirect)
			return
		}
		res, err := ctx.AllURLs.Exec(internal.AllURLsRequest{
			User: user,
		})
		present(res, err)
	}
}

func deleteHandler(ctx Ctx, render views.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		present := PanelPresenter(w, req, render)
		str := chi.URLParam(req, "id")
		id, err := strconv.Atoi(str)
		if err != nil {
			present(nil, internal.MalformedRequestErr)
		}
		_, err = ctx.Delete.Exec(internal.DeleteRequest{
			URLID: uint(id),
		})
		if err != nil {
			present(nil, internal.MalformedRequestErr)
		}
		panelHandler(ctx, render)(w, req)
	}
}
