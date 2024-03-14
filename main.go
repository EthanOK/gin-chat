package main

import (
	"gin-chat/router"
	"gin-chat/utils"
)

func main() {

	utils.InitConfig()

	utils.InitMysql()

	utils.InitRedis()

	r := router.Router()

	r.Run(":8080") // (for windows "localhost:8080")
}
