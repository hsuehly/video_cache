package yun139

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"video_cache/pkg/feishu"
	"video_cache/pkg/logger"

	"video_cache/pkg/cron"

	"github.com/goccy/go-json"
)

type Yun139 struct {
	authorization string
	account       string
	client        *http.Client
	cron          *cron.Cron
}
type Json map[string]any

var (
	NoFileErr   = errors.New("文件已被删除")
	NoUPloadErr = errors.New("文件已被删除")
)

func New139Yun(authorization string) (*Yun139, error) {

	decode, err := base64.StdEncoding.DecodeString(authorization)
	if err != nil {
		return nil, err
	}
	decodeStr := string(decode)
	splits := strings.Split(decodeStr, ":")
	if len(splits) < 2 {
		return nil, errors.New("authorization is invalid, splits < 2")
	}
	//cutups := strings.Split(splits[2], "|")
	//if len(cutups) < 3 {
	//	return nil, errors.New("expiretime is invalid, splits < 3")
	//}

	//currentTimestamp := time.Now().UnixMilli()
	//futureTimestamp, err := strconv.Atoi(cutups[3])
	if err != nil {
		return nil, err
	}
	//diffTime := int((futureTimestamp-int(currentTimestamp))/1000) - 172800
	yun139 := &Yun139{
		account:       splits[1],
		authorization: authorization,
		client:        &http.Client{},
	}
	yun139.cron = cron.NewCron(time.Second * time.Duration(86400))
	yun139.cron.Do(func() {
		expTime, err := yun139.refreshToken()
		if err != nil {
			logger.Warn("refreshToken err", logger.Err(err))
		}
		fmt.Println("refreshToken load", expTime)
		//yun139.cron.UpdateDuration(time.Duration(expTime) * time.Second)

	})
	err = yun139.GetUserInfo()
	if err != nil {
		return nil, err
	}
	return yun139, err
}
func (d *Yun139) request(pathname string, method string, data any) ([]byte, error) {
	url := "https://yun.139.com" + pathname
	randStr := String(16)
	ts := time.Now().Format("2006-01-02 15:04:05")
	bodyData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	sign := calSign(string(bodyData), ts, randStr)
	body := bytes.NewReader(bodyData)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("CMS-DEVICE", "default")
	req.Header.Add("Authorization", "Basic "+d.authorization)
	req.Header.Add("mcloud-channel", "1000101")
	req.Header.Add("mcloud-client", "10701")
	req.Header.Add("mcloud-sign", fmt.Sprintf("%s,%s,%s", ts, randStr, sign))
	req.Header.Add("mcloud-version", "6.6.0")
	req.Header.Add("Origin", "https://yun.139.com")
	req.Header.Add("Referer", "https://yun.139.com/w/")
	req.Header.Add("x-DeviceInfo", "||9|6.6.0|chrome|95.0.4638.69|uwIy75obnsRPIwlJSd7D9GhUvFwG96ce||macos 10.15.2||zh-CN|||")
	req.Header.Add("x-huawei-channelSrc", "10000034")
	req.Header.Add("x-inner-ntwk", "2")
	req.Header.Add("x-m4c-caller", "PC")
	req.Header.Add("x-m4c-src", "10002")
	req.Header.Add("x-SvcType", "1")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")

	//req.Header.Add("Accept", "application/json, text/plain, */*")
	////req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	//req.Header.Add("Authorization", "Basic "+d.Authorization)
	//req.Header.Add("CMS-DEVICE", "default")
	////req.Header.Add("Connection", "keep-alive")
	////req.Header.Add("INNER-HCY-ROUTER-HTTPS", "1")
	//req.Header.Add("Origin", "https://yun.139.com")
	//req.Header.Add("Referer", "https://yun.139.com/w/")
	////req.Header.Add("Sec-Fetch-Dest", "empty")
	////req.Header.Add("Sec-Fetch-Mode", "cors")
	////req.Header.Add("Sec-Fetch-Site", "same-origin")
	////req.Header.Add("caller", "web")
	//req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	////req.Header.Add("caller", "web")
	//req.Header.Add("mcloud-channel", "1000101")
	//req.Header.Add("mcloud-client", "10701")
	////req.Header.Add("mcloud-route", "001")
	//req.Header.Add("mcloud-sign", fmt.Sprintf("%s,%s,%s", ts, randStr, sign))
	////req.Header.Add("mcloud-skey", "IyliZ0EQ4FjerN2VWbUQm4oFZAvOszJC4PDOcsQIp5OsjAut/AmH0r7H2mKmylJAMa6ombNDN5lot47GZbzGrEiVQUJfbZLogtsVBoo7tPdMDi3lqQVhLrB11VP4300HdDGcYBCDl/Sn2fbZMu2UsiDGtNEFnxDay+xNZeqGOzI=")
	//req.Header.Add("mcloud-version", "7.10.0")
	////req.Header.Add("sec-ch-ua", "Not.A/Brand;v=8, Chromium;v=114, Google Chrome;v=114")
	////req.Header.Add("sec-ch-ua-mobile", "?0")
	////req.Header.Add("sec-ch-ua-platform", "macOS")
	//req.Header.Add("x-DeviceInfo", "||9|7.10.0|chrome|114.0.0.0|c74c333b81e50027dc2f253c887ff7c8||macos 10.15.7||zh-CN|||")
	//req.Header.Add("x-SvcType", "1")
	//req.Header.Add("x-huawei-channelSrc", "10000034")
	//req.Header.Add("x-inner-ntwk", "2")
	//req.Header.Add("x-m4c-caller", "PC")
	//req.Header.Add("x-m4c-src", "10002")
	//req.Header.Add("x-yun-api-version", "v1")
	//req.Header.Add("x-yun-app-channel", "10000034")
	//req.Header.Add("x-yun-channel-source", "10000034")
	//req.Header.Add("x-yun-client-info", "||9|7.10.0|chrome|114.0.0.0|c74c333b81e50027dc2f253c887ff7c8||macos 10.15.7||zh-CN|||dW5kZWZpbmVk||")
	////req.Header.Add("x-yun-module-type", "100")
	//req.Header.Add("x-yun-svc-type", "1")
	//req.Header.Add("Content-Type", "application/json")
	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return respData, nil
}

