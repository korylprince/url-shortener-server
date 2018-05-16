package auth

//Auth represents an authentication mechanism
type Auth interface {
	//Authenticate returns whether or not the given username and password
	//successfully authenticates or an error if one occurred
	Authenticate(username, password string) (bool, error)
}
