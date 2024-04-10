package main

import (
	"fmt"
	"gin-chat/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMysql()

	// utils.DB.AutoMigrate(&models.UserBasic{})

	// utils.DB.AutoMigrate(&models.Message{})

	// utils.DB.AutoMigrate(&models.Contact{})

	// utils.DB.AutoMigrate(&models.GroupBasic{})

	fmt.Println(utils.GetUUID())

}
