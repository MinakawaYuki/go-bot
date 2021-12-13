package tools

import (
	"bytes"
	"github.com/tidwall/gjson"
)

// Bytes2Map 将ws获取的字节流转换为map并返回
func Bytes2Map(p []byte) (data map[string]interface{}) {
	json := bytes.NewBuffer(p).String()
	data, err := gjson.Parse(json).Value().(map[string]interface{})
	if !err {
		panic("解析json错误")
	}

	return data
}

// Bytes2Map2 将ws获取的字节流转换为map并返回
func Bytes2Map2(p []byte) (data []interface{}) {
	data = gjson.ParseBytes(p).Value().([]interface{})

	return data
}
