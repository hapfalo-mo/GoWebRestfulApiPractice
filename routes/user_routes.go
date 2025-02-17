package routes

import (
	"my-gin-app/controllers"

	"github.com/gin-gonic/gin"
)

func SetUsertRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.GET("/getAllUser", controllers.GetAllUser)
			users.GET("/getUserByID/:user_id", controllers.GetUserByID)
			users.POST("/createUser", controllers.InsertUser)
			users.PUT("/deleteUser/:user_id", controllers.DeleteUser)
			users.PUT("/updateUser/:user_id", controllers.UpdateUser)
		}
	}
}
