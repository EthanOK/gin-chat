package service

import (
	"gin-chat/models"
	"html/template"

	"strconv"

	"github.com/gin-gonic/gin"
)

// GetIndex
// @Tags 首页
// @Accept json
// @Produce json
// @Success 200 {string} Welcome!!!
// @Router /index [get]
func GetIndex(c *gin.Context) {
	temp, err := template.ParseFiles("index.html", "views/chat/head.html")

	if err != nil {
		panic(err)

	}
	temp.Execute(c.Writer, "index")
}

func GetRegister(c *gin.Context) {
	temp, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		panic(err)

	}
	temp.Execute(c.Writer, "register")
}

func GetChat(c *gin.Context) {

	temp, err := template.ParseFiles("views/chat/index.html",
		"views/chat/head.html",
		"views/chat/tabmenu.html",
		"views/chat/concat.html",
		"views/chat/group.html",
		"views/chat/profile.html",
		"views/chat/main.html",
		"views/chat/createcom.html",
		"views/chat/userinfo.html",
		"views/chat/foot.html")
	if err != nil {
		panic(err)

	}
	user := models.UserBasic{}

	userId, _ := strconv.Atoi(c.Query("userId"))
	user.ID = uint(userId)

	user.Identity = c.Query("token")

	user = models.FindUserByIDAndIdentity(user.ID, user.Identity)
	if user.Name == "" {
		c.Redirect(302, "/")
		return
	}

	temp.Execute(c.Writer, user)
}
