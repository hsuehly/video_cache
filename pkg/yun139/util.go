package yun139

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

//var Rand *rand.Rand

const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

var Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func String(n int) string {
	// 初始化 Rand，使用当前时间作为种子

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[Rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func encodeURIComponent(str string) string {
	r := url.QueryEscape(str)
	r = strings.Replace(r, "+", "%20", -1)
	r = strings.Replace(r, "%21", "!", -1)
	r = strings.Replace(r, "%27", "'", -1)
	r = strings.Replace(r, "%28", "(", -1)
	r = strings.Replace(r, "%29", ")", -1)
	r = strings.Replace(r, "%2A", "*", -1)
	return r
}

func calSign(body, ts, randStr string) string {
	body = encodeURIComponent(body)
	strs := strings.Split(body, "")
	sort.Strings(strs)
	body = strings.Join(strs, "")
	body = base64.StdEncoding.EncodeToString([]byte(body))
	res := GetMD5Encode(body) + GetMD5Encode(ts+":"+randStr)
	res = strings.ToUpper(GetMD5Encode(res))
	return res
}

func GetMD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func unicode(str string) string {
	textQuoted := strconv.QuoteToASCII(str)
	textUnquoted := textQuoted[1 : len(textQuoted)-1]
	return textUnquoted
}
