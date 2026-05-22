package aliyundrive

import (
	"context"
	"fmt"
	"github.com/tickstep/aliyunpan-api/aliyunpan"
	"sync"
	"time"

	"github.com/tickstep/aliyunpan-api/aliyunpan_web"
)

type TokenManager interface {
	KeepAlive(ctx context.Context, t time.Duration)
}
type keepAliveTokenManager struct {
	drive                 *aliyunpan_web.WebPanClient
	refreshToken          string
	accessTokenExpireTime time.Time
	lock                  *sync.Mutex
	wg                    *sync.WaitGroup
}

func NewKeepAliveTokenManager(drive *aliyunpan_web.WebPanClient, refreshToken string, accessTokenExpireTime time.Time) TokenManager {
	return &keepAliveTokenManager{
		drive:                 drive,
		refreshToken:          refreshToken,
		accessTokenExpireTime: accessTokenExpireTime,
		wg:                    new(sync.WaitGroup),
		lock:                  new(sync.Mutex),
	}
}

func (m *keepAliveTokenManager) AccessToken() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	fmt.Println("进入AccessToken", m.accessTokenExpireTime)
	now := time.Now()
	if now.Before(m.accessTokenExpireTime) {
		fmt.Println("不需要更新")

		return nil
	}
	err := m.refresh()
	_, err = m.drive.CreateSession(nil)
	res, err := m.drive.GetFileDownloadUrl(&aliyunpan.GetFileDownloadUrlParam{
		FileId: "661e90dc020425e65dbf4cf4b1632d4d55155f52",
	})
	fmt.Println(res.InternalUrl)
	fmt.Println(res.Url)
	fmt.Println(res.CdnUrl)
	if err != nil {
		return err
	}
	return nil
}
func (m *keepAliveTokenManager) refresh() error {
	fmt.Println("进入refresh")
	now := time.Now()
	result, err := aliyunpan_web.GetAccessTokenFromRefreshToken(m.refreshToken)
	if err != nil {
		fmt.Println(err, "err")
		return err
	}
	fmt.Println(result.ExpiresIn, "ExpiresIn")
	m.refreshToken = result.RefreshToken
	//m.accessToken = result.AccessToken
	fmt.Println("time", now.Add(time.Second*time.Duration(result.ExpiresIn-60)))
	m.accessTokenExpireTime = now.Add(time.Second * time.Duration(result.ExpiresIn-60))
	m.drive.UpdateToken(aliyunpan_web.WebLoginToken{
		AccessTokenType: result.AccessTokenType,
		AccessToken:     result.AccessToken,
		RefreshToken:    result.RefreshToken,
		ExpiresIn:       result.ExpiresIn,
		ExpireTime:      result.ExpireTime,
	})
	fmt.Println("更新token完成")

	return nil
}
func (m *keepAliveTokenManager) KeepAlive(ctx context.Context, t time.Duration) {
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		ticker := time.NewTicker(t)
	keepaliveLoop:
		for {
			select {
			case <-ticker.C:
			case <-ctx.Done():
				break keepaliveLoop
			}
			_ = m.AccessToken()
		}
		ticker.Stop()
	}()
}

func (m *keepAliveTokenManager) WaitStop() {
	m.wg.Wait()
}
