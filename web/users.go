package web

import (
	"microurl/internal"
	"net/http"

	phx "github.com/deltegui/phoenix"
)

func GenericMapper(req *http.Request) interface{} {
	return struct{}{}
}

func LoginViewModelMapper(req *http.Request) interface{} {
	return LoginViewModel{HadError: false}
}

var LoginViewConfig phx.ViewConfig = phx.ViewConfig{
	Layout: "/base.html",
	View:   "/home/index.html",
	Name:   "base.html",
}

var PanelViewConfig phx.ViewConfig = phx.ViewConfig{
	Layout: "/base.html",
	View:   "/home/panel.html",
	Name:   "base.html",
}

type LoginViewModel struct {
	Name         string
	HadError     bool
	ErrorMessage string
}

type loginPresenter struct {
	login phx.HTMLRenderer
	w     http.ResponseWriter
	r     *http.Request
}

func (presenter loginPresenter) Present(data interface{}) {
	presenter.login.Redirect(
		presenter.w,
		presenter.r,
		"/panel",
		http.StatusMovedPermanently)
}

func (presenter loginPresenter) PresentError(err error) {
	presenter.login.Render(presenter.w, LoginViewModel{
		Name:         "",
		HadError:     true,
		ErrorMessage: err.Error(),
	})
}

func LoginHandler(loginCase internal.UseCase) http.HandlerFunc {
	loginRenderer := phx.NewHTMLRenderer(LoginViewConfig)
	return func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		presenter := loginPresenter{
			login: loginRenderer,
			w:     w,
			r:     req,
		}
		loginCase.Exec(presenter, internal.LoginRequest{
			Name:     req.Form.Get("name"),
			Password: req.Form.Get("password"),
		})
	}
}
