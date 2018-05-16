package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/korylprince/url-shortener-server/session"
)

func (s *Server) authenticateHandler() http.Handler {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type response struct {
		SessionID string `json:"session_id"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := new(request)
		d := json.NewDecoder(r.Body)

		if err := d.Decode(req); err != nil {
			jsonResponse(http.StatusBadRequest, fmt.Errorf("Unable to decode request body: %v", err)).ServeHTTP(w, r)
			return
		}

		ok, err := s.auth.Authenticate(req.Username, req.Password)
		if err != nil {
			jsonResponse(http.StatusInternalServerError, fmt.Errorf("Unable to authenticate: %v", err)).ServeHTTP(w, r)
			return
		}

		if !ok {
			jsonResponse(http.StatusUnauthorized, errors.New("Invalid username or password")).ServeHTTP(w, r)
			return
		}

		id, err := s.sessionStore.Create(&session.Session{Username: req.Username})
		if err != nil {
			jsonResponse(http.StatusInternalServerError, fmt.Errorf("Unable to create session: %v", err)).ServeHTTP(w, r)
			return
		}

		jsonResponse(http.StatusOK, &response{SessionID: id}).ServeHTTP(w, r)
	})
}

func (s *Server) requireAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			jsonResponse(http.StatusUnauthorized, errors.New("No Authorization header")).ServeHTTP(w, r)
			return
		}

		if !strings.HasPrefix(auth, `Session id="`) || len(auth) < 13 {
			jsonResponse(http.StatusBadRequest, errors.New("Invalid Authorization header")).ServeHTTP(w, r)
			return
		}

		id := auth[12 : len(auth)-1]

		session, err := s.sessionStore.Check(id)
		if err != nil {
			jsonResponse(http.StatusInternalServerError, fmt.Errorf("Unexpected error when checking session id %s: %v", id, err)).ServeHTTP(w, r)
			return
		}

		if session == nil {
			jsonResponse(http.StatusUnauthorized, fmt.Errorf("Session doesn't exist for id %s", id)).ServeHTTP(w, r)
			return
		}

		(r.Context().Value(contextKeyLogData)).(*logData).User = session.Username

		ctx := context.WithValue(r.Context(), contextKeyUser, session.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
