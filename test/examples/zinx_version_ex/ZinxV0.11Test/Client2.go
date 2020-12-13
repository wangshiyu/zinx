package main

import (
	"fmt"
	"github.com/wangshiyu/zinx/ziface"
	ziface2 "github.com/wangshiyu/zinx/ziface/client"
	"github.com/wangshiyu/zinx/znet"
	"github.com/wangshiyu/zinx/znet/client"
	"time"
)

type HelloZinxRouter_ struct {
	znet.BaseRouter
}

//HelloZinxRouter Handle
func (this *HelloZinxRouter_) Handle(request ziface.IRequest) {
	fmt.Println(string(request.GetData()))
}


//创建连接的时候执行
func ClientDoConnectionBegin(conn ziface2.IConnection) {
	//zlog.Debug("DoConnecionBegin is Called ... ")
	//
	////设置两个链接属性，在连接创建之后
	//zlog.Debug("Set conn Name, Home done!")
	//conn.SetProperty("Name", "Aceld")
	//err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	//if err != nil {
	//	zlog.Error(err)
	//}
}

//连接断开的时候执行
func ClientDoConnectionLost(conn ziface2.IConnection) {
	////在连接销毁之前，查询conn的Name，Home属性
	//if name, err := conn.GetProperty("Name"); err == nil {
	//	zlog.Error("Conn Property Name = ", name)
	//}
	//
	//if home, err := conn.GetProperty("Home"); err == nil {
	//	zlog.Error("Conn Property Home = ", home)
	//}
	//
	//zlog.Debug("DoConneciotnLost is Called ... ")
}

/*
	模拟客户端
*/
func main() {

	for n := 0; n <= 3000; n++{
		go func (){
			client := client.NewClient("test")

			//注册链接hook回调函数
			client.SetOnConnStart(ClientDoConnectionBegin)
			client.SetOnConnStop(ClientDoConnectionLost)

			//配置路由
			client.AddRouter(2, &HelloZinxRouter_{})
			//配置路由
			client.Start()

		}()
		time.Sleep(10)
	}
	select {}

}
