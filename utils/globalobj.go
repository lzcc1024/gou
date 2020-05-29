package utils

import (
	"encoding/json"
	"io/ioutil"

	"github.com/lzcc1024/gou/uiface"
)

// 定义 GlobalObj 全局对象结构体
type GlobalObj struct {
	// 全局Server对象
	TcpServer uiface.IServer
	// 服务器主机
	Host string
	// 服务器端口
	Port int
	// 服务器名称
	Name string
	// 当前框架版本
	Version string
	// 允许最大连接数
	MaxConn int
	// 数据包的最大值
	MaxPacketSize uint32
	// 业务工作Worker池的数量
	WorkerPoolSize uint32
	//业务工作Worker对应负责的任务队列最大任务存储数量
	MaxWorkerTaskLen uint32

	MaxMsgChanLen uint32

	// 配置文件路径
	ConfFilePath string
}

func (g *GlobalObj) Reload() {
	// 读取配置文件
	data, err := ioutil.ReadFile(g.ConfFilePath)
	if err != nil {
		panic(err)
	}
	// json解析
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 定义 GlobalObject 全局对象
var GlobalObject *GlobalObj

// 通过模块的 init 方法 初始化 GlobalObject 全局对象
func init() {
	// 初始化 GlobalObject 默认值
	GlobalObject = &GlobalObj{
		Name:             "gouServer",
		Version:          "V0.0.1",
		Host:             "0.0.0.0",
		Port:             1510,
		MaxConn:          20000,
		MaxPacketSize:    4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    1024,
		ConfFilePath:     "conf/config.json",
	}

	// 通过配置文件覆盖 GlobalObject 的值
	GlobalObject.Reload()
}
