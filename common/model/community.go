package model

import "time"

const (
	COMMUNITY_CATE_COM = 0x01
)

// Community 群聊结构体
type Community struct {
	Id       int64     `xorm:"pk autoincr bigint(20)" form:"id" json:"id"` // 群聊id
	Name     string    `xorm:"varchar(30)" form:"name" json:"name"`        // 群聊名称
	Ownerid  int64     `xorm:"bigint(20)" form:"ownerid" json:"ownerid"`   // 群聊创建者
	Icon     string    `xorm:"varchar(250)" form:"icon" json:"icon"`       // 群聊头像
	Cate     int       `xorm:"nt(11)" form:"cate" json:"cate"`             // 群聊类型
	Memo     string    `xorm:"varchar(120)" form:"memo" json:"memo"`       // 群聊描述
	Createat time.Time `xorm:"datetime" form:"createat" json:"createat"`   // 群聊创建时间
}
