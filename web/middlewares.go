package web

import (
	"context"
	"microurl/internal"
	"net/http"
	"strings"
	"time"

	"github.com/deltegui/phoenix"
)

type ctxKey string

const UserCtxKey ctxKey = "user"

// OptionsCors creates a middleware that handles all request using
// Options method and returns 204. This is used to return all
// CORS headers.
func OptionsCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, req)
	})
}

// JWTAuth middleware to check using a JWT Bearer token if user is authorized or is admin.
type JWTAuth struct {
	tokenizer internal.Tokenizer
}

// NewJWTAuth create JWTAuth middleware.
func NewJWTAuth(tokenizer internal.Tokenizer) JWTAuth {
	return JWTAuth{tokenizer}
}

// Authorize middleware that checks if exists the header 'Authorization' with
// valid JWT bearer token.
func (authMiddle JWTAuth) Authorize(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authMiddle.handleAndCheckToken(w, req, next, func(role temsys.Role) error { return nil })
	})
}

func (authMiddle JWTAuth) handleAndCheckToken(w http.ResponseWriter, req *http.Request, next http.Handler) {
	presenter := phoenix.NewJSONPresenter(w)
	token, err := authMiddle.getToken(req)
	if err != nil {
		presenter.PresentError(err)
		return
	}
	if token.Expires.Before(time.Now()) {
		presenter.PresentError(internal.ExpiredTokenErr)
		return
	}
	ctx := context.WithValue(req.Context(), UserCtxKey, token.Owner)
	next.ServeHTTP(w, req.WithContext(ctx))
}

func (authMiddle JWTAuth) getToken(req *http.Request) (internal.Token, error) {
	const bearerPrefix string = "Bearer "
	bearerToken := req.Header.Get("Authorization")
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
