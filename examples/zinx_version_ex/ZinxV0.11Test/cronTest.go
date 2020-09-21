package main

import (
	"fmt"
	"github.com/robfig/cron"
	"log"
)

func main() {
	i := 0
	c := cron.New()
	// 也可以秒级任务
	//c := cron.New(cron.WithSeconds()) // 创建定时任务 秒
	//   spec := "0 */1 * * * *" // 每一分钟 如果用到分钟, 秒要设置为0
	//spec := "* */1 * * * *"
	fmt.Println(c)
	spec := "*/1 * * * * *" // 每一分钟，
	c.AddFunc(spec, func() {
		i++
		log.Println("cron running:", i)
		fmt.Println("cron running:", i)
	})
	c.Start()
	select {}
}
