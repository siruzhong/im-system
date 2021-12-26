package respond

import (
	"encoding/json"
	"log"
	"net/http"
)

// Respond 通用响应结构体
type Respond struct {
	Code  int         `json:"code"`            // 响应码
	Msg   string      `json:"msg"`             // 响应消息
	Data  interface{} `json:"data,omitempty"`  // 响应数据
	Rows  interface{} `json:"rows,omitempty"`  // 所有响应行
	Total interface{} `json:"total,omitempty"` // 响应数据条数
}

// ResponseFail 通用错误响应方法
func ResponseFail(w http.ResponseWriter, msg string) {
	Response(w, -1, msg, nil)
}

// ResponseOk 通用成功响应方法
func ResponseOk(w http.ResponseWriter, data interface{}, msg string) {
	Response(w, 0, msg, data)
}

// ResponseOkList 通用响应多条数据成功方法
func ResponseOkList(w http.ResponseWriter, lists interface{}, total interface{}) {
	//分页数目,
	ResponseList(w, 0, lists, total)
}

// Response 通用响应单数据方法
func Response(w http.ResponseWriter, code int, msg string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res := Respond{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	ret, err := json.Marshal(res)
	if err != nil {
		log.Println(err.Error())
	}
	_, _ = w.Write(ret)
}

// ResponseList 响应多条数据方法
func ResponseList(w http.ResponseWriter, code int, data interface{}, total interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// 满足条件的全部记录数目
	res := Respond{
		Code:  code,
		Rows:  data,
		Total: total,
	}
	ret, err := json.Marshal(res)
	if err != nil {
		log.Println(err.Error())
	}
	_, _ = w.Write(ret)
}
