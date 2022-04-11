package util

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/Khighness/entry-task/tcp/common"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

// GenerateToken 生成token
// '0' 48
// 'A' 65
func GenerateToken() string {
	//var sessionId bytes.Buffer

	// 生成随机16位byte
	buf := make([]byte, common.TokenBytes)
	for i := 0; i < common.TokenBytes; i++ {
		buf[i] = byte(Uint32())
	}

	// md5计算消息摘要
	hash := md5.New()
	hash.Write(buf)
	buf = hash.Sum(nil)

	return hex.EncodeToString(buf) //test: 670ns

	// 转换为十六进制大写字符串
	//for i := 0; i < common.TokenBytes; i++ {
	//	var b1 byte = (buf[i] & 0xf0) >> 4
	//	var b2 byte = buf[i] & 0x0f
	//	if b1 < 10 {
	//		sessionId.WriteByte(48 + b1)
	//	} else {
	//		sessionId.WriteByte(55 + b1)
	//	}
	//	if b2 < 10 {
	//		sessionId.WriteByte(48 + b2)
	//	} else {
	//		sessionId.WriteByte(55 + b2)
	//	}
	//}
	//
	//return sessionId.String() // test: 1 us
}
