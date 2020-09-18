/**
* @Author: Aceld(刘丹冰)
* @Date: 2019/5/9 10:14
* @Mail: danbing.at@gmail.com
*
*  时间轮定时器调度器单元测试
 */
package test

import (
	"fmt"
	"github.com/wangshiyu/zinx/zlog"
	"github.com/wangshiyu/zinx/ztimer"
	"testing"
	"time"
)

//触发函数
func foo(args ...interface{}) {
	fmt.Printf("I am No. %d function, delay %d ms\n", args[0].(int), args[1].(int))
}

//手动创建调度器运转时间轮
func TestNewTimerScheduler(t *testing.T) {
	timerScheduler := ztimer.NewTimerScheduler()
	timerScheduler.Start()

	//在scheduler中添加timer
	for i := 1; i < 2000; i++ {
		f := ztimer.NewDelayFunc(foo, []interface{}{i, i * 3})
		tid, err := timerScheduler.CreateTimerAfter(f, time.Duration(3*i)*time.Millisecond)
		if err != nil {
			zlog.Error("create timer error", tid, err)
			break
		}
	}

	//执行调度器触发函数
	go func() {
		delayFuncChan := timerScheduler.GetTriggerChan()
		for df := range delayFuncChan {
			df.Call()
		}
	}()

	//阻塞等待
	select {}
}

//采用自动调度器运转时间轮
func TestNewAutoExecTimerScheduler(t *testing.T) {
	autoTS := ztimer.NewAutoExecTimerScheduler()

	//给调度器添加Timer
	for i := 0; i < 2000; i++ {
		f := ztimer.NewDelayFunc(foo, []interface{}{i, i * 3})
		tid, err := autoTS.CreateTimerAfter(f, time.Duration(3*i)*time.Millisecond)
		if err != nil {
			zlog.Error("create timer error", tid, err)
			break
		}
	}

	//阻塞等待
	select {}
}
