package main

import (
	"log"
	"microurl/internal"
	"microurl/internal/config"
	"microurl/internal/persistence"
	"microurl/internal/shortener"
	"microurl/internal/token"
	"microurl/web"
	"microurl/web/session"
	"net/http"
	"time"

	"github.com/deltegui/phoenix"
	phxHash "github.com/deltegui/phoenix/hash"
	"github.com/deltegui/phoenix/validator"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	ctx := wire()
	router := createRouter()
	mount(ctx, router)
	log.Println("Listening on :3000")
	if err := http.ListenAndServe(":3000", router); err != nil {
		log.Fatalln("Error while creating server")
	}
}

func redirectToRoot(w http.ResponseWriter, req *http.Request) phoenix.Present {
	return func(data interface{}, err error) {
		log.Printf("Error while auth: %s", err)
		http.Redirect(w, req, "/", http.StatusMovedPermanently)
	}
}

func wire() web.Ctx {
	conf := config.Load()
	conn := persistence.Connect(conf)
	conn.MigrateAll()
	val := validator.New()
	userRepo := persistence.NewGormUserRepository(conn)
	urlRepo := persistence.NewGormURLRepository(conn)
	hasher := phxHash.BcryptHasher{}
	shortHasher := shortener.Base62{}
	tokenizer := token.New(conf.JWTKey)
	sessionManager := session.New(conf.SessionKey)
	return web.Ctx{
		Session: sessionManager,
		Auth:    web.NewSessionJWTAuth(tokenizer, redirectToRoot, sessionManager),
		Login:   internal.NewLoginCase(val, userRepo, hasher, tokenizer),
		Shorten: internal.NewShortenCase(val, userRepo, urlRepo, shortHasher),
	}
}

func createRouter() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))
	return router
}

func mount(ctx web.Ctx, router chi.Router) {
	router.Mount("/", web.CreateUserRoutes(ctx))
	router.Mount("/panel", web.CreatePanelRoutes(ctx))
}