func (d *Yun139) GetUserInfo() error {
	data := Json{
		"qryUserExternInfoReq": Json{
			"commonAccountInfo": Json{
				"account":     d.account,
				"accountType": 1,
			},
		},
	}
	respData, err := d.request("/orchestration/personalCloud/user/v1.0/qryUserExternInfo", http.MethodPost, data)
	if err != nil {
		return err
	}
	fmt.Println(string(respData), "139data")
	return nil
}

func (d *Yun139) newRequest(pathname string, method string, data any) ([]byte, error) {

	url := "https://personal-kd-njs.yun.139.com" + pathname
	randStr := String(16)
	ts := time.Now().Format("2006-01-02 15:04:05")

	bodyData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	sign := calSign(string(bodyData), ts, randStr)
	body := bytes.NewReader(bodyData)
	req, err := http.NewRequest(method, url, body)
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Authorization", "Basic "+d.authorization)
	req.Header.Add("Caller", "web")
	req.Header.Add("Cms-Device", "default")
	req.Header.Add("Mcloud-Channel", "1000101")
	req.Header.Add("Mcloud-Client", "10701")
	req.Header.Add("Mcloud-Route", "001")
	req.Header.Add("Mcloud-Sign", fmt.Sprintf("%s,%s,%s", ts, randStr, sign))
	req.Header.Add("Mcloud-Version", "7.13.0")
	req.Header.Add("Origin", "https://yun.139.com")
	req.Header.Add("Referer", "https://yun.139.com/w/")
	req.Header.Add("x-DeviceInfo", "||9|7.13.0|chrome|120.0.0.0|||windows 10||zh-CN|||")
	req.Header.Add("x-huawei-channelSrc", "10000034")
	req.Header.Add("x-inner-ntwk", "2")
	req.Header.Add("x-m4c-caller", "PC")
	req.Header.Add("x-m4c-src", "10002")
	req.Header.Add("x-SvcType", "1")
	req.Header.Add("X-Yun-Api-Version", "v1")
	req.Header.Add("X-Yun-App-Channel", "10000034")
	req.Header.Add("X-Yun-Channel-Source", "10000034")
	req.Header.Add("X-Yun-Client-Info", "||9|7.13.0|chrome|120.0.0.0|||windows 10||zh-CN|||dW5kZWZpbmVk||")
	req.Header.Add("X-Yun-Module-Type", "100")
	req.Header.Add("X-Yun-Svc-Type", "1")
	req.Header.Add("Content-Type", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return respData, nil
}

func (d *Yun139) GetLink(contentId string) (string, error) {
	data := Json{
		"operation": "0",
		"inline":    "0",
		"contentID": contentId,
		"extInfo": Json{
			"isReturnCdnDownloadUrl": "1",
		},
		"commonAccountInfo": Json{
			"account":     d.account,
			"accountType": 1,
		},
	}
	resp, err := d.request("/orchestration/personalCloud/uploadAndDownload/v1.0/downloadRequest",
		http.MethodPost, data)
	if err != nil {
		return "", err
	}
	var down DownResult
	err = json.Unmarshal(resp, &down)
	if err != nil {
		return "", err
	}
	if down.Success {
		return down.Data.DownloadURL, nil
	}
	if down.Code == "1010399999" && down.Data.Result.ResultCode == "9149" {
		return "", NoFileErr
	}
	return "", errors.New("获取链接出错:" + down.Message + down.Data.Result.ResultCode)

}

func (d *Yun139) GetNewLink(contentId string) (string, error) {
	data := Json{
		"fileId": contentId,
	}
	resp, err := d.newRequest("/hcy/file/getDownloadUrl", http.MethodPost, data)
	if err != nil {
		return "", err
	}

	return string(resp), nil
}

func (d *Yun139) GetD(contentId string) (string, error) {

	data := Json{
		"getFlvOnlineAddrReq": Json{
			"commonAccountInfo": Json{
				"account":     d.account,
				"accountType": "1",
			},
			"contentID": contentId,
		},
	}
	resp, err := d.request("/orchestration/personalCloud/content/v1.2/getFlvOnlineAddr", http.MethodPost, data)
	if err != nil {
		return "", err
	}

	return string(resp), nil

}

func (d *Yun139) Put(ParentId, fileName string, fileData []byte) (string, error) {
	data := Json{
		"fileCount":       1,
		"manualRename":    2,
		"newCatalogName":  "",
		"operation":       0,
		"parentCatalogID": ParentId,
		"totalSize":       0,
		"commonAccountInfo": Json{
			"account":     d.account,
			"accountType": 1,
		},
		"uploadContentList": []Json{{
			"contentName": fileName,
			"contentSize": 0, // 去除上传大小限制
		}},
	}
	respData, err := d.request("/orchestration/personalCloud/uploadAndDownload/v1.0/pcUploadFileRequest", http.MethodPost, data)
	if err != nil {
		return "", err
	}
	var uploadData UploadResp
	err = json.Unmarshal(respData, &uploadData)
	if err != nil {
		return "", err
	}
	if uploadData.Data.UploadResult.RedirectionURL != "" {
		byteSize := int64(len(fileData))
		req, err := http.NewRequest(http.MethodPost, uploadData.Data.UploadResult.RedirectionURL, bytes.NewReader(fileData))
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "text/plain;name="+unicode(fileName))
		req.Header.Set("contentSize", strconv.FormatInt(byteSize, 10))
		req.Header.Set("range", fmt.Sprintf("bytes=%d-%d", 0, byteSize-1))
		req.Header.Set("Uploadtaskid", uploadData.Data.UploadResult.UploadTaskID)
		req.Header.Set("Rangetype", "0")
		req.Header.Set("Origin", "https://yun.139.com")
		req.Header.Set("Referer", "https://yun.139.com/")
		req.ContentLength = byteSize
		resp, err := d.client.Do(req)
		if err != nil {
			return "", err
		}

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("上传失败: %d", resp.StatusCode)
		}
		return uploadData.Data.UploadResult.NewContentIDList[0].ContentID, nil
	}
	return "", NoUPloadErr
}
func (d *Yun139) Remove(contentId string) error {

	var contentInfoList []string
	var catalogInfoList []string
	contentInfoList = append(contentInfoList, contentId)
	data := Json{
		"createBatchOprTaskReq": Json{
			"taskType":   2,
			"actionType": 202,
			"taskInfo": Json{
				"newCatalogID":    "",
				"contentInfoList": contentInfoList,
				"catalogInfoList": catalogInfoList,
			},
			"commonAccountInfo": Json{
				"account":     d.account,
				"accountType": 1,
			},
		},
	}
	_, err := d.request("/orchestration/personalCloud/batchOprTask/v1.0/createBatchOprTask", http.MethodPost, data)
	if err != nil {
		return err
	}
	return nil
}

