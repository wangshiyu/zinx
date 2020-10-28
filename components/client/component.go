package client

import "github.com/wangshiyu/zinx/ziface/client"

type ClientComponent struct {
	//当前Conn属于哪个Server
	Client client.IClient
	//Cron表达式 运行时间
	Crons []string
}

func (c *ClientComponent) Run() {

}

func (c *ClientComponent) Init() {

}

func (c *ClientComponent) GetCrons() []string {
	return c.Crons
}
