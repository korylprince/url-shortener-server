package main

import (
	"embed"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	_ "embed"

	"github.com/kelseyhightower/envconfig"
	auth "github.com/korylprince/go-ad-auth/v3"
	"github.com/korylprince/httputil/auth/ad"
	"github.com/korylprince/httputil/session/memory"
	"github.com/korylprince/url-shortener-server/v2/db/bbolt"
	"github.com/korylprince/url-shortener-server/v2/httpapi"
)

//go:embed client
var httpEmbed embed.FS

func main() {
	config := new(Config)
	err := envconfig.Process("SHORTENER", config)
	if err != nil {
		log.Fatalln("Error reading configuration from environment:", err)
	}

	rand.Seed(time.Now().Unix())

	db, err := bbolt.New(config.DatabasePath, config.URLIDLength)
	if err != nil {
		log.Fatalln("Unable to create database:", err)
	}

	authConfig := &auth.Config{
		Server:   config.LDAPServer,
		Port:     config.LDAPPort,
		BaseDN:   config.LDAPBaseDN,
		Security: config.SecurityType(),
	}

	auth := ad.New(authConfig, nil, []string{config.LDAPGroup, config.LDAPAdminGroup})

	sessionStore := memory.New(time.Minute * time.Duration(config.SessionExpiration))

	s := httpapi.NewServer(config.AppTitle, db, auth, config.LDAPAdminGroup, sessionStore, httpEmbed, os.Stdout)

	log.Println("Listening on:", config.ListenAddr)

	if config.TLSCert != "" && config.TLSKey != "" {
		log.Println(http.ListenAndServeTLS(config.ListenAddr, config.TLSCert, config.TLSKey, s.Router()))
	} else {
		log.Println(http.ListenAndServe(config.ListenAddr, s.Router()))
	}
}
