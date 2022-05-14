package web

import (
	"microurl/internal"
	"microurl/web/session"
	"net/http"
	"strings"

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
	res := data.(internal.LoginResponse)
	session.StoreSession(presenter.w, presenter.r, res.Token.Value)
	presenter.login.Redirect(
		presenter.w,
		presenter.r,
		"/panel",
		http.StatusMovedPermanently)
}

func (presenter loginPresenter) PresentError(err error) {
	caseErr := err.(internal.UseCaseError)
	msg := strings.ReplaceAll(caseErr.Reason, "\n", "<br />")
	presenter.login.Render(presenter.w, LoginViewModel{
		Name:         "",
		HadError:     true,
		ErrorMessage: msg,
	})
}

func LoginHandler(loginCase internal.UseCase) http.HandlerFunc {
	loginRenderer := phx.NewHTMLRenderer(LoginViewConfig)
	return func(w http.ResponseWriter, req *http.Request) {
		presenter := loginPresenter{
			login: loginRenderer,
			w:     w,
			r:     req,
		}
		if _, err := session.GetSession(w, req); err == nil {
			presenter.login.Redirect(
				presenter.w,
				presenter.r,
				"/panel",
				http.StatusTemporaryRedirect)
			return
		}
		req.ParseForm()
		loginCase.Exec(presenter, internal.LoginRequest{
			Name:     req.Form.Get("name"),
			Password: req.Form.Get("password"),
		})
	}
}
