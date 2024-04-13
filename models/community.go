package models

import (
	"gin-chat/utils"

	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	Name     string
	OwnerId  uint
	Category uint
	Icon     string
	Desc     string
}

func (table *Community) TableName() string {
	return "community"
}

func CreateCommunity(community *Community) (code int, message string) {
	if community.Name == "" {
		return -1, "群名称不能为空"
	}
	if community.OwnerId == 0 {
		return -1, "群主不能为空"
	}

	// 保证事务的一致性
	tx := utils.DB.Begin()

	// 在事务中执行第一个操作
	if err := tx.Create(&community).Error; err != nil {
		// 如果第一个操作失败，则回滚事务并返回错误
		tx.Rollback()
		return -1, "创建群失败"
	}

	// 在事务中执行第二个操作
	if err := tx.Create(&Contact{
		OwnerId:  community.OwnerId,
		TargetId: community.ID,
		Type:     2,
		Desc:     community.Name,
	}).Error; err != nil {
		// 如果第二个操作失败，则回滚事务并返回错误
		tx.Rollback()
		return -1, "创建群失败"
	}

	// 如果两个操作都成功，则提交事务
	tx.Commit()

	return 0, "群创建成功"
}

func GetCommunityList(userId uint) (communityList []*Community, code int, message string) {
	communityIds := GetCommunityIds(userId)

	if len(communityIds) != 0 {
		if err := utils.DB.Where("id in ?", communityIds).Find(&communityList).Error; err != nil {
			return nil, -1, "获取群列表失败"
		}
	}

	return communityList, 0, "获取群列表成功"

}

func FindCommunityById(id uint) (community Community) {
	utils.DB.Where("id = ?", id).First(&community)
	return
}
