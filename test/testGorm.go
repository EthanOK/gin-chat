package main

import (
	"gin-chat/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/ginchat"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 通过AutoMigrate函数可以快速建表，如果表已经存在不会重复创建。
	// 不支持修改字段和删除字段，可以增加字段。
	db.AutoMigrate(&models.UserBasic{})

	user := models.UserBasic{}
	user.Name = "张三"

	// Create
	db.Create(&user)

	// db.First(&user, 1) // 根据整型主键查找

	// Update
	db.Model(&user).Update("PassWord", "123456")
	// Update - 更新多个字段
	// db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	// db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - 删除 product
	// db.Delete(&product, 1)
}
