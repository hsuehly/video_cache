package m3task

import (
	"fmt"
	"strings"
	"testing"
)

func TestM3(t *testing.T) {
	//u := "https://ts.instantts.top/hls/m3u8_71e80a62590a2c2a6920407955483325_e9582d582167"
	u1 := "https://ts.instantts.top/hls/cts_rWjmAQp0Aep0AeLjF3mHkAQOVFmB1YYLNNQCxOFN_bca51862872d"
	fmt.Println(strings.Contains(u1, "/hls/m3u8_"))
	fmt.Println(strings.Contains(u1, "/hls/cts_"))
}

func TestFor(t *testing.T) {
loop:
	for i := 0; i < 17; i++ {

		switch i {
		case 4:
			continue
		case 6:
			break loop

		}
		t.Log(i, "i")

	}
	fmt.Println("oopp")
}
