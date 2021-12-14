package tools

import "math/big"

// GetString 类型断言
func GetString(v interface{}) string {
	var str string
	switch v.(type) {
	case float64:
		vv := v.(float64)
		data := big.NewRat(1, 1)
		data.SetFloat64(vv)
		str = data.FloatString(0)
		break
	case string:
		str = v.(string)
		break
	}
	return str
}
