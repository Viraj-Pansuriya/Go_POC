package Module1

type User struct {
	ID   int
	Name string
}

type UserService struct {
	users []User
}

func (s *UserService) GetUsers() []User {
	return s.users
}

func (s *UserService) AddUser(user User) {
	s.users = append(s.users, user)
}

func (s *UserService) removeUser(user User) {
	for i, u := range s.users {
		if u.ID == user.ID {
			s.users = append(s.users[:i], s.users[i+1:]...)
			break
		}
	}
}
