package tools

func InterfaceToStrMap(v []interface{}) map[string]interface{} {
	i := 0
	var key string
	ret := make(map[string]interface{})
	for k, _v := range v {
		if i % k == 0 {
			key = _v.(string)
		} else {
			ret[key] = _v
		}
	}
	return ret
}
