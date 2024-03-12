package models

import (
	"database/sql"
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

	// for _, v := range userList {
	// 	fmt.Println(v)
	// }
	return
}

func CreateUser(user *UserBasic) *gorm.DB {

	return utils.DB.Create(&user)
}

func DeleteUser(user *UserBasic) *gorm.DB {

	return utils.DB.Delete(&user)
}

func UpdateUser(user *UserBasic) *gorm.DB {

	return utils.DB.Model(&user).Updates(UserBasic{Name: user.Name, PassWord: user.PassWord})
}
