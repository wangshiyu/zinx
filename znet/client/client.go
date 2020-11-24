package client

import (
	"fmt"
	"github.com/wangshiyu/zinx/utils"
	"github.com/wangshiyu/zinx/ziface"
	"github.com/wangshiyu/zinx/ziface/client"
	"github.com/wangshiyu/zinx/znet"
	"net"
)

type Client struct {
	//客户端的名称
	Name string
	//远程服务的IP地址
	IP string
	//远程服务的端口
	Port uint32
	//链接接口
	Connection client.IConnection
	//该Server的连接创建时Hook函数
	OnConnStart func(conn client.IConnection)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn client.IConnection)
	//组件管理
	ComponentManager ziface.IComponentManager
	//加密
	Encryption ziface.IEncryption
	//当前Server的消息管理模块，用来绑定MsgId和对应的处理方法
	MsgHandler ziface.IMsgHandle
}

func (c *Client) Start() {
	address := fmt.Sprintf("%s:%d", c.IP, c.Port)
	c.MsgHandler.StartWorkerPool()
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("resolve tcp addr err: ", err)
		return
	}
	fmt.Println("create tcp addr: ", address)
	c.Connection = NewConntion(c, conn)
	go c.ComponentManager.Runs()
	go c.Connection.Start()
	select {}
}

func (c *Client) Stop() {
	fmt.Println("[STOP] Zinx clinet , name ", c.Name)
	c.Connection.Stop()
}

func (c *Client) GetConnection() client.IConnection {
	return c.Connection
}

//获取组件管理器
func (c *Client) GetComponentMgr() ziface.IComponentManager {
	return c.ComponentManager
}

func (c *Client) GetMsgHandler() ziface.IMsgHandle {
	return c.MsgHandler
}

//获取加密
func (c *Client) GetEncryption() ziface.IEncryption {
	return c.Encryption
}

//设置该Server的连接创建时Hook函数
func (c *Client) SetOnConnStart(hookFunc func(client.IConnection)) {
	c.OnConnStart = hookFunc
}

//设置该Server的连接断开时的Hook函数
func (c *Client) SetOnConnStop(hookFunc func(client.IConnection)) {
	c.OnConnStop = hookFunc
}

//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (c *Client) AddRouter(msgId int32, router ziface.IRouter) {
	c.MsgHandler.AddRouter(msgId, router)
}

func NewClient(Name string) *Client {
	Client := &Client{
		Name:       Name,
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: znet.NewMsgHandle(),
		Encryption: znet.NewRSA2(),
	}
	Client.ComponentManager = NewComponentManager(Client)
	return Client
}
