package main

import (
	"github.com/gin-gonic/gin"
	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/09-middleware-auth/handler"
	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/09-middleware-auth/middleware"
	repo2 "github.com/viraj/go-mono-repo/projects/hello-world/exercises/09-middleware-auth/repo"
)

func main() {
	r := gin.Default()

	repo := repo2.GetUserRepoInstance()
	userHandler := handler.UserHandler{repo}

	r.Group("api/v1")
	public := r.Group("api/v1")
	{
		public.POST("/register", userHandler.RegisterHandler)
		public.GET("/login", userHandler.LoginHandler)
		public.GET("/refresh", userHandler.RefreshToken) // in case of refresh pass RT and based on that create a new AT;
	}
	private := r.Group("api/v1")
	{
		private.GET("/profile", middleware.AuthMiddleware(), userHandler.GetProfileHandler)
		private.GET("/order", middleware.AuthMiddleware(), userHandler.OrderHandler)
	}

	admin := r.Group("api/v1")
	{
		admin.Use(middleware.AuthMiddleware())
		admin.Use(middleware.RoleRequired("admin"))
		// make sure it is after registration of middlewares;
		admin.GET("/admin")
	}
	err := r.Run(":8082")
	if err != nil {
		return
	}
}
