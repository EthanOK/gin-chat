package models

import (
	"database/sql"
	"fmt"
	"gin-chat/utils"

	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name          string
	PassWord      string
	Phone         string
	Email         string
	Identity      string
	ClientIp      string
	ClientPort    string
	LoginTime     sql.NullTime
	HeartbeatTime sql.NullTime
	LoginOutTime  sql.NullTime
	IsLogout      bool
	DeviceInfo    string
}

func (user *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() (userList []*UserBasic, err error) {
	err = utils.DB.Find(&userList).Error

	for _, v := range userList {
		fmt.Println(v)
	}
	return
}
