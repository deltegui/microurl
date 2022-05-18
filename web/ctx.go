package web

import (
	"microurl/internal"
	"microurl/web/session"
)

type Ctx struct {
	ListenURL   string
	Session     session.Manager
	SessionAuth SessionJWTAuth
	HeaderAuth  HeaderJWTAuth
	Login       internal.UseCase
	Shorten     internal.UseCase
	Access      internal.UseCase
	AllURLs     internal.UseCase
	Delete      internal.UseCase
	GenQR       internal.UseCase
}
