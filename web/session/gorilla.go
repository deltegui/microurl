package session

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

type Manager struct {
	store *sessions.CookieStore
}

func New(key string) Manager {
	return Manager{sessions.NewCookieStore([]byte(key))}
}

func (manager Manager) sessionTransaction(req *http.Request, w http.ResponseWriter, exec func(session *sessions.Session)) {
	session, _ := manager.store.Get(req, "session-name")
	exec(session)
	err := session.Save(req, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (manager Manager) Store(w http.ResponseWriter, req *http.Request, token string) {
	manager.sessionTransaction(req, w, func(session *sessions.Session) {
		session.Values["token"] = token
	})
}

func (manager Manager) Get(w http.ResponseWriter, req *http.Request) (token string, err error) {
	manager.sessionTransaction(req, w, func(session *sessions.Session) {
		if _, ok := session.Values["token"]; !ok {
			err = fmt.Errorf("session not found")
			return
		}
		err = nil
		token = session.Values["token"].(string)
	})
	return
}

func (manager Manager) Reset(w http.ResponseWriter, req *http.Request) {
	manager.sessionTransaction(req, w, func(session *sessions.Session) {
		delete(session.Values, "token")
	})
}
