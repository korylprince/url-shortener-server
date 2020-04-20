package main

import (
	"log"
	"strings"

	auth "github.com/korylprince/go-ad-auth/v3"
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

	TLSCert string
	TLSKey  string

	ListenAddr string `default:":8080" required:"true"` //addr format used for net.Dial; required
	Prefix     string //url prefix to mount api to without trailing slash
}

//SecurityType returns the auth.SecurityType for the config
func (c *Config) SecurityType() auth.SecurityType {
	switch strings.ToLower(c.LDAPSecurity) {
	case "", "none":
		return auth.SecurityNone
	case "tls":
		return auth.SecurityTLS
	case "starttls":
		return auth.SecurityStartTLS
	default:
		log.Fatalln("Invalid SHORTENER_LDAPSECURITY:", c.LDAPSecurity)
	}
	panic("unreachable")
}
