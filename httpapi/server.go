package httpapi

import (
	"github.com/korylprince/url-shortener-server/auth"
	"github.com/korylprince/url-shortener-server/db"
	"github.com/korylprince/url-shortener-server/session"
)

//Server represents shared resources
type Server struct {
	db           db.DB
	auth         auth.Auth
	sessionStore session.Store
}

//NewServer returns a new server with the given resources
func NewServer(db db.DB, auth auth.Auth, sessionStore session.Store) *Server {
	return &Server{db: db, auth: auth, sessionStore: sessionStore}
}
