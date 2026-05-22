package httpclient

import (
	"io"
	"net/http"
	"time"
)

type HttpClient struct {
	client *http.Client // http客户端实例
}

func NewHttpClient() *HttpClient {
	return &HttpClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 60,
				IdleConnTimeout:     30 * time.Second,
				TLSHandshakeTimeout: 10 * time.Second,
				//TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

func (c *HttpClient) Do(method, url string, body io.Reader, follow bool, headers map[string]string) (*http.Response, error) {
	// 创建http请求
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.43")
	//req.Header.Set("Connection", "keep-alive")
	// 设置是否跟随重定向
	if !follow {
		c.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // 不跟随重定向
		}
	}
	// 设置请求头
	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	// 发送http请求
	return c.client.Do(req)
}
