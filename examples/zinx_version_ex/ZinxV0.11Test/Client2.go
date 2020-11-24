package main

import (
	"github.com/wangshiyu/zinx/ziface"
	"github.com/wangshiyu/zinx/znet"
	"github.com/wangshiyu/zinx/znet/client"
)


//ping test 自定义路由
type PingRouter_ struct {
	znet.BaseRouter
}

//Ping Handle
func (this *PingRouter_) Handle(request ziface.IRequest) {

	//zlog.Debug("Call PingRouter Handle")
	////先读取客户端的数据，再回写ping...ping...ping
	//zlog.Debug("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
	//err := request.GetConnection().(ziface2.IConnection).SendBuffMsg(0, []byte("ping...ping...ping"))
	//if err != nil {
	//	zlog.Error(err)
	//}
}

/*
	模拟客户端
*/
func main() {
	client := client.NewClient("test")

	//配置路由
	client.AddRouter(-1, &PingRouter_{})
	//配置路由
	client.Start()
}
