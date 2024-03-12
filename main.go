package main

import (
	"gin-chat/router"
	"gin-chat/utils"
)

func main() {

	utils.InitConfig()
	utils.InitMysql()

	r := router.Router()

	r.Run(":8081") // (for windows "localhost:8080")
}