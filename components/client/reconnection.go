package client

import (
	"github.com/wangshiyu/zinx/utils"
	"github.com/wangshiyu/zinx/ziface/client"
)

type Reconnection struct {
	Client client.IClient
	Crons  []string
}

func (c *Reconnection) Run() {
	Connection := c.Client.GetConnection()
	if utils.IsNil(Connection) {
		c.Client.Link()
	}
}

func (c *Reconnection) Init() {
	//一秒运行一次
	c.Crons = []string{"0/5 * * * * *"}
}

func (c *Reconnection) GetCrons() []string {
	return c.Crons
}

func NewReconnection(Client client.IClient) *Reconnection {
	reconnection := Reconnection{}
	reconnection.Client = Client
	return &reconnection
}
