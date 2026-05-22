package rmq

import (
	"video_cache/pkg/logger"
)

// 死信队列，当消息处理失败达到指定次数时，会被投递到此处
type DeadLetterMailbox interface {
	Deliver(msg *MsgEntity) error
}

// 默认使用的死信队列，仅仅对消息失败的信息进行日志打印
type DeadLetterLogger struct{}

func NewDeadLetterLogger() *DeadLetterLogger {
	return &DeadLetterLogger{}
}

func (d *DeadLetterLogger) Deliver(msg *MsgEntity) error {
	logger.Warn("fail execeed retry limit", logger.String("msg id:", msg.MsgID))
	return nil
}
