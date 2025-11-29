package main

import (
	"github.com/gin-gonic/gin"
	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/07-dependency-injection/Repository"
	handler2 "github.com/viraj/go-mono-repo/projects/hello-world/exercises/07-dependency-injection/handler"
	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/07-dependency-injection/service"
)

func main() {

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"Status": "UP"})
	})

	repo := Repository.GetUserRepositoryInstance()
	userService := service.GetUserServiceInstance(repo)
	handler := handler2.NewUserHandler(userService)
	handler.RegisterRoutes(r)

	err := r.Run(":8080")
	if err != nil {
		return
	}

	//ch := make(chan struct{}, 3)
	//wg := sync.WaitGroup{}
	//for index := 0; index < 10; index++ {
	//	wg.Add(1)
	//	go func(index int) {
	//		ch <- struct{}{}
	//		fmt.Println("printing a element", index)
	//		time.Sleep(1 * time.Second)
	//		defer func() { <-ch }()
	//		defer wg.Done()
	//	}(index)
	//}
	//wg.Wait()
}
