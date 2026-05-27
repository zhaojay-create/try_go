package user

import (
	"sync"
	"time"
)

// 这一层操作数据库
type User struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	CreateAt time.Time `json:"create_at"`
	Active   bool      `json:"active"`
}

type UserStore struct {
	users  map[int]User
	mutex  sync.RWMutex
	nextId int
}

func NewUserStore() *UserStore {
	return &UserStore{
		users:  make(map[int]User),
		mutex:  sync.RWMutex{},
		nextId: 1,
	}
}

func (s *UserStore) CreateUser(name string, email string) (User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user := User{
		ID:       s.nextId,
		Name:     name,
		Email:    email,
		CreateAt: time.Now(),
		Active:   true,
	}
	s.nextId++
	s.users[user.ID] = user
	return user, nil
}

func (s *UserStore) GetUser(id int) (User, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, exists := s.users[id]
	return user, exists
}

func (s *UserStore) GetAllUser() []User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	users := make([]User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}
