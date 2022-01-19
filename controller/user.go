package controller

import (
	"IMChat/common/model"
	"IMChat/common/respond"
	"IMChat/common/utils"
	"IMChat/service"
	"log"
	"net/http"
)

// Login 用户登陆
func Login(writer http.ResponseWriter, request *http.Request) {
	_ = request.ParseForm()
	mobile := request.PostForm.Get("mobile")     // 获取前端mobile参数
	password := request.PostForm.Get("password") // 获取前端参数
	user, err := service.UserService{}.Login(mobile, password)
	if err != nil {
		respond.ResponseFail(writer, err.Error())
	} else {
		respond.ResponseOk(writer, user, "登陆成功")
	}
	return
}

// Register 用户注册
func Register(writer http.ResponseWriter, request *http.Request) {
	_ = request.ParseForm()
	mobile := request.PostForm.Get("mobile")
	password := request.PostForm.Get("password")
	nickname := request.PostForm.Get("nickname")
	avatar := ""
	sex := model.SEX_UNKNOW
	user, err := service.UserService{}.Register(mobile, password, nickname, avatar, sex)
	if err != nil {
		respond.ResponseFail(writer, err.Error())
	} else {
		respond.ResponseOk(writer, user, "注册成功")
	}
	return
}

// ChangeInfo 修改用户信息
func ChangeInfo(writer http.ResponseWriter, request *http.Request) {
	var arg model.User
	if err := utils.Bind(request, &arg); err != nil {
		log.Println(err)
		respond.ResponseFail(writer, err.Error())
		return
	}
	err := service.UserService{}.ModifyUserById(arg)
	if err != nil {
		respond.ResponseFail(writer, err.Error())
	} else {
		respond.ResponseOk(writer, nil, "用户信息修改成功")
	}
	return
}

// GetUser 查找指定用户
func GetUser(writer http.ResponseWriter, request *http.Request) {
	var user model.User
	if err := utils.Bind(request, &user); err != nil {
		respond.ResponseFail(writer, err.Error())
	}
	user, err := service.UserService{}.GetUserById(user.Id)
	if err != nil {
		respond.ResponseFail(writer, err.Error())
	} else {
		respond.ResponseOk(writer, user, "获取指定用户成功")
	}
	return
}
