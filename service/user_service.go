package service

import (
	"fmt"
	"gin-chat/models"
	"gin-chat/utils"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
		"code":    0,
		"message": "获取成功",
		"data":    data,
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
	if user.Name == "" || password == "" || repassword == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "参数不能为空",
		})
		return
	}

	// 通过用户名查用户信息
	data := models.FindUserByName(user.Name)
	if data.Name != "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "用户名已存在",
		})
		return
	}

	if password != repassword {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "密码不一致",
		})
		return
	}

	// 生成 md5 密码
	user.PassWord = utils.MakePassword(password, salt)
	user.Salt = salt
	models.CreateUser(&user)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "新增用户成功",
		"data":    user,
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
		"code":    0,
		"message": "删除用户成功",
		"data":    user,
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
			"code":    -1,
			"message": err.Error(),
		})
		return
	}

	models.UpdateUser(&user)

	c.JSON(http.StatusOK, gin.H{

		"code":    0,
		"message": "修改用户成功",
		"data":    user,
	})
}

// LoginUser
// @Summary 用户登陆
// @Tags 用户模块
// @param name formData string true "name"
// @param password formData string true "password"
// @Success 200 {string} json{"code","message"}
// @Router /user/login [post]
func LoginUser(c *gin.Context) {

	name := c.PostForm("name")
	plainpwd := c.PostForm("password")

	// 通过用户名查用户信息
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "用户名不存在",
		})
		return
	}

	// 验证密码是否正确
	valid := utils.ValidPassword(plainpwd, user.Salt, user.PassWord)

	if !valid {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "密码错误",
		})
		return

	}

	// 生成 token
	models.GenerateToken(&user)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登陆成功",
		"data":    user,
	})
}

// GetUserByToken
// @Summary 解析Token
// @Tags 测试解析Token
// @param token query string true "token"
// @Success 200 {string} json{"code","message","data"}
// @Router /getUserByToken [get]
func GetUserByToken(c *gin.Context) {
	token := c.Query("token")
	claims, err := utils.ParseToken(token)

	if err != nil {

		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "token 不正确",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "解析成功",
		"data":    claims,
	})
}

var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMessage(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			fmt.Println(err)
		}

	}(ws)

	MsgHander(ws, c)

}

func MsgHander(ws *websocket.Conn, c *gin.Context) {

	msg, err := utils.Subscribe(c, utils.PublishChannel)
	if err != nil {
		fmt.Println(err)
	}
	tm := time.Now().Format("2006-01-02 15:04:05")
	m := fmt.Sprintf("[ws][%s] %s", tm, msg)

	err = ws.WriteMessage(websocket.TextMessage, []byte(m))
	if err != nil {
		fmt.Println(err)
	}
}

func SendUserMessage(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

func SearchFriends(c *gin.Context) {
	userId, _ := strconv.Atoi(c.PostForm("userId"))

	friends := models.SearchFriends(uint(userId))

	// c.JSON(http.StatusOK, gin.H{
	// 	"code":    0,
	// 	"message": "查找好友列表成功",
	// 	"data":    friends,
	// })
	utils.ResponseOKList(c.Writer, friends, len(friends))
}
