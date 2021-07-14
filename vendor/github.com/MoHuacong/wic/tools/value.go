package tools

import "reflect"

func ValueToSliceInterface(v []reflect.Value) []interface{} {
	data := make([]interface{}, len(v))
	for k, v := range v {
		data[k] = v.Interface()
	}
	return data
}
