package m3task

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/grafov/m3u8"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
	"video_cache/pkg/httpclient"
)

type Parse struct {
	wg            *sync.WaitGroup
	threadLimiter chan struct{}
	client        *httpclient.HttpClient
	ctx           context.Context
}

func NewParse() *Parse {
	return &Parse{
		wg:            &sync.WaitGroup{},
		threadLimiter: make(chan struct{}, 45),
		client:        httpclient.NewHttpClient(),
		ctx:           context.Background(),
	}

}
func (p *Parse) ZA(data bytes.Buffer, baseUrl string) ([]byte, error) {
	playlist, listType, err := m3u8.Decode(data, false)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancelCause(p.ctx)
	defer func() {
		cancel(nil)
	}()
	if listType == m3u8.MEDIA {
		mediaPlaylist := playlist.(*m3u8.MediaPlaylist)
		for _, segment := range mediaPlaylist.GetAllSegments() {
			p.wg.Add(1)
			p.threadLimiter <- struct{}{}
			go func(ctx context.Context, segment *m3u8.MediaSegment, client *httpclient.HttpClient) {
				defer func() {
					p.wg.Done()
					<-p.threadLimiter
				}()
				select {
				case <-ctx.Done():
					return
				default:
					if strings.HasPrefix(segment.URI, "http") {
						if strings.Contains(segment.URI, "/hls/cts_") || strings.Contains(segment.URI, "/hls/m3u8_") {
							for i := 0; i <= 3; i++ {
								if i == 3 {
									cancel(fmt.Errorf("重试最大 链接：%s", segment.URI))
									return
								}
								resp, err := client.Do(http.MethodGet, segment.URI, nil, false, nil)
								if err != nil {
									errorMsg := err.Error()
									//// 查找错误信息中的关键字
									keyword := `failed to parse Location header`
									locationIndex := strings.Index(errorMsg, keyword)
									if locationIndex != -1 {
										// 定位到关键字的末尾位置
										redirectIndex := locationIndex + len(keyword) + 1
										// 截取 URL 地址
										redirectErr := errorMsg[redirectIndex:]
										urlIndex := strings.Index(redirectErr, "https://")
										colonIndex := strings.Index(redirectErr[urlIndex:], `"`)
										if colonIndex != -1 {
											redirectURL := redirectErr[urlIndex : urlIndex+colonIndex]
											segment.URI = redirectURL
											return
										}
									}
									time.Sleep(300 * time.Millisecond)
									continue
								}
								_, _ = io.Copy(io.Discard, resp.Body)
								if resp != nil {
									resp.Body.Close()
								}
								switch resp.StatusCode {
								case 200:
									cancel(fmt.Errorf("获取链接200 链接：%s", segment.URI))
									return
								case 302, 307, 301, 308:
									redirectUrl := resp.Header.Get("Location")
									if redirectUrl == "" {
										cancel(fmt.Errorf("没有获取到链接 链接：%s", segment.URI))
										return
									}
									segment.URI = redirectUrl
									return
								case 503:
									time.Sleep(300 * time.Millisecond)
									continue
								default:
									time.Sleep(300 * time.Millisecond)
									continue
								}
							}
						}

					} else {
						if segment.URI != "" {
							segment.URI = baseUrl + "/" + segment.URI
						}
					}

				}

			}(ctx, segment, p.client)
		}
	} else {
		return nil, errors.New("m3u8列表解析错误")
	}
	p.wg.Wait()
	select {
	case <-ctx.Done():
		return nil, context.Cause(ctx)
	default:
		return playlist.Encode().Bytes(), nil
	}

}
func (p *Parse) AiKu(data bytes.Buffer, baseUrl string) ([]byte, error) {
	playlist, listType, err := m3u8.Decode(data, false)
	if err != nil {
		return nil, err
	}
	if listType == m3u8.MEDIA {
		mediaPlaylist := playlist.(*m3u8.MediaPlaylist)
		for _, segment := range mediaPlaylist.GetAllSegments() {
			if strings.HasPrefix(segment.URI, "http") {
				if strings.Contains(segment.URI, "cdn.json.icu") {
				loop:
					for i := 0; i <= 3; i++ {
						if i == 3 {
							return nil, fmt.Errorf("重试最大 链接：%s", segment.URI)
						}
						resp, err := p.client.Do(http.MethodGet, segment.URI, nil, false, nil)
						if err != nil {
							time.Sleep(300 * time.Millisecond)
							continue
						}
						_, _ = io.Copy(io.Discard, resp.Body)
						if resp != nil {
							resp.Body.Close()
						}
						switch resp.StatusCode {
						case 200:
							return nil, fmt.Errorf("获取链接200 链接：%s", segment.URI)
						case 301, 308, 302, 307:
							redirectUrl := resp.Header.Get("Location")
							if redirectUrl == "" {
								return nil, fmt.Errorf("无重定向地址 链接：%s", segment.URI)
							}
							segment.URI = redirectUrl
							break loop
						case 503:
							time.Sleep(180 * time.Millisecond)
							continue
						default:
							time.Sleep(180 * time.Millisecond)
							continue
						}

					}
				}
			} else {
				if segment.URI != "" {
					segment.URI = baseUrl + "/" + segment.URI
				}
			}
		}
		return playlist.Encode().Bytes(), nil
	} else {
		return nil, errors.New("m3u8列表解析错误")
	}
}

func (p *Parse) M3(data bytes.Buffer, baseUrl string) ([]byte, error) {
	playlist, listType, err := m3u8.Decode(data, false)
	if err != nil {
		return nil, err
	}
	if listType == m3u8.MEDIA {
		mediaPlaylist := playlist.(*m3u8.MediaPlaylist)
		for _, segment := range mediaPlaylist.GetAllSegments() {
			if !strings.HasPrefix(segment.URI, "http") {
				if segment.URI != "" {
					segment.URI = baseUrl + "/" + segment.URI
				}
			}
		}
		return playlist.Encode().Bytes(), nil
	} else {
		return nil, errors.New("m3u8列表解析错误")
	}
}
