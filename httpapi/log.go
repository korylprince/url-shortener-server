package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

//Debug will pass debugging information to the client if true
var Debug = false

type logData struct {
	Action string      `json:"action"`
	URLID  string      `json:"url_id,omitempty"`
	User   string      `json:"user,omitempty"`
	Admin  bool        `json:"admin,omitempty"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
	Time   time.Time   `json:"time"`
}

func setAction(action string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		(r.Context().Value(contextKeyLogData)).(*logData).Action = action
		if id, ok := mux.Vars(r)["id"]; ok {
			(r.Context().Value(contextKeyLogData)).(*logData).URLID = id
		}
		next.ServeHTTP(w, r)
	})
}

func logRequest(w io.Writer, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := new(logData)

		ctx := context.WithValue(r.Context(), contextKeyLogData, l)
		next.ServeHTTP(w, r.WithContext(ctx))

		l.Time = time.Now()
		j, err := json.Marshal(l)
		if err != nil {
			log.Println("Unable to marshal JSON:", err)
		}
		fmt.Println(string(j))
	})
}
