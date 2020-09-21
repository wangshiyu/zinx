package ziface

type IComponent interface {
	//初始化
	Init()
	//运行
	Run()
	//获取crons 表达式
	GetCrons() []string
}
