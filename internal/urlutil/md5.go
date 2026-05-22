package urlutil

import (
	"crypto/md5"
	"encoding/hex"
)

// EncodeMD5 md5 encryption
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}

// EncodeMD5WithSalt 对字符串进行 MD5 加密，并添加随机盐值 salt
func EncodeMD5WithSalt(value string, salt []byte) string {
	// 将盐值和原始字符串拼接起来
	saltedValue := append([]byte(value), salt...)

	// 创建 MD5 实例，写入数据并计算哈希值
	m := md5.New()
	m.Write(saltedValue)
	hash := m.Sum(nil)

	// 将哈希值进行十六进制编码并返回
	return hex.EncodeToString(hash)
}
