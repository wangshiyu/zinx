package server

import (
	"errors"
	"fmt"
	"github.com/wangshiyu/zinx/ziface/server"
	"github.com/wangshiyu/zinx/zlog"
	"sync"
)

/*
	连接管理模块
*/
type ConnManager struct {
	connections map[string]server.IConnection //管理的连接信息
	connLock    sync.RWMutex                  //读写连接的读写锁
}

/*
	创建一个链接管理
*/
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[string]server.IConnection),
	}
}

//添加链接
func (connMgr *ConnManager) Add(conn server.IConnection) {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//将conn连接添加到ConnMananger中
	connMgr.connections[conn.GetConnName()] = conn

	fmt.Println("connection add to ConnManager successfully: conn num = ", connMgr.Len())
}

//删除连接
func (connMgr *ConnManager) Remove(conn server.IConnection) {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除连接信息
	delete(connMgr.connections, conn.GetConnName())
	//fmt.Println("connection Remove ConnID=", conn.GetConnID(), " successfully: conn num = ", connMgr.Len())
	zlog.Info("connection Remove ConnName=", conn.GetConnName(), " successfully: conn num = ", connMgr.Len())
}

//利用ConnID获取链接
func (connMgr *ConnManager) Get(ConnName string) (server.IConnection, error) {
	//保护共享资源Map 加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[ConnName]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

//获取全部链接
func (connMgr *ConnManager) Gets() map[string]server.IConnection {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()
	return connMgr.connections
}

//获取当前连接
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

//清除并停止所有连接
func (connMgr *ConnManager) ClearConn() {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//停止并删除全部的连接信息
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()
		//删除
		delete(connMgr.connections, connID)
	}

	fmt.Println("Clear All Connections successfully: conn num = ", connMgr.Len())
}
