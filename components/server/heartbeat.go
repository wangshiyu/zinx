package server

import (
	"github.com/wangshiyu/zinx/ziface/server"
)

type Heartbeat struct {
	TcpServer server.IServer
	Crons     []string
}

func (c *Heartbeat) Run() {
	connectionMap := c.TcpServer.GetConnMgr().Gets()
	if len(connectionMap) > 0 {
		for _, value := range connectionMap {
			value.SendMsg(-1, []byte("ping"))
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

func NewHeartbeat(TcpServer server.IServer) *Heartbeat {
	heartbeat := Heartbeat{}
	heartbeat.TcpServer = TcpServer
	return &heartbeat
}
