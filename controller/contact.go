package controller

import (
	"IMChat/common/model"
	"IMChat/common/request"
	"IMChat/common/respond"
	"IMChat/common/utils"
	"IMChat/service"
	"log"
	"net/http"
)

// LoadFriend 加载所有好友
func LoadFriend(w http.ResponseWriter, req *http.Request) {
	var arg request.Request
	if err := utils.Bind(req, &arg); err != nil {
		respond.ResponseFail(w, err.Error())
		return
	}
	users := service.ContactService{}.SearchFriend(arg.Userid)
	respond.ResponseOkList(w, users, len(users))
}

// LoadCommunity 加载所有群聊
func LoadCommunity(w http.ResponseWriter, req *http.Request) {
	var arg request.Request
	if err := utils.Bind(req, &arg); err != nil {
		respond.ResponseFail(w, err.Error())
		return
	}
	communities := service.ContactService{}.SearchCommunity(arg.Userid)
	GroupPeopleMap, _ := service.ContactService{}.GetCommunityPeopleNum(communities)
	respond.ResponseOk(w, map[string]interface{}{
		"communities": communities,
		"group_map":   GroupPeopleMap,
	}, "获取社群信息成功")
}

// CreateCommunity 创建群聊
func CreateCommunity(w http.ResponseWriter, req *http.Request) {
	var arg model.Community
	if err := utils.Bind(req, &arg); err != nil {
		log.Println(err)
		respond.ResponseFail(w, err.Error())
		return
	}
	com, err := service.ContactService{}.CreateCommunity(arg)
	if err != nil {
		respond.ResponseFail(w, err.Error())
	} else {
		respond.ResponseOk(w, com, "")
	}
	return
}

// JoinCommunity 加入群聊
func JoinCommunity(w http.ResponseWriter, req *http.Request) {
	var arg request.Request
	if err := utils.Bind(req, &arg); err != nil {
		respond.ResponseFail(w, err.Error())
		return
	}
	err := service.ContactService{}.JoinCommunity(arg.Userid, arg.Dstid)
	AddGroupId(arg.Userid, arg.Dstid)
	if err != nil {
		respond.ResponseFail(w, err.Error())
	} else {
		respond.ResponseOk(w, nil, "群聊加入成功")
	}
}

// AddFriend 添加好友
func AddFriend(w http.ResponseWriter, req *http.Request) {
	var arg request.Request
	if err := utils.Bind(req, &arg); err != nil {
		respond.ResponseFail(w, err.Error())
		return
	}
	err := service.ContactService{}.AddFriend(arg.Userid, arg.Dstid)
	if err != nil {
		respond.ResponseFail(w, err.Error())
	} else {
		respond.ResponseOk(w, nil, "好友添加成功")
	}
}
