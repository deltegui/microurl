package main

import (
	"fmt"
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
	phoenix.PrintLogo("banner")
	conf := config.Load()
	ctx := wire(conf)
	router := createRouter()
	mount(ctx, router)
	log.Println("Listening on :3000")
	// phoenix.FileServerStatic(router, "/static")
	listen(router, conf)
}

func listen(r chi.Router, config config.Configuration) {
	log.Printf("Listening on %s with tls? %t\n", config.ListenURL, config.TLS.Enabled)
	// log.Println("CORS allow origin: ", config.CORS)
	var err error
	if config.TLS.Enabled {
		err = http.ListenAndServeTLS(config.ListenURL, config.TLS.CRT, config.TLS.KEY, r)
	} else {
		err = http.ListenAndServe(config.ListenURL, r)
	}
	if err != nil {
		log.Fatal("Error listening: ", err)
	}
}

func redirectToRoot(w http.ResponseWriter, req *http.Request) phoenix.Present {
	return func(data interface{}, err error) {
		log.Printf("Error while auth: %s", err)
		http.Redirect(w, req, "/", http.StatusMovedPermanently)
	}
}

func wire(conf config.Configuration) web.Ctx {
	conn := persistence.Connect(conf)
	conn.MigrateAll()
	val := validator.New()
	userRepo := persistence.NewGormUserRepository(conn)
	urlRepo := persistence.NewGormURLRepository(conn)
	hasher := phxHash.BcryptHasher{}
	shortHasher := shortener.Base62{}
	tokenizer := token.New(conf.JWTKey)
	sessionManager := session.New(conf.SessionKey)
	genURL := func(path string) string {
		return fmt.Sprintf("%s/%s", conf.ListenURL, path)
	}
	return web.Ctx{
		Session: sessionManager,
		Auth:    web.NewSessionJWTAuth(tokenizer, redirectToRoot, sessionManager),
		Login:   internal.NewLoginCase(val, userRepo, hasher, tokenizer),
		Shorten: internal.NewShortenCase(val, userRepo, urlRepo, shortHasher, genURL),
		Access:  internal.NewAccessCase(val, urlRepo, shortHasher),
		AllURLs: internal.NewAllURLsCase(val, urlRepo, shortHasher, genURL),
		Delete:  internal.NewDeleteCase(val, urlRepo),
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
	router.Mount("/user", web.CreateUserRoutes(ctx))
	router.Mount("/", web.CreateURLRoutes(ctx))
	router.Mount("/panel", web.CreatePanelRoutes(ctx))
}
