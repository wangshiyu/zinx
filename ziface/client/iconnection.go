package client

import "net"

//定义连接接口
type IConnection interface {
	//启动连接，让当前连接开始工作
	Start()
	//停止连接，结束当前连接状态
	Stop()
	//从当前连接获取原始的socket TCPConn
	GetTCPConnection() *net.Conn

	//直接将Message数据发送数据给远程的TCP服务端(无缓冲)
	SendMsg(msgId int32, data []byte) error
	//直接将Message数据发送给远程的TCP服务端(有缓冲)
	SendBuffMsg(msgId int32, data []byte) error

	//设置链接属性
	SetProperty(key string, value interface{})
	//获取链接属性
	GetProperty(key string) (interface{}, error)
	//移除链接属性
	RemoveProperty(key string)
	//是否授权
	IsAuth() bool
	//是否关闭
	IsClosed() bool
}