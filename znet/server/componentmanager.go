package server

import (
	"errors"
	"github.com/robfig/cron"
	server2 "github.com/wangshiyu/zinx/components/server"
	"github.com/wangshiyu/zinx/ziface"
	"github.com/wangshiyu/zinx/ziface/server"
)

/*
	组件管理模块
*/
type ComponentManager struct {
	TcpServer server.IServer
	//组件数据
	ConnectionDataMap map[string]interface{}
	//当前Server对应的组件
	componentMap map[string]ziface.IComponent
}

/*
	创建一个组件管理
*/
func NewComponentManager(TcpServer server.IServer) *ComponentManager {
	ComponentManager := &ComponentManager{
		componentMap: make(map[string]ziface.IComponent),
		TcpServer:    TcpServer,
	}
	Heartbeat := server2.NewHeartbeat(TcpServer)
	ComponentManager.Add("Heartbeat", Heartbeat)
	CheckHeartbeat := server2.NewCheckHeartbeat(TcpServer)
	ComponentManager.Add("CheckHeartbeat", CheckHeartbeat)
	return ComponentManager
}

func (componentMgr *ComponentManager) Runs() {
	maps := componentMgr.Gets()
	if len(maps) > 0 {
		crontab := cron.New()
		for _, Comp := range maps {
			Comp.Init()
			for _, Cron := range Comp.GetCrons() {
				crontab.AddJob(Cron, Comp)
			}
		}
		crontab.Start()
	}
}

//添加组件
func (componentMgr *ComponentManager) Add(key string, conn ziface.IComponent) {
	//将conn连接添加到ConnMananger中
	componentMgr.componentMap[key] = conn
}

//利用ConnID获取链接
func (componentMgr *ComponentManager) Get(key string) (ziface.IComponent, error) {
	if conn, ok := componentMgr.componentMap[key]; ok {
		return conn, nil
	} else {
		return nil, errors.New("component not found")
	}
}

//获取全部链接
func (componentMgr *ComponentManager) Gets() map[string]ziface.IComponent {
	return componentMgr.componentMap
}
