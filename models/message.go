package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	FromId   uint   // 发送者id
	TargetId uint   // 接收者id
	Type     string // 消息类型 群聊 私聊 广播
	Media    int    // 文字 图片 音频
	Content  string // 内容
	Pic      string // 图片
	Url      string // 链接
	Amount   int    //
}

func (msg *Message) TableName() string {
	return "message"

}
