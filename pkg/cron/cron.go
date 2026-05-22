package cron

import (
	"time"
)

type Cron struct {
	ch chan struct{}
	t  *time.Ticker
}

func NewCron(d time.Duration) *Cron {
	return &Cron{
		ch: make(chan struct{}),
		t:  time.NewTicker(d),
	}
}

func (c *Cron) Do(f func()) {
	go func() {

		defer c.t.Stop()
		for {
			select {
			case <-c.t.C:
				f()
			case <-c.ch:
				return
			}
		}
	}()
}
func (c *Cron) UpdateDuration(d time.Duration) {
	c.t.Reset(d)
}
func (c *Cron) Stop() {
	select {
	case _, _ = <-c.ch:
	default:
		c.ch <- struct{}{}
		close(c.ch)
	}
}
