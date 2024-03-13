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
