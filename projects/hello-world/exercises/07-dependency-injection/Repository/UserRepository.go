package Repository

import "github.com/viraj/go-mono-repo/projects/hello-world/exercises/06-http-server/models"

type UserRepository interface {
	FindById(id string) (*models.User, error)
	AddUser(user *models.User) (*models.User, error)
	DeleteUser(id string) error
}
