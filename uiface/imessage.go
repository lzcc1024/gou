package uiface

// 定义 IMessage 接口
type IMessage interface {
	// 获取消息数据段长度
	GetDataLen() uint32
	// 获取消息ID
	GetMsgId() uint32
	// 获取消息内容
	GetData() []byte

	// 设置消息ID
	SetMsgId(uint32)
	// 设置消息内容
	SetData([]byte)
	// 设置消息数据段长度
	SetDataLen(uint32)
}
