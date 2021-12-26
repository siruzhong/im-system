package model

const (
	CMD_SINGLE_MSG = 10 // 私聊消息
	CMD_ROOM_MSG   = 11 // 群聊消息
)

// Message 消息结构体
type Message struct {
	Id      int64  `json:"id,omitempty" form:"id"`           // 消息id
	Userid  int64  `json:"userid,omitempty" form:"userid"`   // 消息发送者id
	Cmd     int    `json:"cmd,omitempty" form:"cmd"`         // 消息类型(群聊/私聊)
	Dstid   int64  `json:"dstid,omitempty" form:"dstid"`     // 消息接受者id(好友/群聊)
	Media   int    `json:"media,omitempty" form:"media"`     // 消息按照什么样式展示
	Content string `json:"content,omitempty" form:"content"` // 消息内容
	Pic     string `json:"pic,omitempty" form:"pic"`         // 预览图片
	Url     string `json:"url,omitempty" form:"url"`         // 服务的URL
	Memo    string `json:"memo,omitempty" form:"memo"`       // 消息描述
	Amount  int    `json:"amount,omitempty" form:"amount"`   // 数字相关的(红包单位为分、语音单位为秒)
}

/**
MEDIA_TYPE_TEXT			{id:1,userid:2,dstid:3,cmd:10,media:1,content:"hello"}
MEDIA_TYPE_News			{id:1,userid:2,dstid:3,cmd:10,media:2,content:"标题",pic:"http://www.baidu.com/a/log,jpg",url:"http://www.a,com/dsturl","memo":"这是描述"}
MEDIA_TYPE_VOICE		{id:1,userid:2,dstid:3,cmd:10,media:3,url:"http://www.a,com/dsturl.mp3",anount:40}
MEDIA_TYPE_IMG			{id:1,userid:2,dstid:3,cmd:10,media:4,url:"http://www.baidu.com/a/log,jpg"}
MEDIA_TYPE_REDPACKAGR	{id:1,userid:2,dstid:3,cmd:10,media:5,url:"http://www.baidu.com/a/b/c/redpackageaddress?id=100000","amount":300,"memo":"恭喜发财"}
MEDIA_TYPE_EMOJ			{id:1,userid:2,dstid:3,cmd:10,media:6,"content":"cry"}
MEDIA_TYPE_Link			{id:1,userid:2,dstid:3,cmd:10,media:7,"url":"http://www.a,com/dsturl.html"}
MEDIA_TYPE_VIDEO		{id:1,userid:2,dstid:3,cmd:10,media:8,pic:"http://www.baidu.com/a/log,jpg",url:"http://www.a,com/a.mp4"}
MEDIA_TYPE_CONTACT		{id:1,userid:2,dstid:3,cmd:10,media:9,"content":"10086","pic":"http://www.baidu.com/a/avatar,jpg","memo":"胡大力"}
*/
