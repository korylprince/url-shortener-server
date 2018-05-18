package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	neturl "net/url"

	"github.com/gorilla/mux"
	"github.com/korylprince/url-shortener-server/db"
)

func (s *Server) getHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		url, err := s.db.Get(id)
		if err != nil {
			jsonResponse(http.StatusInternalServerError, fmt.Errorf("Unable to get URL %s: %v", id, err)).ServeHTTP(w, r)
			return
		}

		if url == nil {
			jsonResponse(http.StatusNotFound, fmt.Errorf("URL %s does not exist", id)).ServeHTTP(w, r)
			return
		}

		jsonResponse(http.StatusOK, url).ServeHTTP(w, r)
	})
}

func (s *Server) putHandler() http.Handler {
	type response struct {
		URLID string `json:"url_id"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := new(db.URL)
		d := json.NewDecoder(r.Body)

		if err := d.Decode(url); err != nil {
			jsonResponse(http.StatusBadRequest, fmt.Errorf("Unable to decode request body: %v", err)).ServeHTTP(w, r)
			return
		}

		(r.Context().Value(contextKeyLogData)).(*logData).Data = url

		if _, err := neturl.ParseRequestURI(url.URL); err != nil {
			jsonResponse(http.StatusBadRequest, fmt.Errorf(`Unable to parse url "%s": %v`, url.URL, err)).ServeHTTP(w, r)
			return
		}

		user := (r.Context().Value(contextKeyUser)).(string)

		id, err := s.db.Put(url, user)
		if err != nil {
			jsonResponse(http.StatusInternalServerError, fmt.Errorf(`Unable to put URL "%s": %v`, url.URL, err)).ServeHTTP(w, r)
			return
		}

		(r.Context().Value(contextKeyLogData)).(*logData).URLID = id

		jsonResponse(http.StatusOK, &response{URLID: id}).ServeHTTP(w, r)
	})
}

func (s *Server) updateHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		user := (r.Context().Value(contextKeyUser)).(string)

		//read url
		url := new(db.URL)
		d := json.NewDecoder(r.Body)

		if err := d.Decode(url); err != nil {
			jsonResponse(http.StatusBadRequest, fmt.Errorf("Unable to decode request body: %v", err)).ServeHTTP(w, r)
			return
		}

		(r.Context().Value(contextKeyLogData)).(*logData).Data = url

		if _, err := neturl.ParseRequestURI(url.URL); err != nil {
			jsonResponse(http.StatusBadRequest, fmt.Errorf(`Unable to parse url "%s": %v`, url.URL, err)).ServeHTTP(w, r)
			return
		}

		//check that user owns URL
		urls, err := s.db.URLs(user)
		if err != nil {
			jsonResponse(http.StatusInternalServerError, fmt.Errorf(`Unable to get URLs for user %s: %v`, user, err)).ServeHTTP(w, r)
			return
		}

		owned := false
		for _, url := range urls {
			if url.ID == id {
				owned = true
			}
		}

		if !owned {
			jsonResponse(http.StatusForbidden, fmt.Errorf("URL %s is not owned by user %s", id, user)).ServeHTTP(w, r)
			return
		}

		//update url
		if err = s.db.Update(id, url); err != nil {
			jsonResponse(http.StatusInternalServerError, fmt.Errorf(`Unable to update URL %s: %v`, id, err)).ServeHTTP(w, r)
			return
		}

		//read url
		url, err = s.db.Get(id)
		if err != nil || url == nil {
			if err == nil {
				err = errors.New("URL unexpectedly nil")
			}
			jsonResponse(http.StatusInternalServerError, fmt.Errorf("Unable to get URL %s: %v", id, err)).ServeHTTP(w, r)
			return
		}

		jsonResponse(http.StatusOK, url).ServeHTTP(w, r)
	})
}

func (s *Server) deleteHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		user := (r.Context().Value(contextKeyUser)).(string)

		//check that user owns URL
		urls, err := s.db.URLs(user)
		if err != nil {
			jsonResponse(http.StatusInternalServerError, fmt.Errorf(`Unable to get URLs for user %s: %v`, user, err)).ServeHTTP(w, r)
			return
		}

		owned := false
		for _, url := range urls {
			if url.ID == id {
				owned = true
			}
		}

		if !owned {
			jsonResponse(http.StatusForbidden, fmt.Errorf("URL %s is not owned by user %s", id, user)).ServeHTTP(w, r)
			return
		}

		//delete url
		if err := s.db.Delete(id); err != nil {
			jsonResponse(http.StatusInternalServerError, fmt.Errorf(`Unable to delete URL %s: %v`, id, err)).ServeHTTP(w, r)
			return
		}

		jsonResponse(http.StatusOK, nil).ServeHTTP(w, r)
	})
}

func (s *Server) viewHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		url, err := s.db.View(id)
		if err != nil {
			jsonResponse(http.StatusInternalServerError, fmt.Errorf("Unable to get URL %s: %v", id, err)).ServeHTTP(w, r)
			return
		}

		if url == "" {
			jsonResponse(http.StatusNotFound, fmt.Errorf("URL %s does not exist", id)).ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})
}

func (s *Server) urlsHandler() http.Handler {
	type response struct {
		URLs []*db.URL `json:"urls"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := (r.Context().Value(contextKeyUser)).(string)

		urls, err := s.db.URLs(user)
		if err != nil {
			jsonResponse(http.StatusInternalServerError, fmt.Errorf(`Unable to get URLs for user %s: %v`, user, err)).ServeHTTP(w, r)
			return
		}

		jsonResponse(http.StatusOK, &response{URLs: urls}).ServeHTTP(w, r)
	})
}
