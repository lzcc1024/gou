package unet

import (
	"github.com/lzcc1024/gou/uiface"
)

type Request struct {
	conn uiface.IConnection
	msg  uiface.IMessage
}

// 获取当前请求连接
func (r *Request) GetConnection() uiface.IConnection {
	return r.conn
}

// 获取当前请求消息数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

//获取请求的消息的ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
