package server

import (
	"fmt"
	"github.com/wangshiyu/zinx/utils"
	"github.com/wangshiyu/zinx/ziface"
	server "github.com/wangshiyu/zinx/ziface/server"
	"github.com/wangshiyu/zinx/zlog"
	"github.com/wangshiyu/zinx/znet"
	"net"
)

var zinxLogo = `                                        
              ██                        
              ▀▀                        
 ████████   ████     ██▄████▄  ▀██  ██▀ 
     ▄█▀      ██     ██▀   ██    ████   
   ▄█▀        ██     ██    ██    ▄██▄   
 ▄██▄▄▄▄▄  ▄▄▄██▄▄▄  ██    ██   ▄█▀▀█▄  
 ▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀  ▀▀    ▀▀  ▀▀▀  ▀▀▀ 
                                        `
var topLine = `┌───────────────────────────────────────────────────┐`
var borderLine = `│`
var bottomLine = `└───────────────────────────────────────────────────┘`

//iServer 接口实现，定义一个Server服务类
type Server struct {
	//服务器的名称
	Name string
	//tcp4 or other
	IPVersion string
	//服务绑定的IP地址
	IP string
	//服务绑定的端口
	Port uint32
	//当前Server的消息管理模块，用来绑定MsgId和对应的处理方法
	msgHandler ziface.IMsgHandle
	//当前Server的链接管理器
	ConnMgr server.IConnManager
	//该Server的连接创建时Hook函数
	OnConnStart func(conn server.IConnection)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn server.IConnection)
	//组件管理
	ComponentManager ziface.IComponentManager
	//加密
	Encryption ziface.IEncryption
}

/*
  创建一个服务器句柄
*/
func NewServer() server.IServer {

	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		msgHandler: znet.NewMsgHandle(),
		ConnMgr:    NewConnManager(),
		Encryption: znet.NewRSA2(),
	}
	s.ComponentManager = NewComponentManager(s)
	return s
}

//============== 实现 ziface.IServer 里的全部接口方法 ========

//开启网络服务
func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)

	go s.ComponentManager.Runs()

	//开启一个go去做服务端Linster业务
	go func() {
		//0 启动worker工作池机制
		s.msgHandler.StartWorkerPool()

		//1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			zlog.Error("resolve tcp addr err: ", err)
			return
		}

		//2 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			zlog.Error("listen", s.IPVersion, "err", err)
			return
		}

		//已经监听成功
		zlog.Info("start Zinx server  ", s.Name, " succ, now listenning...")

		//3 启动server网络连接业务
		for {
			//3.1 阻塞等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				zlog.Error("Accept err ", err)
				continue
			}
			zlog.Info("lank conn addr = ",conn.RemoteAddr().String())

			//3.2 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				zlog.Warn("ConnLen:" + string(s.ConnMgr.Len()) + " MaxConnLen:" + string(utils.GlobalObject.MaxConn))
				conn.Close()
				continue
			}
			//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
			dealConn := NewConntion(s, conn, conn.RemoteAddr().String(), s.msgHandler)
			//cid++
			//3.4 启动当前链接的处理业务
			go dealConn.Start()
		}
	}()
}

//停止服务
func (s *Server) Stop() {
	zlog.Info("[STOP] Zinx server , name ", s.Name)
	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.ConnMgr.ClearConn()
}

//运行服务
func (s *Server) Serve() {
	s.Start()

	//阻塞,否则主Go退出， listenner的go将会退出
	select {}
}

//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRouter(msgId int32, router ziface.IRouter) {
	if msgId < 0 {
		panic("msgId < 0")
	}
	s.msgHandler.AddRouter(msgId, router)
}

//得到链接管理
func (s *Server) GetConnMgr() server.IConnManager {
	return s.ConnMgr
}

//设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(server.IConnection)) {
	s.OnConnStart = hookFunc
}

//设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(server.IConnection)) {
	s.OnConnStop = hookFunc
}

//调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn server.IConnection) {
	if s.OnConnStart != nil {
		zlog.Info("---> CallOnConnStart.....")
		s.OnConnStart(conn)
	}
}

//调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn server.IConnection) {
	if s.OnConnStop != nil {
		zlog.Info("---> CallOnConnStop....")
		s.OnConnStop(conn)
	}
}

//获取组件管理器
func (s *Server) GetComponentMgr() ziface.IComponentManager {
	return s.ComponentManager
}

func (s *Server) GetEncryption() ziface.IEncryption {
	return s.Encryption
}

func init() {
	fmt.Println(zinxLogo)
	fmt.Println(topLine)
	fmt.Println(fmt.Sprintf("%s [Github] https://github.com/aceld                 %s", borderLine, borderLine))
	fmt.Println(fmt.Sprintf("%s [tutorial] https://www.jianshu.com/p/23d07c0a28e5 %s", borderLine, borderLine))
	fmt.Println(bottomLine)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)
}
