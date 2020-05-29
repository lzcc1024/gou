package uiface

type IConnManager interface {
	//添加链接
	Add(conn IConnection)
	//删除连接
	Remove(conn IConnection)
	//利用ConnID获取链接
	Get(connID uint32) (IConnection, error)
	//获取当前连接数量
	Len() int
	//删除并停止所有链接
	ClearConn()
}
