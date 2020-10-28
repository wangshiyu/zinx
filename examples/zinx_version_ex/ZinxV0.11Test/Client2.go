package main

import "github.com/wangshiyu/zinx/znet/client"

/*
	模拟客户端
*/
func main() {
	client := client.NewClient("test")
	client.Start()
}
