package utils

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	// fmt.Println(viper.Get("mysql"))

}

func InitMysql() {

	db_, err := gorm.Open(mysql.Open(viper.GetString("mysql.database")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	DB = db_

}
