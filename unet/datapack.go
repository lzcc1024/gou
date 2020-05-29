package unet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/lzcc1024/gou/uiface"
	"github.com/lzcc1024/gou/utils"
)

type DataPack struct {
}

//封包拆包实例初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

//获取包头长度方法
func (d *DataPack) GetHeadLen() uint32 {
	// DataLen uint32(4字节) + Id uint32(4字节)
	return 8
}

//封包方法
func (d *DataPack) Pack(msg uiface.IMessage) ([]byte, error) {

	// 创建存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 写dataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	// 写msgId
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	// 写data
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	// 返回数据
	return dataBuff.Bytes(), nil
}

//拆包方法
func (d *DataPack) Unpack(binaryData []byte) (uiface.IMessage, error) {

	//创建 从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	// 只解压head的信息，得到dataLen和msgID
	msg := &Message{}

	// 读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// 读msgId
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断dataLen的长度是否超出允许的最大包长度
	fmt.Println(utils.GlobalObject.MaxPacketSize)
	fmt.Println(msg.DataLen)
	if utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("Too large msg data recieved")
	}
	// 只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return msg, nil
}
