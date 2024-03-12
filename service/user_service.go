package service

import (
	"gin-chat/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserList
// @Summary 获取用户列表
// @Tags 用户模块
// @Accept json
// @Produce json
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data, _ := models.GetUserList()
	c.JSON(http.StatusOK, gin.H{
		"message": data,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string true "name"
// @param password query string true "password"
// @param repassword query string true "repassword"
// @Accept json
// @Produce json
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Query("name")
	password := c.Query("password")
	repassword := c.Query("repassword")
	if password != repassword {
		c.JSON(http.StatusOK, gin.H{
			"message": "密码不一致",
		})
		return
	}
	user.PassWord = password
	models.CreateUser(&user)
	c.JSON(http.StatusOK, gin.H{
		"message": "新增用户成功",
	})
}
