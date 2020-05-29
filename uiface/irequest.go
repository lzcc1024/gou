package uiface

// 定义 IRequest 接口
type IRequest interface {
	// 获取当前请求连接
	GetConnection() IConnection
	// 获取当前请求消息数据
	GetData() []byte
	//获取请求的消息的ID
	GetMsgID() uint32
}
