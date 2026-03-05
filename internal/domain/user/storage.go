package domains

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound      = errors.New("User not found")
	ErrUserAlreadyExists = errors.New("User already exists")
)

type Storage interface {
	GetUser(id uuid.UUID) (User, error)
	GetAllUsers(limit, offset int) ([]User, error)
	Create(user User) error
	Update(user User) error
	Delete(id uuid.UUID) error
}

type storage struct {
	mu    sync.RWMutex
	users map[uuid.UUID]*User
}

func NewStorage() *storage {
	return &storage{users: make(map[uuid.UUID]*User)}
}

func (s *storage) GetUser(id uuid.UUID) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, ok := s.users[id]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return *u, nil
}

func (s *storage) GetAllUsers(limit, offset int) ([]User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]User, 0, len(s.users))
	for _, u := range s.users {
		res = append(res, *u)
	}

	return res, nil
}

func (s *storage) Create(user User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[user.ID]; ok {
		return ErrUserAlreadyExists
	}

	cp := user
	s.users[user.ID] = &cp
	return nil
}

func (s *storage) Update(user User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[user.ID]; !ok {
		return ErrUserNotFound
	}

	cp := user
	s.users[user.ID] = &cp
	return nil
}

func (s *storage) Delete(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[id]; !ok {
		return ErrUserNotFound
	}
	delete(s.users, id)
	return nil
}
