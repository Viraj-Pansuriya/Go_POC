package handler

import (
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/06-http-server/models"
	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/09-middleware-auth/auth"
	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/09-middleware-auth/middleware"
	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/09-middleware-auth/repo"
	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/09-middleware-auth/resp"
)

// Instead of generic "unauthorized"
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type UserHandler struct {
	UserRepo *repo.UserRepo
}

func (uh *UserHandler) RegisterHandler(g *gin.Context) {
	user := models.User{}
	err := g.BindJSON(&user)
	at, err := auth.GenerateToken(user.ID, user.Email, user.Role, time.Minute*5)
	rt, err := auth.GenerateToken(user.ID, user.Email, user.Role, time.Hour*24*30*6)
	if err != nil {
		g.AbortWithError(400, err)
	}
	uh.UserRepo.AddUser(&user)
	g.JSON(200, gin.H{"accessToken": at, "refreshToken": rt})
}

func (uh *UserHandler) LoginHandler(g *gin.Context) {
	user := models.User{}
	g.BindJSON(&user)
	if user.ID == 0 || user.Role == "" || user.Email == "" || user.Name == "" {
		g.JSON(400, gin.H{"Error": "User have not enough details passed"})
		return
	}
	at, err := auth.GenerateToken(user.ID, user.Email, user.Role, time.Minute*5)
	rt, err := auth.GenerateToken(user.ID, user.Email, user.Role, time.Hour*24*30*6)
	if err != nil {
		g.AbortWithError(500, err)
		return
	}
	g.JSON(200, gin.H{"accessToken": at, "refreshToken": rt})
}

func (uh *UserHandler) RefreshToken(g *gin.Context) {
	token := middleware.ExtractToken(g)
	claims, err := auth.ValidateToken(token)

	if err != nil {
		str := "Given Refresh Token is not valid , Please login !!"
		g.AbortWithStatusJSON(401, str)
		return
	}
	at, err := auth.GenerateToken(claims.UserID, claims.Email, claims.Role, time.Hour*24*30*6)
	if err != nil {
		g.AbortWithError(500, err)
	}
	g.JSON(200, gin.H{"accessToken": at})
}

func (uh *UserHandler) GetProfileHandler(context *gin.Context) {
	userId, ok := context.GetQuery("userId")
	if ok == false || userId == "" {
		context.AbortWithStatusJSON(400, "UserId can not be empty")
	}

	id64, _ := strconv.ParseUint(userId, 10, 64)
	usr := uh.UserRepo.GetUserById(id64)
	context.JSON(200, gin.H{"user": usr})
}

func (uh *UserHandler) OrderHandler(context *gin.Context) {

	userId, ok := context.GetQuery("userId")
	if ok == false || userId == "" {
		context.AbortWithStatusJSON(400, "UserId can not be empty")
	}
	id := rand.Int64()
	pn := rand.IntN(5)
	prdcts := make([]resp.Product, pn)
	for i := 0; i < pn; i++ {
		pid := rand.Int()
		gofakeit.Seed(0)
		prdcts[i] = resp.Product{
			Id:          uint64(pid),
			ProductName: gofakeit.Name(),
		}
	}
	order := resp.Order{Id: uint64(id), Products: prdcts}
	context.JSON(200, gin.H{"orders": order})

}
