package repo

import (
	"sync"

	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/06-http-server/models"
)

var (
	userRepo UserRepo
	doOnce   sync.Once
)

type UserRepo struct {
	users map[uint64]*models.User
}

func GetUserRepoInstance() *UserRepo {
	doOnce.Do(func() {
		userRepo = UserRepo{
			users: make(map[uint64]*models.User),
		}
	})
	return &userRepo
}

func (ur *UserRepo) AddUser(user *models.User) {
	ur.users[user.ID] = user
}

func (ur *UserRepo) GetUserById(id uint64) *models.User {
	return ur.users[id]
}

func (ur *UserRepo) DeleteById(id uint64) {
	delete(ur.users, id)
}
