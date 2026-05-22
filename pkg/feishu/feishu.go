package feishu

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func SendMsg(txt string) error {
	client := &http.Client{}
	msg := fmt.Sprintf(`{"msg_type":"text","content":{"text":"%s"}}`, txt)
	payload := strings.NewReader(msg)
	req, err := http.NewRequest(http.MethodPost, "https://open.feishu.cn/open-apis/bot/v2/hook/e223646e-0fbd-4cfa-86f7-c6902795ce35", payload)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	err = resp.Body.Close()
	if err != nil {
		return err
	}
	return nil
}
