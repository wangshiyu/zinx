package components

import (
	"github.com/wangshiyu/zinx/ziface"
)

type Component struct {
	//当前Conn属于哪个Server
	TcpServer ziface.IServer
	//Cron表达式 运行时间
	Crons []string
}

func (c *Component) Run() {

}

func (c *Component) Init() {

}

