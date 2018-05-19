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

	api := r.PathPrefix(apiPath).Subrouter()

	api.NotFoundHandler = http.HandlerFunc(notFound)

	//Authenticate: POST /auth
	api.Methods("POST").Path("/auth").Handler(
		logRequest(output,
			setAction("Authenticate",
				jsonRequest(
					s.authenticateHandler()))))

	//Get: GET /urls/<id>
	api.Methods("GET").Path("/urls/{id:[a-zA-Z0-9]+}").Handler(
		logRequest(output,
			setAction("Get",
				s.requireAuthenticated(
					s.getHandler()))))

	//Put: POST /urls
	api.Methods("POST").Path("/urls").Handler(
		logRequest(output,
			setAction("Put",
				jsonRequest(
					s.requireAuthenticated(
						s.putHandler())))))

	//Update: PUT /urls/<id>
	api.Methods("PUT").Path("/urls/{id:[a-zA-Z0-9]+}").Handler(
		logRequest(output,
			setAction("Update",
				jsonRequest(
					s.requireAuthenticated(
						s.updateHandler())))))

	//Delete: DELETE /urls/<id>
	api.Methods("DELETE").Path("/urls/{id:[a-zA-Z0-9]+}").Handler(
		logRequest(output,
			setAction("Delete",
				s.requireAuthenticated(
					s.deleteHandler()))))

	//URLs: GET /urls
	api.Methods("GET").Path("/urls").Handler(
		logRequest(output,
			setAction("URLs",
				s.requireAuthenticated(
					s.urlsHandler()))))

	//View: GET /<code>
	r.Methods("GET").Path("/{id:[a-zA-Z0-9]+}").Handler(
		logRequest(output,
			setAction("View",
				s.viewHandler())))

	r.PathPrefix("/").Handler(http.FileServer(s.box))

	return handlers.CombinedLoggingHandler(output, r)
}
