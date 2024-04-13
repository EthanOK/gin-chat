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
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"` //正则表达式验证手机号码格式是否正确，如果不正确
	Email         string `valid:"email"`
	Identity      string
	ClientIp      string
	ClientPort    string
	Salt          string
	LoginTime     sql.NullTime
	HeartbeatTime sql.NullTime
	LoginOutTime  sql.NullTime
	IsLogout      bool
	DeviceInfo    string
	Avatar        string
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

func FindUserByNameAndPwd(name string, password string) (user UserBasic) {
	utils.DB.Where("name = ? and password = ?", name, password).First(&user)
	return
}

func FindUserByIDAndIdentity(userId uint, identity string) (user UserBasic) {
	utils.DB.Where("id = ? and identity = ?", userId, identity).First(&user)
	return

}

func FindUserByName(name string) (user UserBasic) {

	utils.DB.Where("name = ?", name).First(&user)
	return
}

func FindUserById(id uint) (user UserBasic) {

	utils.DB.Where("id = ?", id).First(&user)
	return
}

func FindNameByUserId(id uint) string {
	var user UserBasic
	utils.DB.Where("id = ?", id).First(&user)
	return user.Name
}

func FindUserByPhone(phone string) (user UserBasic) {

	utils.DB.Where("phone = ?", phone).First(&user)
	return
}
func FindUserByEmail(email string) (user UserBasic) {
	utils.DB.Where("email = ?", email).First(&user)
	return
}

func CreateUser(user *UserBasic) *gorm.DB {

	return utils.DB.Create(&user)
}

func DeleteUser(user *UserBasic) *gorm.DB {

	return utils.DB.Delete(&user)
}

func UpdateUser(user *UserBasic) *gorm.DB {

	return utils.DB.Model(&user).Updates(
		UserBasic{
			Name:     user.Name,
			PassWord: user.PassWord,
			Phone:    user.Phone,
			Email:    user.Email,
			Avatar:   user.Avatar,
		})
}

func UpdateAvatarByUserId(userId uint, avatar string) *gorm.DB {
	return utils.DB.Model(&UserBasic{}).Where("id = ?", userId).Update("avatar", avatar)
}

func GenerateToken(user *UserBasic) *gorm.DB {

	identity, _ := utils.GenerateToken(user.Name, user.PassWord)

	return utils.DB.Model(&user).
		Where("id = ?", user.ID).
		Update("identity", identity)

}
