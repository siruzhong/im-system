package request

// Request 通用请求结构体
type Request struct {
	Userid int64 `json:"userid" form:"userid"`
	Dstid  int64 `json:"dstid" form:"dstid"`
}
