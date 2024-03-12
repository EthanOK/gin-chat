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
```

二、引入 Gin 框架
