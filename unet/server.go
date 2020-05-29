package unet

import (
	"fmt"
	"net"

	"github.com/lzcc1024/gou/uiface"
	"github.com/lzcc1024/gou/utils"
)

// 实现 iServer 接口，定义 Server 服务类
type Server struct {
	// 服务器名称
	Name string
	// 端口
	Port int
	// IP
	IP string
	//  协议 版本 tcp tcp4 tcp6
	Network string

	//当前Server的消息管理模块，用来绑定MsgId和对应的处理方法
	msgHandler uiface.IMsgHandle

	//当前Server的链接管理器
	ConnMgr uiface.IConnManager

	//两个hook函数原型
	//该Server的连接创建时Hook函数
	OnConnStart func(conn uiface.IConnection)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn uiface.IConnection)
}

// 创建 Server
func NewServer() uiface.IServer {
	// 初始化全局配置文件
	utils.GlobalObject.Reload()

	return &Server{
		Name:       utils.GlobalObject.Name,
		Port:       utils.GlobalObject.Port,
		IP:         utils.GlobalObject.Host,
		Network:    "tcp4",
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
}

// 开启服务器
func (s *Server) Start() {
	// 终端输出 starting 提示
	fmt.Printf("[\033[32mInfo\033[0m] server name: \033[4m%s\033[0m, listenner at IP: \033[4m%s\033[0m, port \033[4m%d\033[0m is starting\n", s.Name, s.IP, s.Port)

	fmt.Printf("[\033[32mInfo\033[0m] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)

	// 开启 goroutine 处理服务端Linster业务
	go func() {

		s.msgHandler.StartWorkerPool()

		// 获取 TCP 的 Addr
		tcpAddr, err := net.ResolveTCPAddr(s.Network, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println(err)
			return
		}

		// 监听服务器地址
		listener, err := net.ListenTCP(s.Network, tcpAddr)
		if err != nil {
			fmt.Printf("[\033[31mError\033[0m] %s\n", err)
			return
		}
		// 终端输出 listening 提示
		fmt.Printf("[\033[32mInfo\033[0m] \033[4m%s\033[0m is listening\n", s.Name)

		//TODO server.go 应该有一个自动生成ID的方法
		var cid uint32
		cid = 0

		// 启动 server 网络连接服务
		for {
			// 阻塞等待客户端建立连接
			tcpConn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Printf("[\033[31mError\033[0m] %s\n", err)
				continue
			}

			// 限制最大连接数
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				tcpConn.Close()
				continue
			}

			fmt.Printf("[\033[32mInfo\033[0m] get conn remote addr: %s\n", tcpConn.RemoteAddr().String())
			//处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
			dealConn := NewConnection(s, tcpConn, cid, s.msgHandler)
			cid++

			//启动当前链接的处理业务
			go dealConn.Start()

		}
	}()
}

// 停止服务器
func (s *Server) Stop() {

	s.ConnMgr.ClearConn()

}

// 运行业务服务
func (s *Server) Run() {
	// 开始服务
	s.Start()

	// 阻塞 避免主线程退出
	for {

	}
}

//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRouter(msgId uint32, router uiface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
}

//得到链接管理
func (s *Server) GetConnMgr() uiface.IConnManager {
	return s.ConnMgr
}

//设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(uiface.IConnection)) {
	s.OnConnStart = hookFunc
}

//设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(uiface.IConnection)) {
	s.OnConnStop = hookFunc
}

//调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn uiface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart....")
		s.OnConnStart(conn)
	}

}

//调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn uiface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> CallOnConnStop....")
		s.OnConnStop(conn)
	}
}
