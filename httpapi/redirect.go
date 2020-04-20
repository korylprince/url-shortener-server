package httpapi

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/korylprince/httputil/jsonapi"
)

func withRedirect(next jsonapi.ReturnHandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code, body := next(r)

		if code == http.StatusOK {
			http.Redirect(w, r, body.(string), http.StatusTemporaryRedirect)
			return
		}

		if err, ok := body.(error); ok {
			log.Printf("Redirecting error (%d %s): %v\n", code, http.StatusText(code), err)
		}

		u := &url.URL{Path: "error.html"}
		v := make(url.Values)
		v.Set("statusCode", strconv.Itoa(code))
		v.Set("statusText", http.StatusText(code))
		u.RawQuery = v.Encode()

		http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
	})
}
