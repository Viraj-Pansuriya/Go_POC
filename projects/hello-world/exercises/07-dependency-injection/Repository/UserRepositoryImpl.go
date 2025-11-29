package Repository

import (
	"errors"
	"sync"

	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/06-http-server/models"
)

type userRepositoryImpl struct {
	users map[string]*models.User
	mu    sync.RWMutex
}

func (ur *userRepositoryImpl) FindById(id string) (*models.User, error) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()
	if ur.users == nil {
		return nil, errors.New("userRepositoryImpl FindById() - users is nil")
	}
	usr, ok := ur.users[id]
	if ok == false {
		return nil, errors.New("userRepositoryImpl FindById() - user not found")
	}
	return usr, nil
}

func (ur *userRepositoryImpl) AddUser(user *models.User) (*models.User, error) {
	ur.mu.Lock()
	defer ur.mu.Unlock()
	_, ok := ur.users[user.ID]
	if ok {
		return nil, errors.New("userRepositoryImpl AddUser() - user already exists")
	}
	ur.users[user.ID] = user
	return user, nil
}

func (ur *userRepositoryImpl) DeleteUser(id string) error {
	ur.mu.Lock()
	defer ur.mu.Unlock()
	_, ok := ur.users[id]
	if ok == false {
		return errors.New("userRepositoryImpl FindById() - user not found")
	}
	delete(ur.users, id)
	return nil
}

func GetUserRepositoryInstance() UserRepository {
	once.Do(func() {
		initUsers := []models.User{
			{"1", "Viraj", "viraj777@gmail.com"},
		}
		users := make(map[string]*models.User, len(initUsers))
		for _, usr := range initUsers {
			users[usr.ID] = &usr
		}
		instance = &userRepositoryImpl{users: users}
	})
	return instance
}

var (
	instance UserRepository
	once     sync.Once
)
