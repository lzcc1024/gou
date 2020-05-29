package unet

import (
	"errors"
	"fmt"
	"sync"

	"github.com/lzcc1024/gou/uiface"
)

type ConnManager struct {
	//管理的连接信息
	connections map[uint32]uiface.IConnection
	//读写连接的读写锁
	connLock sync.RWMutex
}

/*
   创建一个链接管理
*/
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]uiface.IConnection),
	}
}

//添加链接
func (c *ConnManager) Add(conn uiface.IConnection) {

	//保护共享资源Map 加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	//将conn连接添加到ConnMananger中
	c.connections[conn.GetConnID()] = conn

	fmt.Println("connection add to ConnManager successfully: conn num = ", c.Len())

}

//删除连接
func (c *ConnManager) Remove(conn uiface.IConnection) {
	//保护共享资源Map 加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	//删除连接信息
	delete(c.connections, conn.GetConnID())

	fmt.Println("connection Remove ConnID=", conn.GetConnID(), " successfully: conn num = ", c.Len())

}

//利用ConnID获取链接
func (c *ConnManager) Get(connID uint32) (uiface.IConnection, error) {

	//保护共享资源Map 加读锁
	c.connLock.RLock()
	defer c.connLock.RUnlock()

	if conn, ok := c.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

//获取当前连接数量
func (c *ConnManager) Len() int {
	return len(c.connections)
}

//删除并停止所有链接
func (c *ConnManager) ClearConn() {

	//保护共享资源Map 加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	//停止并删除全部的连接信息
	for connID, conn := range c.connections {
		//停止
		conn.Stop()
		//删除
		delete(c.connections, connID)
	}

	fmt.Println("Clear All Connections successfully: conn num = ", c.Len())

}
