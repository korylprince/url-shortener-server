package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/korylprince/url-shortener-server/auth/ad"
	"github.com/korylprince/url-shortener-server/db/bbolt"
	"github.com/korylprince/url-shortener-server/httpapi"
	"github.com/korylprince/url-shortener-server/session/memory"
	adauth "gopkg.in/korylprince/go-ad-auth.v1"
)

func main() {
	httpapi.Debug = config.Debug

	db, err := bbolt.New(config.DatabasePath, config.URLIDLength)
	if err != nil {
		log.Fatalln("Unable to create database:", err)
	}

	authConfig := &adauth.Config{
		Server:   config.LDAPServer,
		Port:     config.LDAPPort,
		BaseDN:   config.LDAPBaseDN,
		Security: config.ldapSecurity,
		Debug:    config.Debug,
	}

	auth := ad.New(authConfig, config.LDAPGroup)

	sessionStore := memory.New(time.Minute * time.Duration(config.SessionExpiration))

	s := httpapi.NewServer(db, auth, sessionStore)

	r := httpapi.NewRouter(s, os.Stdout)

	log.Println("Listening on:", config.ListenAddr)
	log.Println(http.ListenAndServe(config.ListenAddr, r))
}
