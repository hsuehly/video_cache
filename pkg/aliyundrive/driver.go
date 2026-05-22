package aliyundrive

import (
	"context"
	"fmt"
	"github.com/tickstep/aliyunpan-api/aliyunpan"
	"github.com/tickstep/aliyunpan-api/aliyunpan_web"
	"time"
)

type AliDriver struct {
	d       *aliyunpan_web.WebPanClient
	DriveId string
}

// be2bb357c87e46e498322fd281fa70f1
func NewAliClient(refreshToken string) (*AliDriver, error) {
	// get access token
	//refreshToken := "be2bb357c87e46e498322fd281fa70f1"
	ali := new(AliDriver)
	webToken, err := aliyunpan_web.GetAccessTokenFromRefreshToken(refreshToken)
	if err != nil {
		fmt.Println(err, "err3")
		return nil, err
	}
	//fmt.Println(webToken, "webToken")
	// web pan client
	appConfig := aliyunpan_web.AppConfig{
		AppId:     "25dzX3vbYqktVxyX",
		DeviceId:  "500880293",
		UserId:    "8a12d38b3d154327940b10fd0861a9bd",
		Nonce:     0,
		PublicKey: "",
	}
	appToken := aliyunpan_web.AppLoginToken{
		AccessToken:  "",
		RefreshToken: refreshToken,
	}
	appSession := aliyunpan_web.SessionConfig{
		DeviceName: "Chrome浏览器",
		ModelName:  "Windows网页版",
	}
	webPanClient := aliyunpan_web.NewWebPanClient(*webToken, appToken, appConfig, appSession)

	r, err := webPanClient.CreateSession(&aliyunpan_web.CreateSessionParam{
		DeviceName: "Chrome浏览器",
		ModelName:  "Windows网页版",
	})
	fmt.Printf("%v ----", r)
	if err != nil {
		fmt.Println(err, "err2")
		return nil, err
	}
	resp, err := webPanClient.GetUserInfo()
	if err != nil {
		fmt.Println(err, "err1")
		return nil, err
	}
	fmt.Printf("%v --", resp)

	fmt.Println(resp.FileDriveId)
	fmt.Println(resp.SafeBoxDriveId)
	fmt.Println(resp.AlbumDriveId)
	ali.DriveId = resp.FileDriveId
	ali.d = webPanClient
	res, err := webPanClient.GetFileDownloadUrl(&aliyunpan.GetFileDownloadUrlParam{
		FileId: "661e90dc020425e65dbf4cf4b1632d4d55155f52",
	})
	if err != nil {
		fmt.Println(err, "err")
	}
	fmt.Println(res.InternalUrl)
	fmt.Println(res.Url)
	fmt.Println(res.CdnUrl)
	mange := NewKeepAliveTokenManager(webPanClient, refreshToken, time.Now().Add(time.Second*time.Duration(10)))
	mange.KeepAlive(context.Background(), time.Second*1)
	//return
	return ali, nil
}
