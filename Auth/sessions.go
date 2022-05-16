package auth

import (
	"errors"
	"sync"
)

// Sessions
type Sessions struct {
	sync.Mutex
	store map[string]string
}

func (s *Sessions) Init() *Sessions {
	sessions := new(Sessions)

	// Initialize the map
	sessions.store = make(map[string]string)

	return sessions
}

func (s *Sessions) Get(id string) (*string, error) {
	s.Lock()

	val, ok := s.store[id]

	if !ok {
		return nil, errors.New("token not found in sessions map")
	}

	s.Unlock()

	return &val, nil
}

func (s *Sessions) Set(id string, user string) {
	s.Lock()

	s.store[id] = user

	s.Unlock()
}

func (s *Sessions) Remove(id string) error {
	s.Lock()

	_, ok := s.store[id]

	if !ok {
		return errors.New("token not found in sessions map")
	}

	delete(s.store, id)

	s.Unlock()

	return nil
}
