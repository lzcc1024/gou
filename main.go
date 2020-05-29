package main

import (
	"fmt"

	"github.com/lzcc1024/gou/uiface"
	"github.com/lzcc1024/gou/unet"
)

// https://www.jianshu.com/p/23d07c0a28e5

//一定要继承BaseRouter
type PingRouter struct {
	unet.BaseRouter
}

//Test PreHandle
// func (this *PingRouter) PreHandle(request uiface.IRequest) {
// 	fmt.Println("Call Router PreHandle")
// 	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping ....\n"))
// 	if err != nil {
// 		fmt.Println("call back ping ping ping error")
// 	}
// }

//Test Handle
func (this *PingRouter) Handle(request uiface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	err := request.GetConnection().SendMsg(0, []byte("1haoxiaoxiti\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

//Test PostHandle
// func (this *PingRouter) PostHandle(request uiface.IRequest) {
// 	fmt.Println("Call Router PostHandle")
// 	_, err := request.GetConnection().GetTCPConnection().Write([]byte("After ping .....\n"))
// 	if err != nil {
// 		fmt.Println("call back ping ping ping error")
// 	}
// }

type HelloZinxRouter struct {
	unet.BaseRouter
}

func (this *HelloZinxRouter) Handle(request uiface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("Hello Zinx Router V0.6"))
	if err != nil {
		fmt.Println(err)
	}
}

//创建连接的时候执行
func DoConnectionBegin(conn uiface.IConnection) {
	fmt.Println("DoConnecionBegin is Called ... ")

	conn.SetProperty("Name", "Aceld")
	conn.SetProperty("Home", "https://www.jianshu.com/u/35261429b7f1")

	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		fmt.Println(err)
	}
}

//连接断开的时候执行
func DoConnectionLost(conn uiface.IConnection) {

	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Conn Property Name = ", name)
	}

	fmt.Println("DoConneciotnLost is Called ... ")
}

func main() {
	s := unet.NewServer()

	//注册链接hook回调函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	s.Run()
}
