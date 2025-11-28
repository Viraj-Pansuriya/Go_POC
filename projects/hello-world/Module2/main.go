package main

import (
	"fmt"

	"github.com/viraj/go-mono-repo/projects/hello-world/Module1"
)

func main() {
	userService := Module1.UserService{}
	userService.AddUser(Module1.User{ID: 1, Name: "John"})
	userService.AddUser(Module1.User{ID: 2, Name: "Jane"})
	userService.AddUser(Module1.User{ID: 3, Name: "Jim"})
	userService.AddUser(Module1.User{ID: 4, Name: "Jill"})
	userService.AddUser(Module1.User{ID: 5, Name: "Jack"})
	userService.AddUser(Module1.User{ID: 6, Name: "Jill"})
	userService.AddUser(Module1.User{ID: 7, Name: "Jack"})
	userService.AddUser(Module1.User{ID: 8, Name: "Jill"})
	fmt.Println(userService.GetUsers())

}
