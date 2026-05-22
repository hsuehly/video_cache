package aliyundrive

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestAliYun(t *testing.T) {
	NewAliClient("be2bb357c87e46e498322fd281fa70f1")

	time.Sleep(30 * time.Second)
}
func TestFor(t *testing.T) {
loop:
	for i := 0; i < 10; i++ {
		if i == 4 {
			fmt.Println("跳过4")
			break loop // 跳过当前迭代的剩余部分
		}
		fmt.Println(i)
	}
}

//func TestReq(t *testing.T) {
//	d := Drive{}
//	data := d.requestWithoutCredit()
//	//fmt.Print(data["access_token"])
//	//fmt.Print(data["refresh_token"])
//	//fmt.Print(data["expires_in"].(int))
//	var i int
//	i = data["expires_in"].(int)
//
//	fmt.Println(i, "ii-i")
//
//}
//
//func TestTokenMange(t *testing.T) {
//	d := &Drive{driveId: "123jsjsj"}
//	RefreshToken := "redkkelaklkdalada"
//	tokenManager := NewKeepAliveTokenManager(NewRefreshTokenManager(d, RefreshToken))
//	tokenManager.KeepAlive(context.Background(), time.Second*10)
//	//tokenManager.WaitStop()
//
//	time.Sleep(time.Second * 30)
//}

// func (d *Drive) requestWithoutCredit() map[string]any {
//
//	return map[string]any{
//		"access_token":  "123",
//		"refresh_token": "456",
//		"expires_in":    70,
//	}
//
// }
func TestTimee(t *testing.T) {
	now := time.Now()
	fmt.Println(now.Add(time.Second * time.Duration(120-60)))
}

func TestTimer(t *testing.T) {
	timer := time.NewTimer(time.Second * 2)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*28)
	ti := 0
	defer func() {
		timer.Stop()
		cancel()
	}()
loop:
	for {
		select {
		case <-timer.C:
			fmt.Println("时间到")
			ti++
			fmt.Println(ti)
			if ti < 5 {
				timer.Reset(time.Second * 2)
				fmt.Println(ti)
			}

		case <-ctx.Done():
			fmt.Println("超时")
			break loop

		}

	}
	fmt.Println("结束")
}
