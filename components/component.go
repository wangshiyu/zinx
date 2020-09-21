package components

import (
	"github.com/wangshiyu/zinx/ziface/server"
)

type Component struct {
	//当前Conn属于哪个Server
	TcpServer server.IServer
	//Cron表达式 运行时间
	Crons []string
}

func (c *Component) Run() {

}

func (c *Component) Init() {

}

func (c *Component) GetCrons() []string {
	return c.Crons
}
