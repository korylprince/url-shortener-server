package memory

import (
	"sync"
	"time"

	"github.com/korylprince/url-shortener-server/rand"
	"github.com/korylprince/url-shortener-server/session"
)

//SessionIDLength is the length of a Session ID
const SessionIDLength = 22

type memorySession struct {
	session *session.Session
	expires time.Time
}

//Store represents a Store that uses an in-memory map
type Store struct {
	store    map[string]*memorySession
	duration time.Duration
	mu       *sync.Mutex
}

//scavenge removes stale records every hour
func scavenge(s *Store) {
	for {
		time.Sleep(time.Hour)
		now := time.Now()
		s.mu.Lock()
		for id, sess := range s.store {
			if sess.expires.Before(now) {
				delete(s.store, id)
			}
		}
		s.mu.Unlock()
	}
}

//New returns a new SessionStore with the given expiration duration.
func New(duration time.Duration) *Store {
	m := &Store{
		store:    make(map[string]*memorySession),
		duration: duration,
		mu:       new(sync.Mutex),
	}
	go scavenge(m)
	return m
}

//Create returns a new session ID for the given session.
//The returned error will always be nil.
func (s *Store) Create(sess *session.Session) (string, error) {
	id := rand.String(SessionIDLength)
	s.mu.Lock()
	s.store[id] = &memorySession{
		session: sess,
		expires: time.Now().Add(s.duration),
	}
	s.mu.Unlock()
	return id, nil
}

//Check returns the session for the given id or nil if it doesn't exist.
//The returned error will always be nil.
func (s *Store) Check(id string) (*session.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sess, ok := s.store[id]; ok {
		if sess.expires.After(time.Now()) {
			sess.expires = time.Now().Add(s.duration)
			return sess.session, nil
		}
		delete(s.store, id)
	}
	return nil, nil
}
