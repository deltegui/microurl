package web

import (
	"microurl/internal"
	"microurl/web/session"
)

type Ctx struct {
	Session session.Manager
	Auth    SessionJWTAuth
	Login   internal.UseCase
	Shorten internal.UseCase
}
