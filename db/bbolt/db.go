package bbolt

import (
	"encoding/binary"
	"fmt"
	"log"
	"time"

	bolt "github.com/coreos/bbolt"
	"github.com/korylprince/url-shortener-server/db"
	"github.com/korylprince/url-shortener-server/rand"
)

var urlsBucket = []byte("urls")
var usersBucket = []byte("users")

var userKey = []byte("user")
var urlKey = []byte("url")
var viewsKey = []byte("views")
var expiresKey = []byte("expires")
var deletedKey = []byte("deleted")
var modifiedKey = []byte("modified")

func getURL(b *bolt.Bucket) (*db.URL, error) {
	bUser := b.Get(userKey)
	if bUser == nil {
		return nil, fmt.Errorf(`Unable to get "%s" value: value is nil`, userKey)
	}

	bURL := b.Get(urlKey)
	if bURL == nil {
		return nil, fmt.Errorf(`Unable to get "%s" value: value is nil`, urlKey)
	}

	bViews := b.Get(viewsKey)
	if bViews == nil {
		return nil, fmt.Errorf(`Unable to get "%s" value: value is nil`, viewsKey)
	}

	views, read := binary.Uvarint(bViews)
	if read < 1 {
		return nil, fmt.Errorf(`Unable to decode "%s" value "%v": number of bytes read is %d`, viewsKey, bViews, read)
	}

	url := &db.URL{
		User:  string(bUser),
		URL:   string(bURL),
		Views: views,
	}

	if bExpires := b.Get(expiresKey); bExpires != nil {
		expires := new(time.Time)
		if err := expires.UnmarshalBinary(bExpires); err != nil {
			return nil, fmt.Errorf(`Unable to get "%s" value: %v`, expiresKey, err)
		}
		url.Expires = expires
	}

	bModified := b.Get(modifiedKey)
	modified := new(time.Time)
	if err := modified.UnmarshalBinary(bModified); err != nil {
		return nil, fmt.Errorf(`Unable to get "%s" value: %v`, modifiedKey, err)
	}
	url.LastModified = modified

	return url, nil
}

func putURL(b *bolt.Bucket, url *db.URL) error {
	if err := b.Put(userKey, []byte(url.User)); err != nil {
		return fmt.Errorf(`Unable to put "%s" value "%s": %v`, userKey, url.User, err)
	}

	if err := b.Put(urlKey, []byte(url.URL)); err != nil {
		return fmt.Errorf(`Unable to put "%s" value "%s": %v`, urlKey, url.URL, err)
	}

	bViews := make([]byte, 8) //size of uint64
	binary.PutUvarint(bViews, url.Views)

	if err := b.Put(viewsKey, bViews); err != nil {
		return fmt.Errorf(`Unable to put "%s" value "%v": %v`, viewsKey, bViews, err)
	}

	if url.Expires != nil {
		bExpires, err := url.Expires.MarshalBinary()
		if err != nil {
			return fmt.Errorf(`Unable to marshal "%s" value "%v": %v`, expiresKey, url.Expires, err)
		}
		if err = b.Put(expiresKey, bExpires); err != nil {
			return fmt.Errorf(`Unable to put "%s" value "%s": %v`, expiresKey, bExpires, err)
		}
	} else if err := b.Delete(expiresKey); err != nil {
		return fmt.Errorf(`Unable to delete "%s": %v`, expiresKey, err)
	}

	modified := time.Now()
	bModified, err := modified.MarshalBinary()
	if err != nil {
		return fmt.Errorf(`Unable to marshal "%s" value "%v": %v`, expiresKey, modified, err)
	}
	if err = b.Put(modifiedKey, bModified); err != nil {
		return fmt.Errorf(`Unable to put "%s" value "%s": %v`, modifiedKey, bModified, err)
	}

	return nil
}

//DB is a BBolt DB
type DB struct {
	db       *bolt.DB
	idLength int
}

//New returns a new *BBoltDB with the given path and id length
//or an error if one occurred
func New(path string, idLength int) (*DB, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to create database %s: %v", path, err)
	}

	tx, err := db.Begin(true)
	if err != nil {
		return nil, fmt.Errorf("Unable to open database %s for writing: %v", path, err)
	}

	_, err = tx.CreateBucketIfNotExists(urlsBucket)
	if err != nil {
		return nil, fmt.Errorf(`Unable to create database %s "%s" bucket: %v`, path, urlsBucket, err)
	}

	_, err = tx.CreateBucketIfNotExists(usersBucket)
	if err != nil {
		return nil, fmt.Errorf(`Unable to create database %s "%s" bucket: %v`, path, usersBucket, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("Unable to commit transaction %s: %v", path, err)
	}

	return &DB{db: db, idLength: idLength}, nil
}

