package main

import (
	"fmt"
	"github.com/wangshiyu/zinx/ziface"
	"github.com/wangshiyu/zinx/znet"
	"github.com/wangshiyu/zinx/znet/client"
)

type HelloZinxRouter_ struct {
	znet.BaseRouter
}

//HelloZinxRouter Handle
func (this *HelloZinxRouter_) Handle(request ziface.IRequest) {
	fmt.Println(string(request.GetData()))
}

/*
	模拟客户端
*/
func main() {
	client := client.NewClient("test")

	//配置路由
	client.AddRouter(2, &HelloZinxRouter_{})
	//配置路由
	client.Start()
}
