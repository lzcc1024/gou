package unet

type Message struct {
	//消息的ID
	Id uint32
	//消息的长度
	DataLen uint32
	//消息的内容
	Data []byte
}

// 创建 Message
func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		Data:    data,
		DataLen: uint32(len(data)),
	}
}

// 获取消息数据段长度
func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

// 获取消息ID
func (m *Message) GetMsgId() uint32 {
	return m.Id
}

// 获取消息内容
func (m *Message) GetData() []byte {
	return m.Data
}

// 设置消息ID
func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

// 设置消息内容
func (m *Message) SetData(data []byte) {
	m.Data = data
}

// 设置消息数据段长度
func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}
