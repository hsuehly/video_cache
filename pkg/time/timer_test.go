package timer

import (
	"fmt"
	"testing"
	"time"
)

// Example task function for demonstration purposes
func exampleTaskFunction() {
	fmt.Println("定时任务执行")
}

func TestTime(t *testing.T) {
	// 创建定时任务实例，初始定时时间为3800秒
	// 创建并启动定时任务
	task := NewTimedTask(5*time.Second, exampleTaskFunction)
	task.Start()

	// 在适当的时候停止任务
	time.Sleep(10 * time.Second)
	task.UpdateTime(time.Second * 1)
	time.Sleep(5 * time.Second)
	task.Stop()
	select {}
}
