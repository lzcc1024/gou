package unet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/lzcc1024/gou/uiface"
	"github.com/lzcc1024/gou/utils"
)

type Connection struct {
	//当前Conn属于哪个Server
	TcpServer uiface.IServer
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	// 当前连接的ID 全局唯一
	ConnID uint32
	//当前连接的关闭状态
	isClosed bool

	//消息管理MsgId和对应处理方法的消息管理模块
	MsgHandler uiface.IMsgHandle

	//通知当前链接退出/停止的channel
	ExitBuffChan chan bool

	//无缓冲管道，用于读、写两个goroutine之间的消息通信
	msgChan chan []byte

	//有关冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte

	//链接属性
	property map[string]interface{}
	//保护链接属性修改的锁
	propertyLock sync.RWMutex
}

// 创建 Connection
func NewConnection(server uiface.IServer, conn *net.TCPConn, connID uint32, msgHandler uiface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:    server,
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		MsgHandler:   msgHandler,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
		msgBuffChan:  make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
		property:     make(map[string]interface{}),
	}

	//将新创建的Conn添加到链接管理中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

func (c *Connection) StartReader() {
	fmt.Printf("[\033[32mInfo\033[0m] ConnID %d reader goroutine is running\n", c.ConnID)
	defer fmt.Printf("[\033[32mInfo\033[0m] \033[4m%s\033[0m conn reader exit!\n", c.GetRemoteAddr().String())
	defer c.Stop()

	for {
		// 创建拆解包的对象
		dp := NewDataPack()
		// 读取客户端的Msg head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Printf("[\033[33mWarn\033[0m] %s\n", err)
			break
		}
		//拆包，得到msgid 和 datalen 放在msg中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Printf("[\033[33mWarn\033[0m] %s\n", err)
			break
		}
		//根据 dataLen 读取 data，放在msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Printf("[\033[33mWarn\033[0m] %s\n", err)
				break
			}
		}
		msg.SetData(data)

		//得到当前客户端请求的Request数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经启动工作池机制，将消息交给Worker处理
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			//从绑定好的消息和对应的处理方法中执行对应的Handle方法
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// 启动连接 当前连接开始工作
func (c *Connection) Start() {
	//开启用户从客户端读取数据流程的Goroutine
	go c.StartReader()
	//开启用于写回客户端数据流程的Goroutine
	go c.StartWriter()

	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.TcpServer.CallOnConnStart(c)

	// 阻塞
	for {
		select {
		case <-c.ExitBuffChan:
			//得到退出消息，不再阻塞
			return
		}
	}
}

// 停止连接 当前连接结束工作
func (c *Connection) Stop() {
	// 当前连接已经关闭 直接返回
	if c.isClosed == true {
		return
	}

	// 更新关闭状态
	c.isClosed = true

	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.TcpServer.CallOnConnStop(c)

	// 关闭当前连接
	c.Conn.Close()

	// 通知缓冲chan
	c.ExitBuffChan <- true

	//将链接从连接管理器中删除
	c.TcpServer.GetConnMgr().Remove(c)

	// 关闭全部chan
	close(c.ExitBuffChan)
	close(c.msgBuffChan)
}

// 获取当前连接的 socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取当前连接远程客户端 Addr
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}
	// 创建拆解包对象
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		return errors.New(fmt.Sprintf("Pack error msg, msgID: %d", msgId))
	}
	// 响应
	c.msgChan <- msg
	return nil
}

// 写消息Goroutine， 用户将数据发送给客户端
func (c *Connection) StartWriter() {
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:, ", err, " Conn Writer exit")
				return
			}
			//针对有缓冲channel需要些的数据处理
		case data, ok := <-c.msgBuffChan:
			if ok {
				//有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				break
			}
		case <-c.ExitBuffChan:
			return
		}
	}
}

//直接将Message数据发送给远程的TCP客户端(有缓冲)
func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {

	if c.isClosed == true {
		return errors.New("Connection closed when send buff msg")
	}
	//将data封包，并且发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}

	//写回客户端
	c.msgBuffChan <- msg

	return nil
}

//设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {

	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value

}

//获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {

	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}

}

//移除链接属性
func (c *Connection) RemoveProperty(key string) {

	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)

}
