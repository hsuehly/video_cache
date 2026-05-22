package rmq

import (
	"errors"
	"fmt"
	"testing"
)

func TestRmq(t *testing.T) {

	fmt.Println(!errors.Is(ErrNoMsg, ErrNoMsg))
}
