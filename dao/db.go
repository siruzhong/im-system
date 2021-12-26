package dao

import (
	"IMChat/common/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
)

var (
	Driver = "mysql"
	DsName = "root:4389589zsr@tcp(127.0.0.1:3306)/IMCHat?charset=utf8"
	DB     *xorm.Engine
	DBErr  error
)

// init 初始化mysql连接
func init() {
	DB, DBErr = xorm.NewEngine(Driver, DsName)
	if DBErr != nil {
		log.Fatal(DBErr)
	}
	// 设置最大打开连接数
	DB.SetMaxOpenConns(2)
	// 将结构体转换为数据表
	_ = DB.Sync2(
		new(model.User),
		new(model.Community),
		new(model.Contact),
	)
	log.Println("initial database success")
}
