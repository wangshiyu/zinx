package client

import (
	"github.com/wangshiyu/zinx/ziface"
	"github.com/wangshiyu/zinx/ziface/client"
)

type Request struct {
	conn client.IConnection //已经和客户端建立好的 链接
	msg  ziface.IMessage    //客户端请求的数据
}

//获取请求连接信息
func (r *Request) GetConnection() interface{} {
	return r.conn
}

//获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

//获取请求的消息的ID
func (r *Request) GetMsgID() int32 {
	return r.msg.GetMsgId()
}
