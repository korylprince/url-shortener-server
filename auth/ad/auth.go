package ad

import (
	"fmt"
	"strings"

	auth "gopkg.in/korylprince/go-ad-auth.v1"
)

//Auth represents an Active Directory authentication mechanism
type Auth struct {
	config     *auth.Config
	group      string
	adminGroup string
}

//New returns a new *Auth with the given configuration. If group is not empty
//membership in the given group is necessary to authenticate
func New(config *auth.Config, group string, adminGroup string) *Auth {
	return &Auth{config: config, group: group, adminGroup: adminGroup}
}

//Authenticate returns whether or not the given username and password
//successfully authenticates and if they are an admin, or an error if one occurred
func (a *Auth) Authenticate(username, password string) (status bool, admin bool, err error) {
	if a.adminGroup == "" {
		status, err = auth.Login(username, password, a.group, a.config)
		return status, false, err
	}

	status, attrs, err := auth.LoginWithAttrs(username, password, a.group, a.config, []string{"memberOf"})
	if err != nil {
		return false, false, fmt.Errorf("Error attempting to authenticate with attributes: %v", err)
	}

	if !status {
		return false, false, nil
	}

	memberOf := attrs["memberOf"]

	if memberOf == nil {
		return true, false, nil
	}

	for _, dn := range memberOf {
		if strings.HasPrefix(dn, fmt.Sprintf("CN=%s,", a.adminGroup)) {
			admin = true
			break
		}
	}

	return true, admin, nil
}
