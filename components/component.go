package components

import (
	"github.com/wangshiyu/zinx/ziface/server"
)

type ServerComponent struct {
	//当前Conn属于哪个Server
	TcpServer server.IServer
	//Cron表达式 运行时间
	Crons []string
}

func (c *ServerComponent) Run() {

}

func (c *ServerComponent) Init() {

}

func (c *ServerComponent) GetCrons() []string {
	return c.Crons
}
