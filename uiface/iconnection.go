package uiface

import "net"

// 定义 Iconnection 连接接口
type IConnection interface {
	// 启动连接 当前连接开始工作
	Start()
	// 停止连接 当前连接结束工作
	Stop()
	// 获取当前连接的 socket TCPConn
	GetTCPConnection() *net.TCPConn
	// 获取当前连接ID
	GetConnID() uint32
	// 获取当前连接远程客户端 Addr
	GetRemoteAddr() net.Addr

	//直接将Message数据发送数据给远程的TCP客户端
	SendMsg(msgId uint32, data []byte) error

	//直接将Message数据发送给远程的TCP客户端(有缓冲)
	SendBuffMsg(msgId uint32, data []byte) error //添加带缓冲发送消息接口

	//设置链接属性
	SetProperty(key string, value interface{})
	//获取链接属性
	GetProperty(key string) (interface{}, error)
	//移除链接属性
	RemoveProperty(key string)
}

// 定义 HandFunc 统计已处理业务
// param *net.TCPConn  socket原生连接
// param []byte        客户端请求的数据
// param int           客户端请求的数据长度
type HandFunc func(*net.TCPConn, []byte, int) error
