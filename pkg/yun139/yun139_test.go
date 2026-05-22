package yun139

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestYun129(t *testing.T) {
	// cGM6MTc2Mjk4ODY2MDU6VWtuODU2ZHZ8MXxSQ1N8MTcxNTkxNzgxMDIwOHxEOHJXTk5aTy5SMDl2NXJfNkxUdlVXdkJWU1lSa0xvMmVZM2U2VTNMRmFtNS5GYVB5bFI0WnpPeXR6Lk14R2pDMFZ6eEZ5UFJlTWFScFh4MVVCc2R6UUx1eHhMbnRYT0s2LlhQQjd4S0oxeUY5WmpZZDdLWGxhQ0hCOTFIbkxFdWZmajRmZFl1bWoxWGEzWF9TZTY4NkFtRGtUQ2h3Z1NSMzU5TGc0VDZXeU0t
	au := "cGM6MTc2Mjk4ODY2MDU6VWtuODU2ZHZ8MXxSQ1N8MTcxNTkxNzgxMDIwOHxEOHJXTk5aTy5SMDl2NXJfNkxUdlVXdkJWU1lSa0xvMmVZM2U2VTNMRmFtNS5GYVB5bFI0WnpPeXR6Lk14R2pDMFZ6eEZ5UFJlTWFScFh4MVVCc2R6UUx1eHhMbnRYT0s2LlhQQjd4S0oxeUY5WmpZZDdLWGxhQ0hCOTFIbkxFdWZmajRmZFl1bWoxWGEzWF9TZTY4NkFtRGtUQ2h3Z1NSMzU5TGc0VDZXeU0t"
	//yun, err := New139Yun(au)
	//if err != nil {
	//	t.Log(err, "ee")
	//
	//}
	////
	//t.Log(yun.expireTime)
	//t.Log(yun.account)
	//yun.refreshToken()
	//url, err := yun.GetLink("1L11qfFSrFNi30020240324134531uqk")
	//if err != nil {
	//	t.Log(err, "DErr")
	//}
	////
	//t.Log(url, "url")
	//data, err := yun.GetNewLink("1L11qfFSrFNi302202403241329444qq")
	//if err != nil {
	//	t.Log(err, "DErr")
	//}
	//data2, err := yun.GetD("1L11qfFSrFNi302202403241329444qq")
	//if err != nil {
	//	t.Log(err, "eee")
	//}
	//t.Log(data2, "data2")
	//fileData, err := os.ReadFile("./bbb.m3u8")
	//if err != nil {
	//	t.Log("fileErr", err)
	//}
	//conid, err := yun.Put("1L11qfFSrFNi300202404171701336l2", "bbb.m3u8", fileData)
	//if err != nil {
	//	t.Log(err, "e344")
	//}
	//t.Log(conid, "conid")

	//err = yun.Remove("1L11qfFSrFNi29720240417172425mtg")
	//if err != nil {
	//	t.Log(err, "eeee")
	//}
	decode, err := base64.StdEncoding.DecodeString(au)
	if err != nil {
		t.Log(err, "ee")
	}
	decodeStr := string(decode)
	t.Log(decodeStr, "decodeStr")
	splits := strings.Split(decodeStr, ":")
	t.Log(splits, "splits")
	if len(splits) < 2 {
		t.Log("authorization is invalid, splits < 2")
	}
	t.Log(splits[1], " splits[1]")
	t.Log(splits[2], " splits[2]")
	spt := strings.Split(splits[2], "|")
	t.Log(spt[1], " spt[1]")
	t.Log(spt[2], " spt[2]")
	t.Log(spt[3], " spt[3]")
	t.Log(spt[4], " spt[4]")
	u := spt[3]
	fmt.Println(u)

	//if len(spt) < 3 {
	//	t.Log("authorization is invalid, splits < 3")
	//}
}

func TestStr(t *testing.T) {
	str := String(16)
	str1 := String(16)
	str2 := String(16)
	str3 := String(16)
	t.Log(str)
	t.Log(str1)
	t.Log(str2)
	t.Log(str3)
}

func TestTime(t *testing.T) {
	//currentTimestamp := int64(1713347533750) // 示例当前时间戳
	//futureTimestamp := int64(1715939214353)  // 示例未来时间戳

	//// 将时间戳转换为 time.Time 类型
	//currentTime := time.Unix(currentTimestamp, 0)
	//futureTime := time.Unix(futureTimestamp, 0)
	//
	//// 计算时间差
	//difference := futureTime.Sub(currentTime)
	//
	//// 输出结果
	//fmt.Println("时间差为:", difference)
	//f := 123.4560999
	//
	//// 直接将浮点数转换为 int 类型，这将截断小数点后的所有位数
	//i := int(f)
	//
	//fmt.Println("原始浮点数:", f)
	//fmt.Println("取整后的整数:", i)
	ctime := time.Now().UnixMilli()
	wtime := 1715939214353

	diffTime := int((wtime - int(ctime) - 10000) / 1000)
	diffTime2 := int((wtime-int(ctime))/1000) - 86400
	t.Log(diffTime)
	t.Log(diffTime2)
}

func TestTicker(t *testing.T) {
	au := "cGM6MTc2Mjk4ODY2MDU6VWtuODU2ZHZ8MXxSQ1N8MTcxNTkxNzgxMDIwOHxEOHJXTk5aTy5SMDl2NXJfNkxUdlVXdkJWU1lSa0xvMmVZM2U2VTNMRmFtNS5GYVB5bFI0WnpPeXR6Lk14R2pDMFZ6eEZ5UFJlTWFScFh4MVVCc2R6UUx1eHhMbnRYT0s2LlhQQjd4S0oxeUY5WmpZZDdLWGxhQ0hCOTFIbkxFdWZmajRmZFl1bWoxWGEzWF9TZTY4NkFtRGtUQ2h3Z1NSMzU5TGc0VDZXeU0t"

	_, err := New139Yun(au)
	if err != nil {
		t.Log(err)
	}
	time.Sleep(time.Second * 50)
}

func TestSwitch(t *testing.T) {
	pan := 3
	switch pan {
	case 3:
		break
	case 6:

	}
	fmt.Println("222")
}
