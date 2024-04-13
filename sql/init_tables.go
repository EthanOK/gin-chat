package sql

import (
	"gin-chat/models"
	"gin-chat/utils"
)

func InitTables() {
	// 自动迁移模式 没有表自动创建；有表自动添加新字段
	utils.DB.AutoMigrate(&models.UserBasic{})

	utils.DB.AutoMigrate(&models.Message{})

	utils.DB.AutoMigrate(&models.Contact{})

	utils.DB.AutoMigrate(&models.Community{})
}
