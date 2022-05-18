package web

import (
	"context"
	"microurl/internal"
	"microurl/web/session"
	"net/http"
	"strings"
	"time"

	"github.com/deltegui/phoenix"
)

type ctxKey string

const UserCtxKey ctxKey = "user"

type tokenGetter func(w http.ResponseWriter, req *http.Request) (internal.Token, error)

type PresenterCreator func(w http.ResponseWriter, req *http.Request) phoenix.Present

type JWTAuth struct {
	tokenizer       internal.Tokenizer
	createPresenter func(w http.ResponseWriter, req *http.Request) phoenix.Present
}

func (authMiddle JWTAuth) createHandler(next http.Handler, getToken tokenGetter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authMiddle.handleAndCheckToken(w, req, next, getToken)
	})
}

func (authMiddle JWTAuth) handleAndCheckToken(w http.ResponseWriter, req *http.Request, next http.Handler, getToken tokenGetter) {
	present := authMiddle.createPresenter(w, req)
	token, err := getToken(w, req)
	if err != nil {
		present(nil, err)
		return
	}
	if token.Expires.Before(time.Now()) {
		present(nil, internal.ExpiredTokenErr)
		return
	}
	ctx := context.WithValue(req.Context(), UserCtxKey, token.Owner)
	next.ServeHTTP(w, req.WithContext(ctx))
}

type SessionJWTAuth struct {
	JWTAuth
	session session.Manager
}

func NewSessionJWTAuth(tokenizer internal.Tokenizer, createPresenter PresenterCreator, manager session.Manager) SessionJWTAuth {
	return SessionJWTAuth{
		JWTAuth{tokenizer, createPresenter},
		manager,
	}
}

func (authMiddle SessionJWTAuth) Authorize(next http.Handler) http.Handler {
	return authMiddle.createHandler(next, authMiddle.getToken)
}

func (authMiddle SessionJWTAuth) getToken(w http.ResponseWriter, req *http.Request) (internal.Token, error) {
	rawToken, err := authMiddle.session.Get(w, req)
	if err != nil {
		return internal.Token{}, internal.NotAuthErr
	}
	token, err := authMiddle.tokenizer.Decode(rawToken)
	if err != nil {
		return internal.Token{}, internal.InvalidTokenErr
	}
	return token, nil
}

type HeaderJWTAuth struct {
	JWTAuth
}

func NewHeaderJWTAuth(tokenizer internal.Tokenizer, createPresenter PresenterCreator) HeaderJWTAuth {
	return HeaderJWTAuth{JWTAuth{tokenizer, createPresenter}}
}

func (authMiddle HeaderJWTAuth) Authorize(next http.Handler) http.Handler {
	return authMiddle.createHandler(next, authMiddle.getToken)
}

func (authMiddle HeaderJWTAuth) getToken(w http.ResponseWriter, req *http.Request) (internal.Token, error) {
	bearerToken := req.Header.Get("Authorization")
	const bearerPrefix string = "Bearer "
	if len(bearerToken) == 0 || !strings.HasPrefix(bearerToken, bearerPrefix) {
		return internal.Token{}, internal.NotAuthErr
	}
	rawToken := strings.Replace(bearerToken, bearerPrefix, "", 1)
	token, err := authMiddle.tokenizer.Decode(rawToken)
	if err != nil {
		return internal.Token{}, internal.InvalidTokenErr
	}
	return token, nil

}
