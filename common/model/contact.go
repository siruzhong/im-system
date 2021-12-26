package model

import "time"

const (
	CONCAT_CATE_USER     = 0x01 // 好友类型
	CONCAT_CATE_COMUNITY = 0x02 // 群聊类型
)

// Contact 会话结构体
type Contact struct {
	Id       int64     `xorm:"pk autoincr bigint(20)" form:"id" json:"id"` // 聊天会话id
	Ownerid  int64     `xorm:"bigint(20)" form:"ownerid" json:"ownerid"`   // 聊天发起者id
	Dstid    int64     `xorm:"bigint(20)" form:"dstid" json:"dstid"`       // 聊天对话者id
	Cate     int       `xorm:"int(11)" form:"cate" json:"cate"`            // 聊天类型(好友/群聊)
	Memo     string    `xorm:"varchar(120)" form:"memo" json:"memo"`       // 聊天备注
	Createat time.Time `xorm:"datetime" form:"createat" json:"createat"`   // 聊天会话创建时间
}
