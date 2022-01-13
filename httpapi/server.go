package httpapi

import (
	"io"
	"io/fs"

	"github.com/korylprince/httputil/auth"
	"github.com/korylprince/httputil/session"
	"github.com/korylprince/url-shortener-server/v2/db"
)

//Server represents shared resources
type Server struct {
	AppTitle     string
	db           db.DB
	auth         auth.Auth
	adminGroup   string
	sessionStore session.Store
	files        fs.FS
	output       io.Writer
}

//NewServer returns a new server with the given resources
func NewServer(title string, db db.DB, auth auth.Auth, adminGroup string, sessionStore session.Store, files fs.FS, output io.Writer) *Server {
	return &Server{AppTitle: title, db: db, auth: auth, adminGroup: adminGroup, sessionStore: sessionStore, files: files, output: output}
}
