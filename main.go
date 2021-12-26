package main

import (
	"IMChat/initial"
	"log"
	"net/http"
)

func main() {
	// 1. 提供指定目录的静态文件支持
	initial.InitStaticFile()

	// 2. 初始化template
	initial.InitTemplate()

	// 3. 注册路由
	initial.InitRouter()

	// 4. 开启服务
	log.Println("Api server run success on ", "127.0.0.1:8080")
	_ = http.ListenAndServe("127.0.0.1:8080", nil)
}
