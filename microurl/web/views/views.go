package views

import (
	"embed"
	"html/template"
	"io"
	"microurl/internal"
	"net/http"
	"strings"
)

//go:embed *
var files embed.FS

var funcs = template.FuncMap{
	"uppercase": func(v string) string {
		return strings.ToUpper(v)
	},
}

func parse(file string) *template.Template {
	return template.Must(
		template.
			New("layout.html").
			Funcs(funcs).
			ParseFS(files, "layout.html", file))
}

func RedirectHandler(path string, statusCode int) http.HandlerFunc {
	return http.RedirectHandler(path, statusCode).ServeHTTP
}

func RenderHandler(render Render, m interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		render(w, m)
	}
}

type Render func(w io.Writer, m interface{}) error

var (
	login  = parse("login.html")
	panel  = parse("panel.html")
	urlErr = parse("urlerr.html")
)

type LoginModel struct {
	HadError bool
	Error    string
}

func Login(w io.Writer, m interface{}) error {
	return login.Execute(w, m)
}

type PanelModel struct {
	HadError bool
	Error    string
	URLs     []internal.URLResponse
}

func Panel(w io.Writer, m interface{}) error {
	return panel.Execute(w, m)
}

func URLError(w io.Writer, m interface{}) error {
	return urlErr.Execute(w, m)
}
