package web

import (
	"microurl/internal"
	"microurl/web/session"
)

type Ctx struct {
	ListenURL string
	Session   session.Manager
	Auth      SessionJWTAuth
	Login     internal.UseCase
	Shorten   internal.UseCase
	Access    internal.UseCase
	AllURLs   internal.UseCase
	Delete    internal.UseCase
	GenQR     internal.UseCase
}
