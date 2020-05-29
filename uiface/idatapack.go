package uiface

// 在发包之前打包成有head和body的两部分的包，在收到数据的时候分两次进行读取，先读取固定长度的head部分，得到后续Data的长度，再根据DataLen读取之后的body

type IDataPack interface {
	//获取包头长度方法
	GetHeadLen() uint32
	//封包方法
	Pack(IMessage) ([]byte, error)
	//拆包方法
	Unpack([]byte) (IMessage, error)
}
