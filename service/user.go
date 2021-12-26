package service

import (
	"IMChat/common/model"
	"IMChat/common/utils"
	"IMChat/dao"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// UserService 用户服务
type UserService struct{}

// Login 用户登陆
func (u UserService) Login(mobile, plainpwd string) (user model.User, err error) {
	// 用户查询
	_, err = dao.DB.Where("mobile = ?", mobile).Get(&user)
	if err != nil {
		err = errors.New("该用户不存在")
		return
	}
	// 密码校验
	if !utils.ValidatePasswd(plainpwd, user.Salt, user.Password) {
		err = errors.New("密码不正确")
		return
	}
	// 创建登陆token
	user.Token = utils.GenerateToken()
	// 更新用户的token的字段
	_, err = dao.DB.Id(user.Id).Cols("token").Update(&user)
	return
}

// Register 用户注册
func (u UserService) Register(mobile, plainpwd, nickname, avatar, sex string) (user model.User, err error) {
	// 查询当前手机号的用户
	_, err = dao.DB.Where("mobile = ?", mobile).Get(&user)
	if err != nil {
		return
	}
	if user.Id > 0 {
		err = errors.New("该手机号已经被注册")
		return
	}
	// 创建新用户
	user.Mobile = mobile
	user.Avatar = avatar
	user.Nickname = nickname
	user.Sex = sex
	user.Salt = fmt.Sprintf("%06d", rand.Int31n(10000))   // 生成[0,10000)的随机数
	user.Password = utils.MakePasswd(plainpwd, user.Salt) // 创建加密密码
	user.Createat = time.Now()
	user.Token = utils.GenerateToken() // 创建用户token
	// 数据库中插入新用户
	_, err = dao.DB.InsertOne(&user)
	return user, nil
}

// GetUserById 查询指定用户
func (s UserService) GetUserById(userId int64) (user model.User, err error) {
	// 首先通过手机号查询用户
	_, err = dao.DB.ID(userId).Get(&user)
	return
}

// ModifyUserById 修改指定用户信息
func (s UserService) ModifyUserById(arg model.User) (err error) {
	// 更新用户信息
	_, err = dao.DB.ID(arg.Id).Update(&arg)
	return
}
