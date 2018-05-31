package main

import (
	"log"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/korylprince/go-ad-auth.v1"
)

//Config represents options given in the environment
type Config struct {
	SessionExpiration int `default:"60"` //in minutes

	DatabasePath string `required:"true"`

	URLIDLength int `default:"6" required:"true"`

	AppTitle string

	LDAPServer     string `required:"true"`
	LDAPPort       int    `default:"389" required:"true"`
	LDAPBaseDN     string `required:"true"`
	LDAPGroup      string
	LDAPAdminGroup string
	LDAPSecurity   string `default:"none" required:"true"`
	ldapSecurity   auth.SecurityType

	TLSCert string
	TLSKey  string

	ListenAddr string `default:":8080" required:"true"` //addr format used for net.Dial; required
	Prefix     string //url prefix to mount api to without trailing slash
	Debug      bool   `default:"false"` //return debugging information to client
}

var config = &Config{}

func init() {
	err := envconfig.Process("SHORTENER", config)
	if err != nil {
		log.Fatalln("Error reading configuration from environment:", err)
	}

	switch strings.ToLower(config.LDAPSecurity) {
	case "", "none":
		config.ldapSecurity = auth.SecurityNone
	case "tls":
		config.ldapSecurity = auth.SecurityTLS
	case "starttls":
		config.ldapSecurity = auth.SecurityStartTLS
	default:
		log.Fatalln("Invalid SHORTENER_LDAPSECURITY:", config.LDAPSecurity)
	}
}
