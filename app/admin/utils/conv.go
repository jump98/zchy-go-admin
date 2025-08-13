package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func GetInterfaceInt64Value(vi interface{}) (int64, error) {
	i := int64(0)
	var err error = nil
	switch v := vi.(type) {
	case int64:
		// 直接使用v
		i = v
	case float64:
		// 转换为int64
		i = int64(v)
	case string:
		// 字符串转int64
		i, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, err
		}
	case json.Number:
		i, err = v.Int64()
		if err != nil {
			return 0, err
		}
	default:
		return 0, fmt.Errorf("unexpected radarId type: %T", v)
	}

	return i, nil
}
