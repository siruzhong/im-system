// Package initial 初始化操作
package initial

import (
	"IMChat/controller"
	"html/template"
	"log"
	"net/http"
)

// InitStaticFile 初始化静态文件
func InitStaticFile() {
	http.Handle("/asset/", http.FileServer(http.Dir(".")))
	http.Handle("/mnt/", http.FileServer(http.Dir(".")))
}

// InitTemplate 初始化template
func InitTemplate() {
	tpl, err := template.ParseGlob("view/**/*")
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range tpl.Templates() {
		tName := t.Name()
		http.HandleFunc(tName, func(writer http.ResponseWriter, request *http.Request) {
			_ = tpl.ExecuteTemplate(writer, tName, nil)
		})
	}
}

// InitRouter 注册路由
func InitRouter() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, "/user/login.shtml", http.StatusFound)
		return
	})
	http.HandleFunc("/login", controller.Login)                             // 用户登陆
	http.HandleFunc("/register", controller.Register)                       // 用户注册
	http.HandleFunc("/modification", controller.ChangeInfo)                 // 用户信息修改
	http.HandleFunc("/user/find", controller.GetUser)                       // 查找用户
	http.HandleFunc("/contact/loadcommunity", controller.LoadCommunity)     // 加载群聊信息
	http.HandleFunc("/contact/createcommunity", controller.CreateCommunity) // 创建群聊
	http.HandleFunc("/contact/joincommunity", controller.JoinCommunity)     // 加入群聊
	http.HandleFunc("/contact/loadfriend", controller.LoadFriend)           // 加载好友信息
	http.HandleFunc("/contact/addfriend", controller.AddFriend)             // 添加好友
	http.HandleFunc("/chat", controller.Chat)                               // 聊天
	http.HandleFunc("/attach/upload", controller.Upload)                    // 上传
}
