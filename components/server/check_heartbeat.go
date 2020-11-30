package server

import (
	"github.com/wangshiyu/zinx/ziface/server"
	"github.com/wangshiyu/zinx/znet"
)

type CheckHeartbeat struct {
	TcpServer server.IServer
	Crons     []string
}

func (c *CheckHeartbeat) Run() {
	connectionMap := c.TcpServer.GetConnMgr().Gets()
	if len(connectionMap) > 0 {
		for _, value := range connectionMap {
			var lastTime, _ = value.GetProperty(znet.LAST_MSG_READ_DATE)
			println(lastTime)
		}
	}
}

func (c *CheckHeartbeat) Init() {
	//一秒运行一次
	c.Crons = []string{"0/5 * * * * *"}
}

func (c *CheckHeartbeat) GetCrons() []string {
	return c.Crons
}

func NewCheckHeartbeat(TcpServer server.IServer) *CheckHeartbeat {
	checkHeartbeat := CheckHeartbeat{}
	checkHeartbeat.TcpServer = TcpServer
	return &checkHeartbeat
}