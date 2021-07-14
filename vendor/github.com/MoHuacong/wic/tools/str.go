package tools

import (
	"encoding/json"
	"strconv"
)

func StrToUpper(str string) string {
	if len(str) < 1 {
		return ""
	}
	strArry := []rune(str)
	if strArry[0] >= 97 && strArry[0] <= 122  {
		strArry[0] -= 32
	}
	return string(strArry)
}

func InterfaceToStr(value interface{}) (string, string) {
	var typ, key string
	if value == nil {
		return "nil", key
	}

	switch value.(type) {
	case float64:
		typ = "float64"
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		typ = "float32"
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		typ = "int"
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		typ = "uint"
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		typ = "int8"
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		typ = "uint8"
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		typ = "int16"
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		typ = "uint16"
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		typ = "int32"
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		typ = "uint32"
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		typ = "int64"
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		typ = "uint64"
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		typ = "string"
		key = value.(string)
	case []byte:
		typ = "[]byte"
		key = string(value.([]byte))
	default:
		typ = "default"
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return typ, key
}