package httpapi

import (
	"io"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//API is the current API version
const API = "1.0"
const apiPath = "/api/" + API

func notFound(w http.ResponseWriter, r *http.Request) {
	jsonResponse(http.StatusNotFound, nil).ServeHTTP(w, r)
}

//NewRouter returns a new router for the given Server
func NewRouter(s *Server, output io.Writer) http.Handler {
	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(notFound)

	//Authenticate: POST /auth
	r.Methods("POST").Path(apiPath + "/auth").Handler(
		setAction("Authenticate",
			jsonRequest(
				s.authenticateHandler())))

	//Get: GET /urls/<id>
	r.Methods("GET").Path(apiPath + "/urls/{id:[a-zA-Z0-9]+}").Handler(
		setAction("Get",
			s.requireAuthenticated(
				s.getHandler())))

	//Put: POST /urls
	r.Methods("POST").Path(apiPath + "/urls").Handler(
		setAction("Put",
			jsonRequest(
				s.requireAuthenticated(
					s.putHandler()))))

	//Update: PUT /urls/<id>
	r.Methods("PUT").Path(apiPath + "/urls/{id:[a-zA-Z0-9]+}").Handler(
		setAction("Update",
			jsonRequest(
				s.requireAuthenticated(
					s.updateHandler()))))

	//Delete: DELETE /urls/<id>
	r.Methods("DELETE").Path(apiPath + "/urls/{id:[a-zA-Z0-9]+}").Handler(
		setAction("Delete",
			s.requireAuthenticated(
				s.deleteHandler())))

	//View: GET /<code>
	r.Methods("GET").Path("/{id:[a-zA-Z0-9]+}").Handler(
		setAction("View",
			s.viewHandler()))

	//URLs: GET /urls
	r.Methods("GET").Path(apiPath + "/urls").Handler(
		setAction("URLs",
			s.requireAuthenticated(
				s.urlsHandler())))

	return handlers.CombinedLoggingHandler(output, logRequest(output, r))
}
