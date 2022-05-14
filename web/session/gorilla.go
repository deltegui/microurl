package session

import (
	"fmt"
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

func StoreSession(w http.ResponseWriter, req *http.Request, token string) {
	sessionTransaction(req, w, func(session *sessions.Session) {
		session.Values["token"] = token
	})
}

func GetSession(w http.ResponseWriter, req *http.Request) (token string, err error) {
	sessionTransaction(req, w, func(session *sessions.Session) {
		if _, ok := session.Values["token"]; !ok {
			err = fmt.Errorf("session not found")
			return
		}
		err = nil
		token = session.Values["token"].(string)
	})
	return
}
