package timer

import (
	"fmt"
	"time"
)

// TimedTask 结构体用于管理定时任务
type TimedTask struct {
	timer    *time.Timer // 定时器
	execFun  func()      // 定时执行的任务函数
	execTime time.Duration
}

// NewTimedTask 创建一个新的定时任务，需要传入初始时间间隔和任务执行函数
func NewTimedTask(initialTime time.Duration, execFun func()) *TimedTask {
	return &TimedTask{
		timer:    time.NewTimer(initialTime),
		execFun:  execFun,
		execTime: initialTime,
	}
}

// Start 启动定时任务，无限循环直到 Stop 被调用
func (t *TimedTask) Start() {
	go func() {
		for {
			select {
			case <-t.timer.C:
				t.execFun() // 执行任务函数
				t.timer.Reset(t.execTime)
				fmt.Println("Executing task at", time.Now())
			}
		}
	}()
}

// UpdateTime 更新定时任务的下一次执行时间
func (t *TimedTask) UpdateTime(d time.Duration) {
	t.execTime = d
	t.timer.Reset(t.execTime)
}

// Stop 停止定时任务，释放资源
func (t *TimedTask) Stop() {
	if !t.timer.Stop() {
		// 如果定时器已经触发，则从通道中读取以释放资源
		<-t.timer.C
	}
}
