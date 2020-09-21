package znet

import (
	"errors"
	"github.com/robfig/cron"
	"github.com/wangshiyu/zinx/components"
	"github.com/wangshiyu/zinx/ziface"
)

/*
	组件管理模块
*/
type ComponentManager struct {
	TcpServer ziface.IServer
	//组件数据
	ConnectionDataMap map[string]interface{}
	//当前Server对应的组件
	componentMap map[string]ziface.IComponent
}

/*
	创建一个组件管理
*/
func NewComponentManager(TcpServer ziface.IServer) *ComponentManager {
	ComponentManager := &ComponentManager{
		componentMap: make(map[string]ziface.IComponent),
		TcpServer:    TcpServer,
	}
	Heartbeat := components.NewHeartbeat(TcpServer)
	Heartbeat.Init()
	ComponentManager.Add("Heartbeat", Heartbeat)
	return ComponentManager
}

func (componentMgr *ComponentManager) Runs() {
	maps := componentMgr.Gets()
	if len(maps) > 0 {
		crontab := cron.New()
		for _, Comp := range maps {
			Comp.Init()
			task := func() {
				Comp.Run()
			}
			for _, Cron := range Comp.GetCrons() {
				crontab.AddFunc(Cron, task)
			}
		}
		crontab.Start()
	}
	select {}
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
