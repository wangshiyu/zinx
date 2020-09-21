package ziface

type IComponentManager interface {
	Runs()
	Add(key string, component IComponent) //添加组件
	Get(key string) (IComponent, error)   //获取组件
	Gets() map[string]IComponent        //获取全部组件
}
