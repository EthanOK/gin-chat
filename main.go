package main

import (
	"gin-chat/router"
	"gin-chat/sql"
	"gin-chat/utils"
)

func main() {

	utils.InitConfig()

	utils.InitMysql()

	utils.InitRedis()

	sql.InitTables()

	r := router.Router()

	r.Run(":8080") // (for windows "localhost:8080")
}
