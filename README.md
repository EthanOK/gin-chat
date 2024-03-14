`Go 练手项目`

# Gin 项目实战练习

## 一、引入 GORM

- 安装 GORM

```
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```

- 安装 viper

```
go get -u github.com/spf13/viper
```

- 配置 viper

```go
    // 读取config/app.yaml文件
	viper.SetConfigName("app")
    viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
```

- 配置 GORM

```go
	db_, err := gorm.Open(mysql.Open(viper.GetString("mysql.database")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 通过 UserBasic 自动创建 表
	db_.AutoMigrate(&model.UserBasic{})
```

## 二、引入 Gin 框架

- 安装 Gin

```
go get -u github.com/gin-gonic/gin
```

- 配置 Gin

```go
	// package service
	func GetIndex(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome Gin Chat!!!",
	})
}
```

```go
	r := gin.Default()

	r.GET("/index", service.GetIndex)

```

- 启动 Gin

```go
	r.Run(":8080")
```

## 三、引入 gin-swagger

- 安装 gin-swagger

```
go get -u github.com/swaggo/swag/cmd/swag

go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

- 配置 gin-swagger

`service/index.go`

```go
	// GetIndex
	// @Tags 首页
	// @Accept json
	// @Produce json
	// @Success 200 {string} Welcome!!!
	// @Router /index [get]
	func GetIndex(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome Gin Chat!!!",
		})
	}
```

`router/app.go`

```go
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
```

- 生成 dosc files

```
swag init
```

## 四、日志打印

```go
func InitMysql() {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 彩色打印
		},
	)

	db_, err := gorm.Open(mysql.Open(viper.GetString("mysql.database")), &gorm.Config{
		Logger: newLogger,
	})
}
```

## 五、用户模块基本功能

- 创建用户

```go
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

	// 通过用户名查用户信息
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

	// 生成 md5 密码
	user.PassWord = utils.MakePassword(password, salt)
	// UserBasic 表 添加 Salt 字段
	user.Salt = salt
	models.CreateUser(&user)
	c.JSON(http.StatusOK, gin.H{
		"message": "新增用户成功",
	})
}
```

- 获取用户列表

```go
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
```

- 删除用户

```go
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
```

- 更新用户

```go
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
```

- 校验电话与邮箱

```go
type UserBasic struct {
	gorm.Model
	Name          string
	PassWord      string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"` //正则表达式验证手机号码格式是否正确，如果不正确
	Email         string `valid:"email"`
	Identity      string
	ClientIp      string
	ClientPort    string
	LoginTime     sql.NullTime
	HeartbeatTime sql.NullTime
	LoginOutTime  sql.NullTime
	IsLogout      bool
	DeviceInfo    string
}
```

```go
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	// ......
	user.Email = c.PostForm("email")
	// 校验数据
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
		return
	}
	// ......
}
```

- 用户登陆

```go
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

	// 通过用户名查用户信息
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"message": "用户名不存在",
		})
		return
	}

	// 验证密码是否正确
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
```
