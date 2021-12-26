package model

import "time"

const (
	SEX_WOMEN  = "W" // 男
	SEX_MAN    = "M" // 女
	SEX_UNKNOW = "U" // 未知性别
)

// User 用户结构体
type User struct {
	Id       int64     `xorm:"pk autoincr bigint(64)" form:"id" json:"id"`  // 用户id(主键自增)
	Mobile   string    `xorm:"varchar(20)" form:"mobile" json:"mobile"`     // 用户电话号码
	Password string    `xorm:"varchar(40)" form:"password" json:"-"`        // 用户密码
	Avatar   string    `xorm:"varchar(150)" form:"avatar" json:"avatar"`    // 用户头像
	Sex      string    `xorm:"varchar(2)" form:"sex" json:"sex"`            // 用户性别
	Nickname string    `xorm:"varchar(20)" form:"nickname" json:"nickname"` // 用户昵称
	Salt     string    `xorm:"varchar(10)" form:"salt" json:"-"`            // 用户盐度加密(首次注册需要)
	Online   int       `xorm:"int(10)" form:"online" json:"online"`         // 用户是否在线
	Token    string    `xorm:"varchar(40)" form:"token" json:"token"`       // 用户登陆token
	Memo     string    `xorm:"varchar(140)" form:"memo" json:"memo"`        // 用户备注
	Createat time.Time `xorm:"datetime" form:"createat" json:"createat"`    // 创建时间
}
