package session

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

func sessionTransaction(req *http.Request, w http.ResponseWriter, exec func(session *sessions.Session)) {
	session, _ := store.Get(req, "session-name")
	exec(session)
	err := session.Save(req, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func StoreSession(req *http.Request, w http.ResponseWriter, token string) {
	sessionTransaction(req, w, func(session *sessions.Session) {
		session.Values["token"] = token
	})
}

func GetSession(req *http.Request, w http.ResponseWriter) (token string) {
	sessionTransaction(req, w, func(session *sessions.Session) {
		token = session.Values["token"].(string)
	})
	return
}
