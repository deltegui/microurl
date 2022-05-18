package api

import (
	"encoding/json"
	"log"
	"microurl/internal"
	"microurl/web"
	"net/http"

	"github.com/deltegui/phoenix"
	"github.com/go-chi/chi/v5"
)

func CreateURLRoutes(ctx web.Ctx) chi.Router {
	r := chi.NewRouter()
	r.Use(ctx.HeaderAuth.Authorize)
	r.Get("/all", getAllHandler(ctx))
	r.Post("/", shortenHandler(ctx))
	r.Delete("/{id}", deleteHandler(ctx))
	r.Post("/genqr/{id}", genQRHandler(ctx))
	return r
}

type CaseReq struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func shortenHandler(ctx web.Ctx) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		present := phoenix.JSONPresenter(w, req)
		user, ok := req.Context().Value(web.UserCtxKey).(string)
		if !ok {
			present(nil, internal.MalformedRequestErr)
			return
		}
		var r CaseReq
		if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
			log.Println("Cannot decode input for shorten api handler", err)
			present(nil, internal.MalformedRequestErr)
			return
		}
		present(ctx.Shorten.Exec(internal.ShortenRequest{
			Username: user,
			Name:     r.Name,
			URL:      r.URL,
		}))
	}
}

func getAllHandler(ctx web.Ctx) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		present := phoenix.JSONPresenter(w, req)
		user, ok := req.Context().Value(web.UserCtxKey).(string)
		if !ok {
			present(nil, internal.MalformedRequestErr)
			return
		}
		present(ctx.AllURLs.Exec(internal.AllURLsRequest{User: user}))
	}
}

func deleteHandler(ctx web.Ctx) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		present := phoenix.JSONPresenter(w, req)
		id, err := web.GetURLID(req)
		if err != nil {
			present(nil, internal.MalformedRequestErr)
			return
		}
		present(ctx.Delete.Exec(internal.URLIDRequest{
			URLID: id,
		}))
	}
}

func genQRHandler(ctx web.Ctx) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		present := phoenix.JSONPresenter(w, req)
		id, err := web.GetURLID(req)
		if err != nil {
			present(nil, internal.MalformedRequestErr)
			return
		}
		present(ctx.GenQR.Exec(internal.URLIDRequest{
			URLID: id,
		}))
	}
}
