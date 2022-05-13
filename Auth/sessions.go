package auth

import "sync"

// Sessions
type Sessions struct {
	sync.Mutex
	store map[string]string
}

func (s *Sessions) New() *Sessions {
	sessions := new(Sessions)

	// Initialize the map
	sessions.store = make(map[string]string)

	return sessions
}

func (s *Sessions) Get(id string) (*string, error) {
	return nil, nil
}

func (s *Sessions) Set(id string, user User) error {
	return nil
}

func (s *Sessions) Remove(id string) error {
	return nil
}
