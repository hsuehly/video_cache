package poller

import (
	"context"
	"testing"
)

func TestPoll(t *testing.T) {
	producer := NewPoller(50)
	producer.Poll(context.Background())

}
