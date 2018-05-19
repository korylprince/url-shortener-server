package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gobuffalo/packr"
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

	box := packr.NewBox("./client/dist")

	s := httpapi.NewServer(db, auth, sessionStore, box)

	r := httpapi.NewRouter(s, os.Stdout)

	log.Println("Listening on:", config.ListenAddr)

	if config.TLSCert != "" && config.TLSKey != "" {
		log.Println(http.ListenAndServeTLS(config.ListenAddr, config.TLSCert, config.TLSKey, r))
	} else {
		log.Println(http.ListenAndServe(config.ListenAddr, r))
	}
}