//Get returns the *URL with the given id, or an error if one occurred.
//If a *URL with the given id doesn't exist, url will be nil
func (d *DB) Get(id string) (url *db.URL, err error) {
	tx, err := d.db.Begin(false)
	if err != nil {
		return nil, fmt.Errorf("Unable to open database for reading: %v", err)
	}

	defer func() {
		if rErr := tx.Rollback(); rErr != nil && err == nil {
			err = fmt.Errorf("Unable to rollback read-only transaction: %v", rErr)
		}
	}()

	ub := tx.Bucket(urlsBucket)
	if ub == nil {
		return nil, fmt.Errorf(`Unable to open database "%s" bucket: bucket is nil`, urlsBucket)
	}

	b := ub.Bucket([]byte(id))
	if b == nil {
		return nil, nil
	}

	url, err = getURL(b)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal URL: %v", err)
	}

	//check for deleted URL
	if b.Get(deletedKey) != nil {
		return nil, nil
	}

	url.ID = id

	return url, err
}

//Put saves the given url in the database for the given user, returning the id, or an error if one occurred.
func (d *DB) Put(url *db.URL, user string) (id string, err error) {
	tx, err := d.db.Begin(true)
	if err != nil {
		return "", fmt.Errorf("Unable to open database for writing: %v", err)
	}

	defer func() {
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				log.Println("WARNING: Unable to rollback failed transaction:", rErr)
			}
			return
		}

		if cErr := tx.Commit(); cErr != nil {
			err = fmt.Errorf("Unable to commit transaction: %v", cErr)
		}
	}()

	//store URL
	ub := tx.Bucket(urlsBucket)
	if ub == nil {
		return "", fmt.Errorf(`Unable to open database "%s" bucket: bucket is nil`, urlsBucket)
	}

	id = url.ID

	if id == "" {
		//generate random, unused id if not set
		for {
			id = rand.String(d.idLength)
			if ub.Bucket([]byte(id)) == nil {
				break
			}
		}
	} else if ub.Bucket([]byte(id)) != nil {
		return "", fmt.Errorf("URL %s already exists", id)
	}

	b, err := ub.CreateBucketIfNotExists([]byte(id))
	if err != nil {
		return "", fmt.Errorf(`Unable to create url "%s" bucket: %v`, id, err)
	}

	url.User = user
	url.Views = 0

	if err = putURL(b, url); err != nil {
		return "", fmt.Errorf(`Unable to marshal "%s" URL: %v`, id, err)
	}

	//store user
	ub = tx.Bucket(usersBucket)
	if ub == nil {
		return "", fmt.Errorf(`Unable to open database "%s" bucket: bucket is nil`, usersBucket)
	}

	b, err = ub.CreateBucketIfNotExists([]byte(user))
	if err != nil {
		return "", fmt.Errorf(`Unable to create user "%s" bucket: %v`, user, err)
	}

	if err = b.Put([]byte(id), nil); err != nil {
		return "", fmt.Errorf(`Unable to add url "%s" to user "%s": %v`, id, user, err)
	}

	return id, nil
}

//Update updates the *URL with the given id or returns an error if one occurred
func (d *DB) Update(id string, url *db.URL) error {
	u, err := d.Get(id)
	if err != nil {
		return fmt.Errorf(`Unable to get URL "%s": %v`, id, err)
	}
	if u == nil {
		return fmt.Errorf(`Unable to get URL "%s": URL doesn't exist`, id)
	}

	url.ID = id
	url.User = u.User
	url.Views = u.Views

	tx, err := d.db.Begin(true)
	if err != nil {
		return fmt.Errorf("Unable to open database for writing: %v", err)
	}

	defer func() {
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				log.Println("WARNING: Unable to rollback failed transaction:", rErr)
			}
			return
		}

		if cErr := tx.Commit(); cErr != nil {
			err = fmt.Errorf("Unable to commit transaction: %v", cErr)
		}
	}()

	ub := tx.Bucket(urlsBucket)
	if ub == nil {
		return fmt.Errorf(`Unable to open database "%s" bucket: bucket is nil`, urlsBucket)
	}

	b := ub.Bucket([]byte(id))
	if b == nil {
		return fmt.Errorf(`Unable to open url "%s" bucket: bucket is nil`, id)
	}

	if err = putURL(b, url); err != nil {
		return fmt.Errorf(`Unable to update "%s" URL: %v`, id, err)
	}

	return nil
}