func (d *Yun139) refreshToken() (int, error) {
	url := "https://aas.caiyun.feixin.10086.cn:443/tellin/authTokenRefresh.do"
	var resp RefreshTokenResp
	decode, err := base64.StdEncoding.DecodeString(d.authorization)
	if err != nil {
		return 0, err
	}
	decodeStr := string(decode)
	splits := strings.Split(decodeStr, ":")
	reqBody := strings.NewReader("<root><token>" + splits[2] + "</token><account>" + splits[1] + "</account><clienttype><![CDATA[414]]></clienttype></root>")
	req, err := http.NewRequest(http.MethodPost, url, reqBody)
	if err != nil {
		return 0, err
	}
	req.Header.Add("x-NationCode", "+86")
	req.Header.Add("x-NetType", "1")
	req.Header.Add("x-yun-app-channel", "10000023")
	req.Header.Add("x-huawei-channelSrc", "10000023")
	req.Header.Add("x-MM-Source", "032")
	req.Header.Add("x-SvcType", "1")
	req.Header.Add("x-ExpRoute-Code", "routeCode="+splits[1]+",type=10")
	req.Header.Add("Host", "aas.caiyun.feixin.10086.cn")
	req.Header.Add("User-Agent", "okhttp/3.11.0")
	req.Header.Add("Content-Type", "application/xml; charset=UTF-8")
	res, err := d.client.Do(req)
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	err = xml.Unmarshal(body, &resp)
	if err != nil {
		return 0, err
	}
	if resp.Return != "0" {
		_ = feishu.SendMsg("定时刷新token失败" + resp.Desc)
		return 0, fmt.Errorf("failed to refresh token: %s", resp.Desc)
	}
	d.authorization = base64.StdEncoding.EncodeToString([]byte(splits[0] + ":" + splits[1] + ":" + resp.Token))
	_ = feishu.SendMsg("定时刷新token 成功，token " + d.authorization)
	return 86400, nil
}

func (d *Yun139) Destroy() {
	d.cron.Stop()
}
