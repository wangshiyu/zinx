package client

import (
	"github.com/wangshiyu/zinx/utils"
	"github.com/wangshiyu/zinx/ziface/client"
	"github.com/wangshiyu/zinx/znet"
	"time"
)

type CheckHeartbeat struct {
	Client client.IClient
	Crons  []string
}

func (c *CheckHeartbeat) Run() {
	connection := c.Client.GetConnection()
	now := time.Now()
	if utils.IsNotNil(connection) {
		var lastTime, _ = connection.GetProperty(znet.LAST_MSG_READ_DATE)
		if utils.IsNotNil(lastTime) {
			timeDifference := now.Unix() - lastTime.(time.Time).Unix()
			if timeDifference > 5 {
				connection.Stop()
			}
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

func NewCheckHeartbeat(Client client.IClient) *CheckHeartbeat {
	checkHeartbeat := CheckHeartbeat{}
	checkHeartbeat.Client = Client
	return &checkHeartbeat
}
