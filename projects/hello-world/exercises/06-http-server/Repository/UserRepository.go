package Repository

import (
	"errors"
	"sync"

	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/06-http-server/models"
)

type UserRepository interface {
	FindById(id string) (*models.User, error)
	AddUser(user *models.User) (*models.User, error)
	DeleteUser(id string) error
}

type UserRepositoryImpl struct {
	users map[string]*models.User
}

func (ur *UserRepositoryImpl) FindById(id string) (*models.User, error) {
	if ur.users == nil {
		return nil, errors.New("UserRepositoryImpl FindById() - users is nil")
	}
	usr, ok := ur.users[id]
	if ok == false {
		return nil, errors.New("UserRepositoryImpl FindById() - user not found")
	}
	return usr, nil
}

func (ur *UserRepositoryImpl) AddUser(user *models.User) (*models.User, error) {
	_, ok := ur.users[user.ID]
	if ok {
		return nil, errors.New("UserRepositoryImpl AddUser() - user already exists")
	}
	ur.users[user.ID] = user
	return user, nil
}

func (ur *UserRepositoryImpl) DeleteUser(id string) error {
	_, ok := ur.users[id]
	if ok == false {
		return errors.New("UserRepositoryImpl FindById() - user not found")
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
		instance = &UserRepositoryImpl{users: users}
	})
	return instance
}

var (
	instance UserRepository
	once     sync.Once
)
