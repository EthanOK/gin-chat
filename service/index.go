package service

import (
	"html/template"

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
