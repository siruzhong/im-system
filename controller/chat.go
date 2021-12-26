package controller

import (
	"IMChat/common/model"
	"IMChat/service"
	"encoding/json"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
)

func init() {
	go udpSendProc()
	go udpReceiveProc()
}

// Conversation 会话结构体
type Conversation struct {
	Conn         *websocket.Conn // websocket连接
	DataChan     chan []byte     // 数据通道
	CommunityIds set.Interface   // 用户创建的所有群聊id
}

var (
	clientMap   = make(map[int64]*Conversation, 0) // userid和Conversation映射关系表
	rwLock      sync.RWMutex                       // 读写锁
	udpSendChan = make(chan []byte, 1024)          // 用来存放发送的要广播的数据
)

// Chat 聊天
func Chat(writer http.ResponseWriter, request *http.Request) {
	// 检验接入是否合法,例如：http://127.0.0.1/chat?id=1&token=xxx
	query := request.URL.Query()
	id := query.Get("id")
	token := query.Get("token")
	userId, _ := strconv.ParseInt(id, 10, 64)
	// 检查用户token是否有效
	isValidate := checkToken(userId, token)
	// websocket升级
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { // CheckOrigin解决跨域问题
		return isValidate
	}}).Upgrade(writer, request, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	// 创建会话
	conversation := &Conversation{
		Conn:         conn,
		DataChan:     make(chan []byte, 50),
		CommunityIds: set.New(set.ThreadSafe),
	}
	// 获取用户创建的所有群聊id
	communityIds := service.ContactService{}.SearchCommunityIds(userId)
	for _, v := range communityIds {
		conversation.CommunityIds.Add(v)
	}
	// userid和Conversation形成绑定关系
	rwLock.Lock()
	clientMap[userId] = conversation
	rwLock.Unlock()
	// 完成发送逻辑
	go sendProc(conversation)
	// 完成接收逻辑
	go receiveProc(conversation)
	log.Printf("<-%d\n", userId)
	sendMsg(userId, []byte("hello,world!"))
}

// AddGroupId 添加新的群ID到用户的CommunityIds中
func AddGroupId(userId, gid int64) {
	// 取得node
	rwLock.Lock()
	node, ok := clientMap[userId]
	if ok {
		node.CommunityIds.Add(gid)
	}
	// clientMap[userId] = node
	rwLock.Unlock()
	//添加gid到set
}

// sendProc websocket发送协程
func sendProc(node *Conversation) {
	for {
		select {
		// 循环读取会话中的消息发送到conn中
		case data := <-node.DataChan:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}

// receiveProc websocket接收协程
func receiveProc(node *Conversation) {
	for {
		// 循环读取conn中的消息
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			return
		}
		// 把消息广播到局域网
		broadMsg(data)
		log.Printf("[ws]<=%s\n", data)
	}
}

// broadMsg 将消息广播到局域网
func broadMsg(data []byte) {
	udpSendChan <- data
}

// udpSendProc 完成udp数据的发送协程
func udpSendProc() {
	log.Println("start udpSendProc")
	// 使用udp协议拨号
	conn, err := net.DialUDP("udp", nil,
		&net.UDPAddr{
			IP:   net.IPv4(0, 0, 0, 0),
			Port: 3000,
		})
	defer conn.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
	// 通过得到的conn发送消息
	for {
		select {
		case data := <-udpSendChan:
			_, err = conn.Write(data)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}

// udpReceiveProc 完成upd接收并处理功能
func udpReceiveProc() {
	log.Println("start udpReceiveProc")
	// 监听udp广播端口
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		log.Println(err.Error())
	}
	// 处理端口发过来的数据
	for {
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			log.Println(err.Error())
			return
		}
		// 直接数据处理
		dispatch(buf[0:n])
	}
}

// dispatch 后端调度逻辑处理
func dispatch(data []byte) {
	// 将data反序列化为message
	msg := model.Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Println(err.Error())
		return
	}
	// 根据消息类型对消息进行处理
	switch msg.Cmd {
	// 私聊消息
	case model.CMD_SINGLE_MSG:
		sendMsg(msg.Dstid, data)
	// 群聊消息
	case model.CMD_ROOM_MSG:
		for _, v := range clientMap {
			if v.CommunityIds.Has(msg.Dstid) {
				v.DataChan <- data
			}
		}
	}
}

// sendMsg 发送消息
func sendMsg(userId int64, msg []byte) {
	rwLock.RLock()
	node, ok := clientMap[userId]
	rwLock.RUnlock()
	if ok {
		node.DataChan <- msg
	}
}

// checkToken 检测token是否有效
func checkToken(userId int64, token string) bool {
	user, err := service.UserService{}.GetUserById(userId)
	if err != nil {
		return false
	}
	return user.Token == token
}
