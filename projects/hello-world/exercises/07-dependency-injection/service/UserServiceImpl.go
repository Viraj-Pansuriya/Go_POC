package service

import (
	"sync"

	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/06-http-server/Repository"
	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/06-http-server/models"
	Repository2 "github.com/viraj/go-mono-repo/projects/hello-world/exercises/07-dependency-injection/Repository"
)

var (
	once     sync.Once
	instance UserService
)

type userServiceImpl struct {
	repo Repository.UserRepository
}

func (us *userServiceImpl) FindById(id string) (*models.User, error) {
	return us.repo.FindById(id)
}

func (us *userServiceImpl) AddUser(id string, name string, email string) (*models.User, error) {
	usr := &models.User{ID: id, Name: name, Email: email}
	return us.repo.AddUser(usr)
}

func (us *userServiceImpl) DeleteUser(id string) error {
	return us.repo.DeleteUser(id)
}

func GetUserServiceInstance(repo Repository2.UserRepository) UserService {
	once.Do(func() {
		instance = &userServiceImpl{repo: repo}
	})
	return instance
}
