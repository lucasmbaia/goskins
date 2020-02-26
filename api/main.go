package main

import (
	//"github.com/lucasmbaia/punctorum/api/models/interfaces"
	"github.com/gin-gonic/gin"
	//"github.com/lucasmbaia/punctorum/api/models"
	"github.com/lucasmbaia/goskins/api/repository/gorm"
	"github.com/lucasmbaia/goskins/api/controllers"
	"github.com/lucasmbaia/goskins/api/config"
)

func init() {
	config.EnvConfig = config.Config{
		DBFields:   gorm.GormConfig{
			Username:	    "goskins",
			Password:	    "123456",
			Host:		    "127.0.0.1",
			Port:		    "3306",
			DBName:		    "goskins",
			Timeout:	    "30000ms",
			Debug:		    true,
			ConnsMaxIdle:	    5,
			ConnsMaxOpen:	    5,
			ConnsMaxLifetime:   5,
		},
	}

	config.LoadDB()
}

func main() {
	var (
		g	    *gin.Engine
		users	    *controllers.Users
		inventories *controllers.Inventories
	)

	users = controllers.NewUsers()
	inventories = controllers.NewInventories()
	g = gin.Default()

	/*v1 := g.Group("/teste/:ID/users")
	{
		v1.GET("", users.Get)
		v1.POST("", users.Post)
	}*/

	v1 := g.Group("/v1")
	{
		v1.GET("/users", users.Get)
		v1.GET("/users/:user", users.Get)
		v1.POST("/users", users.Post)
		v1.GET("/users/:user/inventories", inventories.Get)
		v1.GET("/users/:user/inventories/:inventory", inventories.Get)
		v1.POST("/users/:user/inventories", inventories.Post)
	}

	g.Run()
}
