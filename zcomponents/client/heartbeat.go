package client

import (
	"fmt"
	"github.com/wangshiyu/zinx/utils"
	"github.com/wangshiyu/zinx/ziface/client"
	"github.com/wangshiyu/zinx/znet"
)

type Heartbeat struct {
	Client client.IClient
	Crons  []string
}

func (c *Heartbeat) Run() {
	Connection := c.Client.GetConnection()
	if utils.IsNotNil(Connection) {
		err := Connection.SendBuffMsg(-1, []byte(znet.PING))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (c *Heartbeat) Init() {
	//一秒运行一次
	c.Crons = []string{"0/5 * * * * *"}
}

func (c *Heartbeat) GetCrons() []string {
	return c.Crons
}

func NewHeartbeat(Client client.IClient) *Heartbeat {
	heartbeat := Heartbeat{}
	heartbeat.Client = Client
	return &heartbeat
}
