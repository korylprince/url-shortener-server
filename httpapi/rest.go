package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	neturl "net/url"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/korylprince/httputil/auth/ad"
	"github.com/korylprince/httputil/jsonapi"
	"github.com/korylprince/url-shortener-server/v2/db"
)

func (s *Server) hasRights(r *http.Request, username, id string) (bool, error) {
	session := jsonapi.GetSession(r)
	user := session.(*ad.User)
	for _, g := range user.Groups {
		if g == s.adminGroup {
			return true, nil
		}
	}

	urls, err := s.db.URLs(username)
	if err != nil {
		return false, fmt.Errorf("Unable to get URLs for user %s: %v", username, err)
	}

	owned := false
	for _, url := range urls {
		if url.ID == id {
			owned = true
		}
	}

	return owned, nil
}

func (s *Server) getHandler(r *http.Request) (int, interface{}) {
	id := mux.Vars(r)["id"]

	session := jsonapi.GetSession(r)
	user := session.Username()

	//check user has rights to url
	ok, err := s.hasRights(r, user, id)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Unable to check if user %s is has rights for URL %s: %v", user, id, err)
	}

	if !ok {
		return http.StatusForbidden, fmt.Errorf("User %s does not have permission to read URL %s", user, id)
	}

	//read url
	url, err := s.db.Get(id)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Unable to get URL %s: %v", id, err)
	}

	if url == nil {
		return http.StatusNotFound, fmt.Errorf("URL %s does not exist", id)
	}

	jsonapi.LogActionID(r, url.ID)

	return http.StatusOK, url
}

func (s *Server) putHandler(r *http.Request) (int, interface{}) {
	type response struct {
		URLID string `json:"url_id"`
	}

	url := new(db.URL)
	d := json.NewDecoder(r.Body)

	if err := d.Decode(url); err != nil {
		return http.StatusBadRequest, fmt.Errorf("Unable to decode request body: %v", err)
	}

	if _, err := neturl.ParseRequestURI(url.URL); err != nil {
		return http.StatusBadRequest, fmt.Errorf(`Unable to parse url "%s": %v`, url.URL, err)
	}

	if url.ID != "" && !regexp.MustCompile(allowedIDRegexp).MatchString(url.ID) {
		return http.StatusBadRequest, fmt.Errorf(`URL ID %s not valid`, url.ID)
	}

	session := jsonapi.GetSession(r)
	user := session.Username()

	id, err := s.db.Put(url, user)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return http.StatusConflict, err
		}

		return http.StatusInternalServerError, fmt.Errorf(`Unable to put URL "%s": %v`, url.URL, err)
	}

	jsonapi.LogActionID(r, url.ID)

	return http.StatusOK, &response{URLID: id}
}

func (s *Server) updateHandler(r *http.Request) (int, interface{}) {
	id := mux.Vars(r)["id"]
	session := jsonapi.GetSession(r)
	user := session.Username()

	//check URL exists
	url, err := s.db.Get(id)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Unable to get URL %s: %v", id, err)
	}

	if url == nil {
		return http.StatusNotFound, fmt.Errorf("URL %s does not exist", id)
	}

	//read url from body
	url = new(db.URL)
	d := json.NewDecoder(r.Body)

	if err = d.Decode(url); err != nil {
		return http.StatusBadRequest, fmt.Errorf("Unable to decode request body: %v", err)
	}

	jsonapi.LogActionID(r, url.ID)

	if _, err = neturl.ParseRequestURI(url.URL); err != nil {
		return http.StatusBadRequest, fmt.Errorf(`Unable to parse url "%s": %v`, url.URL, err)
	}

	//check user has rights to url
	ok, err := s.hasRights(r, user, id)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Unable to check if user %s is has rights for URL %s: %v", user, id, err)
	}

	if !ok {
		return http.StatusForbidden, fmt.Errorf("User %s does not have permission to update URL %s", user, id)
	}

	//update url
	if err = s.db.Update(id, url); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Unable to update URL %s: %v", id, err)
	}

	//re-read url
	url, err = s.db.Get(id)
	if err != nil || url == nil {
		if err == nil {
			err = errors.New("URL unexpectedly nil")
		}
		return http.StatusInternalServerError, fmt.Errorf("Unable to get URL %s: %v", id, err)
	}

	return http.StatusOK, url
}

func (s *Server) deleteHandler(r *http.Request) (int, interface{}) {
	id := mux.Vars(r)["id"]
	session := jsonapi.GetSession(r)
	user := session.Username()

	//check URL exists
	url, err := s.db.Get(id)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Unable to get URL %s: %v", id, err)
	}

	if url == nil {
		return http.StatusNotFound, fmt.Errorf("URL %s does not exist", id)
	}

	//check user has rights to url
	ok, err := s.hasRights(r, user, id)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Unable to check if user %s is has rights for URL %s: %v", user, id, err)
	}

	if !ok {
		return http.StatusForbidden, fmt.Errorf("User %s does not have permission to delete URL %s", user, id)
	}

	jsonapi.LogActionID(r, url.ID)

	//delete url
	if err := s.db.Delete(id); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Unable to delete URL %s: %v", id, err)
	}

	return http.StatusOK, nil
}

func (s *Server) titleHandler(r *http.Request) (int, interface{}) {
	type response struct {
		AppTitle string `json:"app_title"`
	}

	jsonapi.LogActionID(r, s.AppTitle)
	return http.StatusOK, &response{AppTitle: s.AppTitle}
}

func (s *Server) urlsHandler(r *http.Request) (int, interface{}) {
	type response struct {
		URLs []*db.URL `json:"urls"`
	}

	admin := false
	session := jsonapi.GetSession(r)
	username := session.Username()
	user := session.(*ad.User)
	for _, g := range user.Groups {
		if g == s.adminGroup {
			admin = true
			break
		}
	}

	if admin && r.FormValue("all") == "true" {
		username = ""
	}

	urls, err := s.db.URLs(username)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Unable to get URLs for user %s: %v", username, err)
	}

	return http.StatusOK, &response{URLs: urls}
}

func (s *Server) viewHandler(r *http.Request) (int, interface{}) {
	id := mux.Vars(r)["id"]

	url, err := s.db.View(id)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Unable to get URL %s: %v", id, err)
	}

	if url == "" {
		return http.StatusNotFound, fmt.Errorf("URL %s doesn't exist", id)
	}

	return http.StatusOK, url
}
