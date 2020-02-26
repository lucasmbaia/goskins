package main

import (
	//"github.com/lucasmbaia/punctorum/api/models/interfaces"
	"github.com/gin-gonic/gin"
	//"github.com/lucasmbaia/punctorum/api/models"
	"github.com/lucasmbaia/goskins/api/controllers"
)

func main() {
	var (
		g	    *gin.Engine
		users	    *controllers.Users
		//inventories *controllers.Inventories
	)

	users = controllers.NewUsers()
	//inventories = controllers.NewInventories()
	g = gin.Default()

	/*v1 := g.Group("/teste/:ID/users")
	{
		v1.GET("", users.Get)
		v1.POST("", users.Post)
	}*/

	v1 := g.Group("/v1")
	{
		v1.GET("/users", users.Get)
		v1.GET("/users/:User", users.Get)
		v1.POST("/users", users.Post)
		//v1.GET("/users/:User/inventories", inventories.Get)
		//v1.GET("/users/:User/inventories/:ID", inventories.Get)
		//v1.POST("/users/:User/inventories", inventories.Get)
	}

	g.Run()
}
