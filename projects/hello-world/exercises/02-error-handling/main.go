package main

import (
	"errors"
	"fmt"

	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/02-error-handling/registration"
)

func main() {

	userService := &registration.UserService{}
	_, err := userService.Register("viraj", "virajpansuriya777@gmail.com", 1)

	if err != nil {
		notOk := errors.As(err, &registration.ErrInvalidUser)
		fmt.Println(notOk)
	}

}
