package httpclient

import (
	"testing"
)

func TestReq(t *testing.T) {
	u := "https://ts.instantts.top/hls/cts_rWjmAQp0Aer0AQNjYALf3ZfiZgRt2mx5AmtLNETRTut_8a18b7b000f1"
	client := NewHttpClient()
	resp, err := client.Do("GET", u, nil, false, nil)
	if err != nil {
		t.Log(err)
	}
	err = resp.Body.Close()
	if err != nil {
		t.Log(err, "body")
	}

	t.Log(resp.Header.Get("Location"))
	t.Log(resp.StatusCode)
}
