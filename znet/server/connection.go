package server

import (
	"context"
	"errors"
	"github.com/wangshiyu/zinx/utils"
	"github.com/wangshiyu/zinx/ziface"
	"github.com/wangshiyu/zinx/ziface/server"
	"github.com/wangshiyu/zinx/zlog"
	"github.com/wangshiyu/zinx/znet"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

type Connection struct {
	//当前Conn属于哪个Server
	TcpServer server.IServer
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	//当前连接的ID 也可以称作为SessionID，ID全局唯一
	ConnID uint32
	//消息管理MsgId和对应处理方法的消息管理模块
	MsgHandler ziface.IMsgHandle
	//告知该链接已经退出/停止的channel
	ctx    context.Context
	cancel context.CancelFunc
	//无缓冲管道，用于读、写两个goroutine之间的消息通信
	msgChan chan []byte
	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte

	sync.RWMutex
	//链接属性
	property map[string]interface{}
	//当前连接的关闭状态
	isClosed bool
	//组件数据
	ConnectionDataMap map[struct{}]interface{}
}

//创建连接的方法
func NewConntion(server server.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	//初始化Conn属性
	c := &Connection{
		TcpServer:   server,
		Conn:        conn,
		ConnID:      connID,
		isClosed:    false,
		MsgHandler:  msgHandler,
		msgChan:     make(chan []byte),
		msgBuffChan: make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
		property:    make(map[string]interface{}),
	}
	//链接时间
	c.SetProperty(znet.LINK_TIME,time.Now())
	//将新创建的Conn添加到链接管理中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

/*
	写消息Goroutine， 用户将数据发送给客户端
*/
func (c *Connection) StartWriter() {
	zlog.Info("[Writer Goroutine is running]")
	defer zlog.Info(c.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				zlog.Error("Send Data error:, ", err, " Conn Writer exit")
				return
			}
			zlog.Debugf("Send data succ! data = %+v\n", data)
		case data, ok := <-c.msgBuffChan:
			if ok {
				//有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					zlog.Error("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				zlog.Error("msgBuffChan is Closed")
				break
			}
		case <-c.ctx.Done():
			return
		}
	}
}

/*
	读消息Goroutine，用于从客户端中读取数据
*/
func (c *Connection) StartReader() {
	zlog.Info("[Reader Goroutine is running]")
	defer zlog.Info(c.RemoteAddr().String(), "[conn Reader exit!]")
	defer c.Stop()
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			// 创建拆包解包的对象
			dp := znet.NewDataPack()

			//读取客户端的Msg head
			headData := make([]byte, dp.GetHeadLen())
			if _, err := io.ReadFull(c.Conn, headData); err != nil {
				switch {
				case strings.Contains(err.Error(), "connection reset"):
					zlog.Error("Connection refused")
					return
				case strings.Contains(err.Error(), "EOF"):
					zlog.Error("EOF")
					return
				default:
					zlog.Errorf("read msg head error Unknown error:%s", err)
				}
				break
			}
			zlog.Debugf("read headData %+v\n", headData)

			//拆包，得到msgid 和 datalen 放在msg中
			msg, err := dp.Unpack(headData)
			if err != nil {
				//fmt.Println("unpack error ", err)
				zlog.Error("unpack error ", err)
				break
			}

			//根据 dataLen 读取 data，放在msg.Data中
			var data []byte
			if msg.GetDataLen() > 0 {
				data = make([]byte, msg.GetDataLen())
				if _, err := io.ReadFull(c.Conn, data); err != nil {
					zlog.Error("read msg data error ", err)
					break
				}
			}
			if utils.GlobalObject.Encryption {
				data = c.TcpServer.GetEncryption().Decrypt(data)
			}
			msg.SetData(data)
			zlog.Debug("server read data = ", string(data))
			//更新消息接收时间
			c.SetProperty(znet.LAST_MSG_READ_DATE, time.Now())
			c.SetProperty(znet.LAST_MSG_READ_LEN, len(data))
			len_, _ := c.GetProperty(znet.READ_MSG_LEN)
			if len_ == nil {
				len_ = int64(0)
			}
			len_ = len_.(int64) + int64(len(data))
			c.SetProperty(znet.READ_MSG_LEN, len_)
			//得到当前客户端请求的Request数据
			req := Request{
				conn: c,
				msg:  msg,
			}
			if utils.GlobalObject.WorkerPoolSize > 0 {
				//已经启动工作池机制，将消息交给Worker处理
				c.MsgHandler.SendMsgToTaskQueue(&req)
			} else {
				//从绑定好的消息和对应的处理方法中执行对应的Handle方法
				go c.MsgHandler.DoMsgHandler(&req)
			}
		}
	}
}

//启动连接，让当前连接开始工作
func (c *Connection) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	//1 开启用户从客户端读取数据流程的Goroutine
	go c.StartReader()
	//2 开启用于写回客户端数据流程的Goroutine
	go c.StartWriter()
	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.TcpServer.CallOnConnStart(c)
}

//停止连接，结束当前连接状态M
func (c *Connection) Stop() {
	//fmt.Println("Conn Stop()...ConnID = ", c.ConnID)
	zlog.Info("Conn Stop()...ConnID = ", c.ConnID)
	//如果当前链接已经关闭
	c.Lock()
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	c.Unlock()

	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.TcpServer.CallOnConnStop(c)

	c.Lock()
	defer c.Unlock()
	// 关闭socket链接
	c.Conn.Close()
	//关闭Writer
	c.cancel()

	//将链接从连接管理器中删除
	c.TcpServer.GetConnMgr().Remove(c)

	//关闭该链接全部管道
	close(c.msgBuffChan)
}

//从当前连接获取原始的socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(msgId int32, data []byte) error {
	c.RLock()
	if c.isClosed == true {
		c.RUnlock()
		return errors.New("connection closed when send msg")
	}
	c.RUnlock()
	zlog.Debug("Server SendBuffMsg data = ", string(data))
	//将data封包，并且发送
	dp := znet.NewDataPack()
	if utils.GlobalObject.Encryption {
		data = c.TcpServer.GetEncryption().Encryption(data)
	}
	msg, err := dp.Pack(znet.NewMsgPackage(msgId, data))
	if err != nil {
		zlog.Error("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}
	//写回客户端
	c.msgChan <- msg

	return nil
}

func (c *Connection) SendBuffMsg(msgId int32, data []byte) error {
	c.RLock()
	if c.isClosed == true {
		c.RUnlock()
		return errors.New("Connection closed when send buff msg")
	}
	c.RUnlock()
	zlog.Debug("Server SendBuffMsg data = ", string(data))
	//将data封包，并且发送
	dp := znet.NewDataPack()
	if utils.GlobalObject.Encryption {
		data = c.TcpServer.GetEncryption().Encryption(data)
	}
	msg, err := dp.Pack(znet.NewMsgPackage(msgId, data))
	if err != nil {
		zlog.Error("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}
	//写回客户端
	c.msgBuffChan <- msg

	return nil
}

//设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.Lock()
	defer c.Unlock()

	c.property[key] = value
}

//获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.RLock()
	defer c.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

//移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.Lock()
	defer c.Unlock()

	delete(c.property, key)
}

//是否授权
func (c *Connection) IsAuth() bool {
	return false
}

//是否关闭
func (c *Connection) IsClosed() bool {
	return c.isClosed
}
