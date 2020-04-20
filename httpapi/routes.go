package httpapi

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/korylprince/httputil/auth/ad"
	"github.com/korylprince/httputil/jsonapi"
	"github.com/korylprince/httputil/session"
)

const allowedIDRegexp = "[a-zA-Z0-9_\\-.]+"

//API is the current API version
const API = "1.1"
const apiPath = "/api/" + API

//Router returns a new router
func (s *Server) Router() http.Handler {
	r := mux.NewRouter()

	var hook = func(sess session.Session) (bool, interface{}, error) {
		attrs := map[string]bool{"admin": false}
		user := sess.(*ad.User)
		for _, g := range user.Groups {
			if g == s.adminGroup {
				attrs["admin"] = true
				break
			}
		}

		return true, attrs, nil
	}

	apirouter := jsonapi.New(s.output, s.auth, s.sessionStore, hook)
	r.PathPrefix(apiPath).Handler(http.StripPrefix(apiPath, apirouter))

	apirouter.Handle("GET", fmt.Sprintf("/urls/{id:%s}", allowedIDRegexp), s.getHandler, true)
	apirouter.Handle("POST", "/urls", s.putHandler, true)
	apirouter.Handle("PUT", fmt.Sprintf("/urls/{id:%s}", allowedIDRegexp), s.updateHandler, true)
	apirouter.Handle("DELETE", fmt.Sprintf("/urls/{id:%s}", allowedIDRegexp), s.deleteHandler, true)
	apirouter.Handle("GET", "/title", s.titleHandler, false)
	apirouter.Handle("GET", "/urls", s.urlsHandler, true)

	r.Path("/error.html").Handler(http.FileServer(s.box))
	r.Methods("GET").Path(fmt.Sprintf("/{id:%s}", allowedIDRegexp)).Handler(withRedirect(s.viewHandler))
	r.PathPrefix("/").Handler(http.FileServer(s.box))

	return handlers.CombinedLoggingHandler(s.output, r)
}
