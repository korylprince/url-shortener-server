package db

//DB represents a URL shortening database
type DB interface {
	//Get returns the *URL with the given id, or an error if one occurred.
	//If a *URL with the given id doesn't exist, url will be nil
	Get(id string) (url *URL, err error)

	//Put saves the given url in the database for the given user, returning the id, or an error
	//if one occurred.
	Put(url *URL, user string) (id string, err error)

	//Update updates the *URL with the given id or returns an error if one occurred
	Update(id string, url *URL) error

	//Delete deletes the *URL with the given id or returns an error if one occurred
	Delete(id string) error

	//View returns the url with the given id, or an error if one occurred.
	//If a url with the given id doesn't exist, url will be empty.
	//View increments the view counter for the URL and should be used
	//by clients wanting to resolve the shortened URL.
	View(id string) (url string, err error)

	//URLs returns the URLs for the given user or all URLs if user is empty
	//or an error if one occurred
	URLs(user string) ([]*URL, error)
}
