package registration

import (
	"fmt"
	"math/rand"
	"strings"
)

type User struct {
	ID    int
	Name  string
	Email string
	Age   int
}

type UserService struct {
	users map[int]*User // in-memory storage
}

func NewUserService() *UserService {
	return &UserService{}
}

// Register a new user - returns error if validation fails
func (s *UserService) Register(name, email string, age int) (*User, error) {
	user, err := validateUser(name, email, age)

	if err != nil {
		fmt.Println(err.Error())
		return nil, nil
	}
	return user, nil
}

// FindByID - returns ErrUserNotFound if not found
func (s *UserService) FindByID(id int) (*User, error) {
	user, ok := s.users[id]
	if ok {
		return user, nil
	}
	return nil, &ValidationError{
		Field:   "user not found",
		Message: "user nor found",
	}
}

// UpdateEmail - returns validation error if email invalid
func (s *UserService) UpdateEmail(id int, newEmail string) error {

	if strings.Contains(newEmail, "@") == false {
		return ErrInvalidEmail
	}
	s.users[id].Email = newEmail
	return nil
}

// Delete - returns error if user doesn't exist
func (s *UserService) Delete(id int) error {
	_, ok := s.users[id]
	if ok {
		delete(s.users, id)
		return nil
	}
	return ErrUserNotFound
}

func validateUser(name string, email string, age int) (*User, error) {
	if len(name) < 1 {
		return nil, ErrInvalidUser
	}
	if !strings.Contains(email, "@") {
		return nil, ErrInvalidEmail
	}
	if age < 18 {
		return nil, ErrInvalidAge
	}
	return &User{ID: rand.Int(), Name: name, Email: email, Age: age}, nil
}
