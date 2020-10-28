package client

import (
	"github.com/wangshiyu/zinx/ziface/client"
)

type Heartbeat struct {
	Client client.IClient
	Crons  []string
}

func (c *Heartbeat) Run() {
	Connection := c.Client.GetConnection()
	if Connection != nil {
		Connection.SendBuffMsg(-1, []byte("ping"))
	}
}

func (c *Heartbeat) Init() {
	//一秒运行一次
	c.Crons = []string{"* * * * * *"}
}

func (c *Heartbeat) GetCrons() []string {
	return c.Crons
}

func NewHeartbeat(Client client.IClient) *Heartbeat {
	heartbeat := Heartbeat{}
	heartbeat.Client = Client
	return &heartbeat
}
