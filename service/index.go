package service

import (
	"gin-chat/models"
	"gin-chat/utils"
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

	token := c.Query("token")

	claim, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(401, gin.H{

			"code": 401,
			"msg":  "登录已过期,请重新登录",
		})
		return
	}

	user = models.FindUserById(user.ID)

	if user.Name != claim.Name {
		c.Redirect(302, "/")
		return
	}

	temp.Execute(c.Writer, user)
}

func Chat(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
