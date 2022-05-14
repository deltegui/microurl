package main

import (
	"log"
	"microurl/internal"
	"microurl/internal/config"
	"microurl/internal/persistence"
	"microurl/internal/token"
	"microurl/web"
	"net/http"

	phx "github.com/deltegui/phoenix"
	"github.com/deltegui/phoenix/hash"
	"github.com/deltegui/phoenix/validator"
	"github.com/go-chi/chi"
)

func main() {
	conf := config.Load()

	conn := persistence.Connect(conf)
	conn.MigrateAll()
	validator := validator.New()
	hasher := hash.BcryptPasswordHasher{}
	tokenizer := token.New(conf.JWTKey)
	auth := web.NewJWTAuth(tokenizer)

	router := chi.NewRouter()
	router.Get("/login", phx.RenderView(web.LoginViewConfig, web.LoginViewModelMapper))

	router.Post("/login", web.LoginHandler(internal.NewLoginCase(
		validator,
		persistence.NewGormUserRepository(conn),
		hasher,
		tokenizer)))
	router.Get("/panel", auth.Authorize(phx.RenderView(web.PanelViewConfig, web.GenericMapper)))

	log.Println("Listening on", conf.ListenURL)
	if err := http.ListenAndServe(conf.ListenURL, router); err != nil {
		log.Fatalf("Error while listening: %s: %s\n", conf.ListenURL, err)
	}
}
