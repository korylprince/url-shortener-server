package ad

import (
	"gopkg.in/korylprince/go-ad-auth.v1"
)

//Auth represents an Active Directory authentication mechanism
type Auth struct {
	config *auth.Config
	group  string
}

//New returns a new *Auth with the given configuration. If group is not empty
//membership in the given group is necessary to authenticate
func New(config *auth.Config, group string) *Auth {
	return &Auth{config: config, group: group}
}

//Authenticate returns whether or not the given username and password
//successfully authenticates or an error if one occurred
func (a *Auth) Authenticate(username, password string) (bool, error) {
	return auth.Login(username, password, a.group, a.config)
}
