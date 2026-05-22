package cron

import (
	"fmt"
	"testing"
	"time"
)

func TestCron(t *testing.T) {
	c := NewCron(time.Second)
	c.Do(func() {
		t.Logf("cron log")
	})
	//time.Sleep(time.Second * 5)
	c.Stop()
}

func TestUPData(t *testing.T) {
	cron := NewCron(4 * time.Second)

	// 启动定时任务
	cron.Do(func() {
		fmt.Println("定时任务执行：", time.Now())
	})

	// 模拟一段时间后更新时间间隔
	time.Sleep(10 * time.Second)
	cron.UpdateDuration(1 * time.Second)

	//// 在适当的时候停止 Cron
	time.Sleep(10 * time.Second)
	cron.Stop()

	fmt.Println("Cron 定时任务已停止")
	time.Sleep(100 * time.Second)
	//select {}
}