//Delete deletes the *URL with the given id or returns an error if one occurred
func (d *DB) Delete(id string) error {
	tx, err := d.db.Begin(true)
	if err != nil {
		return fmt.Errorf("Unable to open database for writing: %v", err)
	}

	defer func() {
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				log.Println("WARNING: Unable to rollback failed transaction:", rErr)
			}
			return
		}

		if cErr := tx.Commit(); cErr != nil {
			err = fmt.Errorf("Unable to commit transaction: %v", cErr)
		}
	}()

	//
	ub := tx.Bucket(urlsBucket)
	if ub == nil {
		return fmt.Errorf(`Unable to open database "%s" bucket: bucket is nil`, urlsBucket)
	}

	b := ub.Bucket([]byte(id))
	if b == nil {
		return fmt.Errorf(`Unable to open url "%s" bucket: bucket is nil`, id)
	}

	modified := time.Now()
	bModified, err := modified.MarshalBinary()
	if err != nil {
		return fmt.Errorf(`Unable to marshal "%s" value "%v": %v`, expiresKey, modified, err)
	}
	if err = b.Put(modifiedKey, bModified); err != nil {
		return fmt.Errorf(`Unable to put "%s" value "%s": %v`, modifiedKey, bModified, err)
	}

	return b.Put(deletedKey, nil)
}

//View returns the url with the given id, or an error if one occurred.
//If a URL with the given id doesn't exist, url will be empty.
//View increments the view counter for the URL and should be used
//by clients wanting to resolve the shortened URL.
func (d *DB) View(id string) (url string, err error) {
	u, err := d.Get(id)
	if err != nil {
		return "", fmt.Errorf(`Unable to get url "%s": %v`, id, err)
	}

	//check exists
	if u == nil {
		return "", nil
	}

	//check expired
	if u.Expires != nil && time.Now().After(*(u.Expires)) {
		return "", nil
	}

	//increment views
	tx, err := d.db.Begin(true)
	if err != nil {
		return "", fmt.Errorf("Unable to open database for writing: %v", err)
	}

	defer func() {
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				log.Println("WARNING: Unable to rollback failed transaction:", rErr)
			}
			return
		}

		if cErr := tx.Commit(); cErr != nil {
			err = fmt.Errorf("Unable to commit transaction: %v", cErr)
		}
	}()

	ub := tx.Bucket(urlsBucket)
	if ub == nil {
		return "", fmt.Errorf(`Unable to open database "%s" bucket: bucket is nil`, urlsBucket)
	}

	b := ub.Bucket([]byte(id))
	if b == nil {
		return "", fmt.Errorf(`Unable to open url "%s" bucket: bucket is nil`, id)
	}

	bViews := make([]byte, 8) //size of uint64
	binary.PutUvarint(bViews, u.Views+1)

	if err = b.Put(viewsKey, bViews); err != nil {
		return "", fmt.Errorf(`Unable to put "%s" value "%v": %v`, viewsKey, bViews, err)
	}

	return u.URL, nil
}

func (d *DB) getUserIDs(user string) ([]string, error) {
	tx, err := d.db.Begin(false)
	if err != nil {
		return nil, fmt.Errorf("Unable to open database for reading: %v", err)
	}

	defer func() {
		if rErr := tx.Rollback(); rErr != nil && err == nil {
			err = fmt.Errorf("Unable to rollback read-only transaction: %v", rErr)
		}
	}()

	ub := tx.Bucket(usersBucket)
	if ub == nil {
		return nil, fmt.Errorf(`Unable to open database "%s" bucket: bucket is nil`, usersBucket)
	}

	b := ub.Bucket([]byte(user))
	if b == nil {
		return nil, nil
	}

	var ids []string

	err = b.ForEach(func(k, v []byte) error {
		ids = append(ids, string(k))
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf(`Unable to read user "%s" URLs: %v`, user, err)
	}

	return ids, nil
}

func (d *DB) getAllIDs() ([]string, error) {
	tx, err := d.db.Begin(false)
	if err != nil {
		return nil, fmt.Errorf("Unable to open database for reading: %v", err)
	}

	defer func() {
		if rErr := tx.Rollback(); rErr != nil && err == nil {
			err = fmt.Errorf("Unable to rollback read-only transaction: %v", rErr)
		}
	}()

	ub := tx.Bucket(urlsBucket)
	if ub == nil {
		return nil, fmt.Errorf(`Unable to open database "%s" bucket: bucket is nil`, urlsBucket)
	}

	var ids []string

	err = ub.ForEach(func(k, v []byte) error {
		if v == nil {
			ids = append(ids, string(k))
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Unable to read URL IDs: %v", err)
	}

	return ids, nil
}

//URLs returns the URLs for the given user or all URLs if user is empty
//or an error if one occurred
func (d *DB) URLs(user string) ([]*db.URL, error) {
	urls := make([]*db.URL, 0)

	var ids []string
	var err error
	if user == "" {
		ids, err = d.getAllIDs()
		if err != nil {
			return nil, fmt.Errorf("Unable to get URL IDs: %v", err)
		}
	} else {
		ids, err = d.getUserIDs(user)
		if err != nil {
			return nil, fmt.Errorf(`Unable to get user "%s" URL IDs: %v`, user, err)
		}
	}

	for _, id := range ids {
		url, err := d.Get(id)
		if err != nil {
			return nil, fmt.Errorf(`Unable to get URL "%s": %v`, id, err)
		}
		if url != nil {
			url.ID = id
			urls = append(urls, url)
		}
	}

	return urls, nil
}
