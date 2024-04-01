package models

import (
	"gin-chat/utils"

	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	OwnerId  uint
	TargetId uint
	Type     int //1好友 2群组
	Desc     string
}

func (msg *Contact) TableName() string {
	return "contact"

}

func SearchFriends(userId uint) (friends []UserBasic) {
	var contacts []Contact
	utils.DB.Where("owner_id = ? and type = ?", userId, 1).Find(&contacts)
	ids := make([]uint, 0)

	for _, v := range contacts {
		ids = append(ids, v.TargetId)
	}
	utils.DB.Where("id in (?)", ids).Find(&friends)

	return
}
