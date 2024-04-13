package router

import (
	"gin-chat/service"

	"gin-chat/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {

	r := gin.Default()

	//	swagger
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 加载静态文件
	r.Static("/asset", "asset/")
	r.LoadHTMLGlob("views/**/*")

	// 首页 登录界面
	r.GET("/", service.GetIndex)
	// 首页 登录界面
	r.GET("/index", service.GetIndex)
	// 注册界面
	r.GET("/register", service.GetRegister)
	// 聊天界面
	r.GET("/chat", service.GetChat)

	// 用户模块 api 接口
	r.GET("/user/getUserList", service.GetUserList)

	r.GET("/getUserByToken", service.GetUserByToken)

	r.POST("/user/createUser", service.CreateUser)

	r.POST("/user/login", service.LoginUser)

	r.POST("/user/deleteUser", service.DeleteUser)

	r.POST("/user/updateUser", service.UpdateUser)

	// 聊天模块 api 接口
	r.POST("/searchFriends", service.SearchFriends)
	// 添加好友
	r.POST("/contact/addFriend", service.AddFriend)
	// 创建群组
	r.POST("/contact/createCommunity", service.CreateCommunity)
	// // 加入群组
	r.POST("/contact/joinGroup", service.JoinGroup)
	// 获取群组列表
	r.POST("/contact/loadcommunity", service.Loadcommunity)

	// webSocket chat
	r.GET("/wsChat", service.Chat)

	// 发送消息
	r.GET("/user/sendMessage", service.SendMessage)

	r.GET("/user/sendUserMessage", service.SendUserMessage)
	// 上传文件
	r.POST("/attach/upload", service.Upload)
	// 上传并更新头像 Avatar
	r.POST("/attach/uploadAvatar", service.UploadAvatar)

	return r

}
