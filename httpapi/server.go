package httpapi

import (
	"github.com/gobuffalo/packr"
	"github.com/korylprince/url-shortener-server/auth"
	"github.com/korylprince/url-shortener-server/db"
	"github.com/korylprince/url-shortener-server/session"
)

//Server represents shared resources
type Server struct {
	AppTitle     string
	db           db.DB
	auth         auth.Auth
	sessionStore session.Store
	box          packr.Box
}

//NewServer returns a new server with the given resources
func NewServer(title string, db db.DB, auth auth.Auth, sessionStore session.Store, box packr.Box) *Server {
	return &Server{AppTitle: title, db: db, auth: auth, sessionStore: sessionStore, box: box}
}
