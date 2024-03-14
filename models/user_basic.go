package models

import (
	"database/sql"
	"gin-chat/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
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
}

type Claims struct {
	Name     string
	PassWord string
	jwt.RegisteredClaims
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

func FindUserByName(name string) (user UserBasic) {

	utils.DB.Where("name = ?", name).First(&user)
	return
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

	utils.DB.AutoMigrate(&UserBasic{})

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
			Email:    user.Email})
}

func GenerateToken(user *UserBasic) *gorm.DB {

	identity, _ := generateToken(user.Name, user.PassWord)

	return utils.DB.Model(&user).
		Where("id = ?", user.ID).
		Update("identity", identity)

}

func generateToken(name string, password string) (string, error) {

	claims := Claims{
		name,
		password,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "ethan",
		},
	}

	mySigningKey := []byte(viper.GetString("jwt.secret"))

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(mySigningKey)
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("jwt.secret")), nil
	})
	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
