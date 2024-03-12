package router

import (
	"gin-chat/service"

	"gin-chat/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	docs.SwaggerInfo.BasePath = ""

	r := gin.Default()

	r.GET("/index", service.GetIndex)

	r.GET("/user/getUserList", service.GetUserList)

	r.GET("/user/createUser", service.CreateUser)

	r.GET("/user/deleteUser", service.DeleteUser)

	r.POST("/user/updateUser", service.UpdateUser)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r

}
