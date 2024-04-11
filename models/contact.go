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

func AddFriend(userId, targetId uint) bool {
	user := FindUserById(targetId)
	if user.Name == "" {
		return false
	}

	utils.DB.Create(&Contact{
		OwnerId:  userId,
		TargetId: targetId,
		Type:     1,
	})
	return true

}
func AddFriendByName(userId uint, targetName string) string {
	user := FindUserByName(targetName)
	if user.Name == "" {
		return "好友不存在"

	}
	if user.ID == userId {
		return "不能添加自己"
	}
	// 判断是否已经添加
	var contact Contact
	utils.DB.Where("owner_id = ? and target_id = ? and type = ?", userId, user.ID, 1).First(&contact)
	if contact.ID != 0 {
		return "好友已存在"
	}
	// 保证事务的一致性
	tx := utils.DB.Begin()

	// 在事务中执行第一个操作
	if err := tx.Create(&Contact{
		OwnerId:  userId,
		TargetId: user.ID,
		Type:     1,
	}).Error; err != nil {
		// 如果第一个操作失败，则回滚事务并返回错误
		tx.Rollback()
		return "error"
	}

	// 在事务中执行第二个操作
	if err := tx.Create(&Contact{
		OwnerId:  user.ID,
		TargetId: userId,
		Type:     1,
	}).Error; err != nil {
		// 如果第二个操作失败，则回滚事务并返回错误
		tx.Rollback()
		return "error"
	}

	// 如果两个操作都成功，则提交事务
	tx.Commit()

	return ""
}
