package web

import (
	"fmt"
	"microurl/internal"
	"microurl/web/views"
	"net/http"
	"strconv"

	"github.com/deltegui/phoenix"
	"github.com/go-chi/chi/v5"
)

const panelPath = "/panel"

func CreatePanelRoutes(ctx Ctx) chi.Router {
	r := chi.NewRouter()
	r.Use(ctx.Auth.Authorize)
	r.Get("/", panelHandler(ctx, views.Panel))
	r.Post("/shorten", shortenHandler(ctx, views.Panel))
	r.Post("/delete/{id}", deleteHandler(ctx, views.Panel))
	r.Post("/genqr/{id}", genQRHandler(ctx, views.Panel))
	return r
}

func PanelPresenter(w http.ResponseWriter, req *http.Request, render views.Render) phoenix.Present {
	return func(data interface{}, err error) {
		if err != nil {
			urls := []internal.URLResponse{}
			if payload, ok := data.(internal.AllURLsResponse); ok {
				urls = payload.URLs
			}
			render(w, views.PanelModel{
				URLs:     urls,
				HadError: true,
				Error:    err.(internal.UseCaseError).Reason,
			})
			return
		}
		payload := data.(internal.AllURLsResponse)
		render(w, views.PanelModel{
			URLs:     payload.URLs,
			HadError: false,
		})
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
			Username: user,
			Name:     req.Form.Get("name"),
			URL:      req.Form.Get("url"),
		})
		res, errAll := getAllUrls(w, req, ctx)
		if errAll != nil {
			present(res, errAll)
			return
		}
		if err != nil {
			present(res, internal.MalformedRequestErr)
			return
		}
		present(res, nil)
	}
}

func panelHandler(ctx Ctx, render views.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		present := PanelPresenter(w, req, render)
		res, err := getAllUrls(w, req, ctx)
		present(res, err)
	}
}

func getAllUrls(w http.ResponseWriter, req *http.Request, ctx Ctx) (internal.UseCaseResponse, error) {
	user, ok := req.Context().Value(UserCtxKey).(string)
	if !ok {
		http.Redirect(w, req, rootPath, http.StatusTemporaryRedirect)
		return internal.NoResponse, fmt.Errorf("not logged")
	}
	return ctx.AllURLs.Exec(internal.AllURLsRequest{
		User: user,
	})
}

func deleteHandler(ctx Ctx, render views.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		present := PanelPresenter(w, req, render)
		id, err := getURLID(req)
		if err != nil {
			present(nil, internal.MalformedRequestErr)
			return
		}
		_, err = ctx.Delete.Exec(internal.URLIDRequest{
			URLID: id,
		})
		if err != nil {
			present(nil, internal.MalformedRequestErr)
			return
		}
		panelHandler(ctx, render)(w, req)
	}
}

func genQRHandler(ctx Ctx, render views.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		present := PanelPresenter(w, req, render)
		id, err := getURLID(req)
		if err != nil {
			present(nil, internal.MalformedRequestErr)
			return
		}
		ctx.GenQR.Exec(internal.URLIDRequest{
			URLID: id,
		})
		panelHandler(ctx, render)(w, req)
	}
}

func getURLID(req *http.Request) (uint, error) {
	str := chi.URLParam(req, "id")
	id, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
