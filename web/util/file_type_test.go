package util

import (
	"io/ioutil"
	"os"
	"testing"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-03-01

func TestGetFileType(t *testing.T) {
	file, err := os.Open("/Users/zikang.chen/Pictures/khighness.jpg")
	if err != nil {
		t.Log("open file err:", err)
	}
	fSrc, _ := ioutil.ReadAll(file)
	t.Log(GetFileType(fSrc[:10]))
}
