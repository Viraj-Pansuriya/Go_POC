package service

import "github.com/viraj/go-mono-repo/projects/hello-world/exercises/06-http-server/models"

type UserService interface {
	FindById(id string) (*models.User, error)
	AddUser(id string, name string, email string) (*models.User, error)
	DeleteUser(id string) error
}
