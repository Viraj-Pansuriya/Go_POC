package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/06-http-server/Repository"
	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/06-http-server/models"
)

var (
	repo = Repository.GetUserRepositoryInstance()
)

func AddUserHandler(w http.ResponseWriter, r *http.Request) {

}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func GetUser(c *gin.Context) {
	id := c.Param("id")
	resp, _ := repo.FindById(id)
	c.JSON(200, resp)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	_ = repo.DeleteUser(id)
	c.JSON(200, gin.H{"id": id})
}
func AddUser(c *gin.Context) {
	var user models.User
	c.BindJSON(&user)
	usr, _ := repo.AddUser(user)
	c.JSON(200, usr)
}

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
