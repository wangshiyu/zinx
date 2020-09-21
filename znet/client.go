package znet

import (
	"github.com/wangshiyu/zinx/ziface"
	ziface2 "github.com/wangshiyu/zinx/ziface/server"
)

type Client struct {
	//客户端的名称
	Name string
	//tcp4 or other
	IPVersion string
	//远程服务的IP地址
	IP string
	//远程服务的端口
	Port uint32
	//组件管理
	ComponentManager ziface.IComponentManager
	//加密
	Encryption ziface.IEncryption
	//链接接口
	Connection ziface2.IConnection
}

func (c *Client) Start() {
	//conn, err := net.Dial("tcp", c.IP+":"+string(c.Port))
}

func NewClient() *Client {
	return &Client{}
}
