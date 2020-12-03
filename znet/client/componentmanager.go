package client

import (
	"errors"
	"github.com/robfig/cron"
	client2 "github.com/wangshiyu/zinx/zcomponents/client"
	"github.com/wangshiyu/zinx/ziface"
	"github.com/wangshiyu/zinx/ziface/client"
)

/*
	组件管理模块
*/
type ComponentManager struct {
	Client client.IClient
	//组件数据
	ConnectionDataMap map[string]interface{}
	//当前Server对应的组件
	componentMap map[string]ziface.IComponent
}

/*
	创建一个组件管理
*/
func NewComponentManager(Client client.IClient) *ComponentManager {
	ComponentManager := &ComponentManager{
		componentMap: make(map[string]ziface.IComponent),
		Client:       Client,
	}
	Heartbeat := client2.NewHeartbeat(Client)
	ComponentManager.Add("Heartbeat", Heartbeat)
	CheckHeartbeat := client2.NewCheckHeartbeat(Client)
	ComponentManager.Add("CheckHeartbeat", CheckHeartbeat)
	Reconnection := client2.NewReconnection(Client)
	ComponentManager.Add("Reconnection", Reconnection)
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
