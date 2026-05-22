package logger

import "testing"

func TestLogger(t *testing.T) {
	_, _ = Init(WithFormat("json"), WithSave(true), WithLevel(levelDebug))
	Warn("warn", String("name", "hsuehly"))
	Error("err", Int("num", 1))
}
