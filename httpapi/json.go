package httpapi

import (
	"encoding/json"
	"errors"
	"log"
	"mime"
	"net/http"
)

func jsonRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			jsonResponse(http.StatusBadRequest, errors.New("Could not parse Content-Type")).ServeHTTP(w, r)
			return
		}

		if mediaType != "application/json" {
			jsonResponse(http.StatusBadRequest, errors.New("Content-Type not application/json")).ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func jsonResponse(code int, data interface{}) http.Handler {
	type response struct {
		Code        int    `json:"code"`
		Description string `json:"description"`
		Debug       string `json:"debug,omitempty"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err, ok := data.(error); ok || data == nil {
			resp := response{Code: code, Description: http.StatusText(code)}
			if err != nil {
				(r.Context().Value(contextKeyLogData)).(*logData).Error = err.Error()
				if Debug {
					resp.Debug = err.Error()
				}
			}
			data = resp
		}

		if code == http.StatusUnauthorized {
			w.Header().Set("WWW-Authenticate", `Session realm="api"`)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)

		e := json.NewEncoder(w)
		eErr := e.Encode(data)

		if eErr != nil {
			log.Println("Error writing JSON response:", eErr)
		}
	})
}
