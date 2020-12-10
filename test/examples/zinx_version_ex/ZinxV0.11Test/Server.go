/**
* @Author: Aceld
* @Date: 2019/4/30 17:42
* @Mail: danbing.at@gmail.com
*  ZinxV0.11测试，测试Zinx 日志模块功能 zlog模块
 */
package main

import (
	"github.com/wangshiyu/zinx/ziface"
	ziface2 "github.com/wangshiyu/zinx/ziface/server"
	"github.com/wangshiyu/zinx/zlog"
	"github.com/wangshiyu/zinx/znet"
	"github.com/wangshiyu/zinx/znet/server"
)


type HelloZinxRouter struct {
	znet.BaseRouter
}

//HelloZinxRouter Handle
func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	zlog.Debug("Call HelloZinxRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	zlog.Debug("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
	err := request.GetConnection().(ziface2.IConnection).SendBuffMsg(1, []byte("Hello Zinx Router V0.10"))
	if err != nil {
		zlog.Error(err)
	}
}

//创建连接的时候执行
func DoConnectionBegin(conn ziface2.IConnection) {
	zlog.Debug("DoConnecionBegin is Called ... ")

	//设置两个链接属性，在连接创建之后
	zlog.Debug("Set conn Name, Home done!")
	conn.SetProperty("Name", "Aceld")
	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		zlog.Error(err)
	}
}

//连接断开的时候执行
func DoConnectionLost(conn ziface2.IConnection) {
	//在连接销毁之前，查询conn的Name，Home属性
	if name, err := conn.GetProperty("Name"); err == nil {
		zlog.Error("Conn Property Name = ", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		zlog.Error("Conn Property Home = ", home)
	}

	zlog.Debug("DoConneciotnLost is Called ... ")
}

func main() {
	//创建一个server句柄
	s := server.NewServer()

	//注册链接hook回调函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//配置路由
	s.AddRouter(2, &HelloZinxRouter{})

	//开启服务
	s.Serve()
}
