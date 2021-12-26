package service

import (
	"IMChat/common/model"
	"IMChat/dao"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

// ContactService 聊天服务
type ContactService struct{}

// AddFriend 添加好友
func (service ContactService) AddFriend(userid, dstid int64) error {
	// 如果加自己
	if userid == dstid {
		return errors.New("不能添加自己为好友啊")
	}
	var dstUser model.User
	if _, err := dao.DB.Where("id = ?", dstid).Get(&dstUser); err != nil || dstUser.Id == 0 {
		log.Println(err)
		return errors.New("该用户不存在")
	}
	// 判断是否已经加了好友
	tmp := model.Contact{}
	// 查询是否已经是好友（条件的链式操作）
	_, err := dao.DB.Where("ownerid = ?", userid).And("dstid = ?", dstid).And("cate = ?", model.CONCAT_CATE_USER).Get(&tmp)
	if err != nil {
		return err
	}
	// 如果存在记录说明已经是好友了不加
	if tmp.Id > 0 {
		return errors.New("该用户已经被添加过啦")
	}
	// 事务
	session := dao.DB.NewSession()
	_ = session.Begin()
	// 插自己的
	_, e2 := session.InsertOne(model.Contact{
		Ownerid:  userid,
		Dstid:    dstid,
		Cate:     model.CONCAT_CATE_USER,
		Createat: time.Now(),
	})
	// 插对方的
	_, e3 := session.InsertOne(model.Contact{
		Ownerid:  dstid,
		Dstid:    userid,
		Cate:     model.CONCAT_CATE_USER,
		Createat: time.Now(),
	})
	// 没有错误
	if e2 == nil && e3 == nil {
		// 提交事物
		_ = session.Commit()
		return nil
	} else {
		// 回滚事物
		_ = session.Rollback()
		if e2 != nil {
			return e2
		} else {
			return e3
		}
	}
}

// SearchCommunity 查询群聊信息
func (service ContactService) SearchCommunity(userId int64) []model.Community {
	contacts := make([]model.Contact, 0)
	communityIds := make([]int64, 0)
	_ = dao.DB.Where("ownerid = ? and cate = ?", userId, model.CONCAT_CATE_COMUNITY).Find(&contacts)
	for _, v := range contacts {
		communityIds = append(communityIds, v.Dstid)
	}
	communities := make([]model.Community, 0)
	if len(communityIds) == 0 {
		return communities
	}
	_ = dao.DB.In("id", communityIds).Find(&communities)
	return communities
}

// GetCommunityPeopleNum 获取群聊人数
func (service ContactService) GetCommunityPeopleNum(communities []model.Community) (GroupPeopleMap map[int64]int, err error) {
	var communityIds = "("
	for _, row := range communities {
		communityIds += strconv.FormatInt(row.Id, 10) + ","
	}
	communityIds = strings.Trim(communityIds, ",") + ")"
	query := "select dstid, count(*) as num from contact where cate = 2 and dstid in " + communityIds + " GROUP BY dstid"
	results, err := dao.DB.QueryString(query)
	GroupPeopleMap = make(map[int64]int)
	for _, val := range results {
		dstid, _ := strconv.ParseInt(val["dstid"], 10, 64)
		num, _ := strconv.Atoi(val["num"])
		GroupPeopleMap[dstid] = num
	}
	return
}

// SearchCommunityIds 获取用户创建的所有群聊id
func (service ContactService) SearchCommunityIds(userId int64) (communityIds []int64) {
	contacts := make([]model.Contact, 0)
	communityIds = make([]int64, 0)
	_ = dao.DB.Where("ownerid = ? and cate = ?", userId, model.CONCAT_CATE_COMUNITY).Find(&contacts)
	for _, v := range contacts {
		communityIds = append(communityIds, v.Dstid)
	}
	return
}

// JoinCommunity 加入群聊
func (service ContactService) JoinCommunity(userId, comId int64) error {
	cot := model.Contact{
		Ownerid: userId,
		Dstid:   comId,
		Cate:    model.CONCAT_CATE_COMUNITY,
	}
	_, _ = dao.DB.Get(&cot)
	if cot.Id == 0 {
		cot.Createat = time.Now()
		_, err := dao.DB.InsertOne(cot)
		return err
	} else {
		return nil
	}
}

// CreateCommunity 创建群聊
func (service ContactService) CreateCommunity(comm model.Community) (ret model.Community, err error) {
	if len(comm.Name) == 0 {
		err = errors.New("缺少群名称")
		return ret, err
	}
	if comm.Ownerid == 0 {
		err = errors.New("请先登录")
		return ret, err
	}
	com := model.Community{
		Ownerid: comm.Ownerid,
	}
	// 统计该用户创建群聊的数量
	num, err := dao.DB.Count(&com)
	if num > 5 {
		err = errors.New("一个用户最多只能创见5个群")
		return com, err
	} else {
		comm.Createat = time.Now()
		// 开启事务
		session := dao.DB.NewSession()
		_ = session.Begin()
		_, err = session.InsertOne(&comm)
		if err != nil {
			_ = session.Rollback()
			return com, err
		}
		_, err = session.InsertOne(
			model.Contact{
				Ownerid:  comm.Ownerid,
				Dstid:    comm.Id,
				Cate:     model.CONCAT_CATE_COMUNITY,
				Createat: time.Now(),
			})
		if err != nil {
			_ = session.Rollback()
		} else {
			_ = session.Commit()
		}
		return com, err
	}
}

// SearchFriend 查询好友信息
func (service ContactService) SearchFriend(userId int64) []model.User {
	contacts := make([]model.Contact, 0) // 所有聊天
	friendIds := make([]int64, 0)        // 所有好友id
	_ = dao.DB.Where("ownerid = ? and cate = ?", userId, model.CONCAT_CATE_USER).Find(&contacts)
	for _, contact := range contacts {
		friendIds = append(friendIds, contact.Dstid)
	}
	friends := make([]model.User, 0) // 所有好友
	if len(friendIds) == 0 {
		return friends
	}
	_ = dao.DB.In("id", friendIds).Find(&friends)
	return friends
}
