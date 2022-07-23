package util

import (
	"bytes"
	"encoding/hex"
	"strconv"
	"strings"
	"sync"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-03-01

var fileTypeMap sync.Map

// see https://www.garykessler.net/library/file_sigs.html
func init() {
	fileTypeMap.Store("ffd8ff", "jpg")         //JPEG (jpg)
	fileTypeMap.Store("89504e47", "png")       //PNG (png)
	fileTypeMap.Store("47494638", "gif")       //GIF (gif)
	fileTypeMap.Store("49492a00", "tif")       //TIFF (tif)
	fileTypeMap.Store("424d", "bmp")           //Windows Bitmap(bmp)
	fileTypeMap.Store("41433130", "dwg")       //CAD (dwg)
	fileTypeMap.Store("255044462d312e", "pdf") //Adobe Acrobat (pdf)
}

// GetFileType 通过文件前几个字节判断文件内容类型
func GetFileType(fSrc []byte) string {
	var fileType string
	fileCode := bytesToHexString(fSrc)

	fileTypeMap.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		if strings.HasPrefix(fileCode, strings.ToLower(k)) ||
			strings.HasSuffix(k, strings.ToLower(fileCode)) {
			fileType = v
			return false
		}
		return true
	})
	return fileType
}

// bytesToHexString 将字节数组转换为十六进制字符串
func bytesToHexString(src []byte) string {
	res := bytes.Buffer{}
	if src == nil || len(src) <= 0 {
		return ""
	}
	temp := make([]byte, 0)
	for _, v := range src {
		sub := v & 0xFF
		hv := hex.EncodeToString(append(temp, sub))
		if len(hv) < 2 {
			res.WriteString(strconv.FormatInt(int64(0), 10))
		}
		res.WriteString(hv)
	}
	return res.String()
}
