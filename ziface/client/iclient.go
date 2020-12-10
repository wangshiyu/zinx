package client

import (
	"github.com/wangshiyu/zinx/ziface"
)

type IClient interface {
	//启动客户端方法
	Start()
	//连接
	Link()
	//停止客户端方法
	Stop()
	//设置该Server的连接创建时Hook函数
	SetOnConnStart(func(IConnection))
	//设置该Server的连接断开时的Hook函数
	SetOnConnStop(func(IConnection))
	//设置连接
	SetConnection(Connection IConnection)
	//调用连接OnConnStart Hook函数
	CallOnConnStart(Connection IConnection)
	//调用连接OnConnStop Hook函数
	CallOnConnStop(Connection IConnection)
	//获取连接
	GetConnection() IConnection
	//获取组件管理器
	GetComponentMgr() ziface.IComponentManager
	//消息处理
	GetMsgHandler() ziface.IMsgHandle
	//获取加密
	GetEncryption() ziface.IEncryption
	//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
	AddRouter(msgId int32, router ziface.IRouter)
	//获取连接地址
	GetAddress() string
}
