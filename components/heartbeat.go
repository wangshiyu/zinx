package components

import (
	"github.com/wangshiyu/zinx/ziface/server"
)

type Heartbeat struct {
	Component
}

func (c *Heartbeat) Run() {
	connectionMap := c.TcpServer.GetConnMgr().Gets()
	if len(connectionMap)>0 {
		for _, value := range connectionMap {
			value.SendBuffMsg(-1,[]byte("ping"))
		}
	}
}

func (c *Heartbeat) Init() {
	//一秒运行一次
	c.Crons = []string{"* * * * * *"}
}

func NewHeartbeat(TcpServer server.IServer) *Heartbeat {
	heartbeat := Heartbeat{}
	heartbeat.TcpServer = TcpServer
	return &heartbeat
}
