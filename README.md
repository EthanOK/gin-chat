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

- 查看接口文档

```
    http://localhost:8080/swagger/index.html
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
// @param avatar formData string false "icon"
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
	user.Avatar = c.PostForm("icon")

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

- Token 生成与验证

```
go get -u github.com/golang-jwt/jwt/v5
```

```go
package utils

type Claims struct {
	Name     string
	PassWord string
	jwt.RegisteredClaims
}

func GenerateToken(name string, password string) (string, error) {

	claims := Claims{
		name,
		password,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "ethan",
		},
	}

	mySigningKey := []byte(viper.GetString("jwt.secret"))

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(mySigningKey)
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("jwt.secret")), nil
	})
	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
```

- 添加好友

```go
// service
func AddFriend(c *gin.Context) {
	userId, _ := strconv.Atoi(c.PostForm("userId"))
	targetName := c.PostForm("targetName")

	massage := models.AddFriendByName(uint(userId), targetName)

	if massage == "" {
		utils.ResponseOK(c.Writer, "Success", "添加好友成功")
	} else {
		utils.ResponseFail(c.Writer, massage)
	}
}
```

```go
// model
func AddFriendByName(userId uint, targetName string) string {
	user := FindUserByName(targetName)
	if user.Name == "" {
		return "好友不存在"

	}
	if user.ID == userId {
		return "不能添加自己"
	}
	// 判断是否已经添加
	var contact Contact
	utils.DB.Where("owner_id = ? and target_id = ? and type = ?", userId, user.ID, 1).First(&contact)
	if contact.ID != 0 {
		return "好友已存在"
	}
	// 保证事务的一致性
	tx := utils.DB.Begin()

	// 在事务中执行第一个操作
	if err := tx.Create(&Contact{
		OwnerId:  userId,
		TargetId: user.ID,
		Type:     1,
	}).Error; err != nil {
		// 如果第一个操作失败，则回滚事务并返回错误
		tx.Rollback()
		return "error"
	}

	// 在事务中执行第二个操作
	if err := tx.Create(&Contact{
		OwnerId:  user.ID,
		TargetId: userId,
		Type:     1,
	}).Error; err != nil {
		// 如果第二个操作失败，则回滚事务并返回错误
		tx.Rollback()
		return "error"
	}

	// 如果两个操作都成功，则提交事务
	tx.Commit()

	return ""
}
```

- 创建群组

```go
// service
func CreateCommunity(c *gin.Context) {

	community := models.Community{}

	community.Name = c.PostForm("name")
	ownerId, _ := strconv.Atoi(c.PostForm("ownerId"))
	community.OwnerId = uint(ownerId)
	community.Icon = c.PostForm("icon")
	community.Desc = c.PostForm("desc")
	category, _ := strconv.Atoi(c.PostForm("cate"))
	community.Category = uint(category)

	code, message := models.CreateCommunity(&community)

	if code == -1 {
		utils.ResponseFail(c.Writer, message)
	} else if code == 0 {

		utils.ResponseOK(c.Writer, "Success", message)
	}

}
```

```go
// model
func CreateCommunity(community *Community) (code int, message string) {
	if community.Name == "" {
		return -1, "群名称不能为空"
	}
	if community.OwnerId == 0 {
		return -1, "群主不能为空"
	}

	// 保证事务的一致性
	tx := utils.DB.Begin()

	// 在事务中执行第一个操作
	if err := tx.Create(&community).Error; err != nil {
		// 如果第一个操作失败，则回滚事务并返回错误
		tx.Rollback()
		return -1, "创建群失败"
	}

	// 在事务中执行第二个操作
	if err := tx.Create(&Contact{
		OwnerId:  community.OwnerId,
		TargetId: community.ID,
		Type:     2,
		Desc:     community.Name,
	}).Error; err != nil {
		// 如果第二个操作失败，则回滚事务并返回错误
		tx.Rollback()
		return -1, "创建群失败"
	}

	// 如果两个操作都成功，则提交事务
	tx.Commit()

	return 0, "群创建成功"
}
```

- 加入群组

```go
// service
func JoinGroup(c *gin.Context) {
	userId, _ := strconv.Atoi(c.PostForm("userId"))

	communityId, _ := strconv.Atoi(c.PostForm("comId"))

	code, message := models.AddCommunityById(uint(userId), uint(communityId))

	if code == -1 {
		utils.ResponseFail(c.Writer, message)
	} else if code == 0 {
		utils.ResponseOK(c.Writer, "Success", message)
	}
}
```

```go
func AddCommunityById(userId, communityId uint) (code int, message string) {
	user := FindUserById(userId)
	if user.Name == "" {
		return -1, "用户不存在"
	}

	community := FindCommunityById(communityId)

	if community.Name == "" {
		return -1, "群组不存在"
	}

	// 判断是否已经添加群组
	var contact Contact
	utils.DB.Where("owner_id = ? and target_id = ? and type = ?", userId, communityId, 2).First(&contact)
	if contact.ID != 0 {
		return -1, "早已经添加了"
	}

	utils.DB.Create(&Contact{
		OwnerId:  userId,
		TargetId: communityId,
		Type:     2,
		Desc:     community.Name,
	})

	return 0, "添加群组成功"

}
```

## 六、WebSocket 实现 消息通信

```go
func Chat(writer http.ResponseWriter, request *http.Request) {

	// 1. 获取参数
	query := request.URL.Query()
	userId, _ := strconv.ParseInt(query.Get("userId"), 10, 64)
	isValid := true

	// 2. 建立一个 WebSocket 连接
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isValid
		},
	}).Upgrade(writer, request, nil)

	if err != nil {
		fmt.Println(err)

		return
	}

	// 3. 定义 WebSocket 节点
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	// 4. userId 与 node绑定 并加锁 谁进来谁在线
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	// 5. 完成发送消息到 WebSocket 逻辑
	go sendProc(node)

	// 6. 完成从 WebSocket 接收消息逻辑
	go receiveProc(node)

	// 7. 后台向登录用户发送欢迎消息
	hello := "欢迎用户" + FindNameByUserId(uint(userId)) + "进入聊天室"
	sendMsg(uint(userId), []byte(hello))

}
```

### 1. 接受 WebSocket 消息的处理逻辑

```go
func receiveProc(node *Node) {
	for {
		// 从 Websocket 中读取一条消息
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		// 处理接受到的数据
		disPatch(data)
		// 广播数据
		broadMsg(data)

		fmt.Println("[ws] receiveProc <<<<", string(data))

	}
}
```

### 2. 发送 WebSocket 消息的处理逻辑

```go
func sendProc(node *Node) {
	// 从 node.DataQueue 通道中读取数据，并将其发送到 WebSocket 中
	for {
		select {
		case data := <-node.DataQueue:
			fmt.Println("[ws] sendProc >>>>", string(data))
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
```

### 3. 处理接受到的数据

```go
func disPatch(data []byte) {
	msg := Message{}
	// json 转换
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch msg.Type {
	case 1: // 私聊
		sendMsg(msg.TargetId, data)
		// case 2:
		// 	sendGroupMsg(msg)
		// case 3:
		// 	broadMsg(data)
	}
}

func sendMsg(userId uint, msg []byte) {
	fmt.Println("sendMsg >>> userId:", userId, "message:", string(msg))

	rwLocker.RLock()
	node, ok := clientMap[int64(userId)]
	rwLocker.RUnlock()

	if ok {
		node.DataQueue <- msg
	}
}
```

## 七、Redis 实现 消息存储

```go
// 初始化 Redis 连接
```
