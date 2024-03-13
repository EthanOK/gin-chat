package service

import (
	"fmt"
	"gin-chat/models"
	"gin-chat/utils"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
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
// @param name formData string true "name"
// @param password formData string true "password"
// @param repassword formData string true "repassword"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [post]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.PostForm("name")
	password := c.PostForm("password")
	repassword := c.PostForm("repassword")
	salt := fmt.Sprintf("%06d", rand.Int31())

	data := models.FindUserByName(user.Name)
	if data.Name != "" {
		c.JSON(http.StatusOK, gin.H{
			"message": "用户名已存在",
		})
		return
	}

	if password != repassword {
		c.JSON(http.StatusOK, gin.H{
			"message": "密码不一致",
		})
		return
	}

	user.PassWord = utils.MakePassword(password, salt)
	user.Salt = salt
	models.CreateUser(&user)
	c.JSON(http.StatusOK, gin.H{
		"message": "新增用户成功",
	})
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id formData string true "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [post]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	models.DeleteUser(&user)
	c.JSON(http.StatusOK, gin.H{
		"message": "删除用户成功",
	})
}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string true "id"
// @param name formData string false "name"
// @param password formData string fase "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
		return
	}

	models.UpdateUser(&user)

	c.JSON(http.StatusOK, gin.H{
		"message": "更新用户成功",
	})
}

// LoginUser
// @Summary 用户登陆
// @Tags 用户模块
// @param name formData string true "name"
// @param password formData string true "password"
// @Success 200 {string} json{"code","message"}
// @Router /user/loginUser [post]
func LoginUser(c *gin.Context) {

	name := c.PostForm("name")
	plainpwd := c.PostForm("password")

	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"message": "用户名不存在",
		})
		return
	}
	valid := utils.ValidPassword(plainpwd, user.Salt, user.PassWord)

	if !valid {
		c.JSON(http.StatusOK, gin.H{
			"message": "密码错误",
		})
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登陆成功",
	})
}
