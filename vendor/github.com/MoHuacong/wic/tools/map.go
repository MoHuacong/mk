package tools

type AnyMaps map[Any]Any
type AnyListMap map[Any]Any

type ListMap struct {
	idUrl AnyListMap
	urlId map[Any]map[Any]bool
}

func NewListMap() *ListMap {
	return &ListMap{make(AnyListMap), make(map[Any]map[Any]bool)}
}

func (lm *ListMap) Is(args ...Any) bool {
	key := args[0]
	var value Any
	if len(args) >= 2 { value = args[1] }

	if value == nil {
		if lm.IdKey(key) != nil { return true }
		if lm.IdKey(value) != nil { return true }
		return false
	}

	if lm.UrlKey(key, value) != nil { return true }
	if lm.UrlKey(value, key) != nil { return true }
	return false
}

func (lm *ListMap) IdKey(key Any) Any {
	if lm.idUrl[key] != nil {
		return lm.idUrl[key]
	}
	return nil
}

func (lm *ListMap) UrlKey(key, value Any) Any {
	v := lm.UrlKeyList(key)
	if v == nil { return nil }
	if v[value] == false { return nil }
	return v[value]
}

func (lm *ListMap) UrlKeyList(key Any) map[Any]bool {
	if lm.urlId[key] == nil { return nil }
	return lm.urlId[key]
}

func (lm *ListMap) Key(args ...Any) Any {
	key := args[0]
	var value Any
	if len(args) >= 2 { value = args[1] }

	if value == nil {
		v1 := lm.IdKey(key)
		if v1 != nil { return v1 }

		v2 := lm.IdKey(value)
		if v2 != nil { return v2 }
		return false
	}

	v3 := lm.UrlKey(key, value)
	if v3 != nil { return v3 }

	v4 := lm.UrlKey(value, key)
	if v4 != nil { return v4 }

	return nil
}

func (lm *ListMap) IdSet(key, value Any) bool {
	lm.idUrl[key] = value
	return true
}

func (lm *ListMap) UrlSet(key, value Any) bool {
	if lm.urlId[key] == nil {
		lm.urlId[key] = make(map[Any]bool)
	}
	lm.urlId[key][value] = true
	return lm.urlId[key][value]
}

func (lm *ListMap) Set(key, value Any) bool {
	if !lm.IdSet(key, value) { return false }
	if !lm.UrlSet(value, key) { return false }
	return true
}

func (lm *ListMap) IdDelete(key Any) bool {
	if lm.idUrl[key] != nil  {
		delete(lm.idUrl, key)
		return true
	}
	return true
}

func (lm *ListMap) UrlDelete(args ...Any) bool {
	key := args[0]
	var value Any
	if len(args) >= 2 { value = args[1] }

	if value == nil {
		if lm.urlId[key] == nil { return false }
		delete(lm.urlId, key)
		return true
	}

	if lm.urlId[key][value] == false { return false }
	delete(lm.urlId[key], value)
	return true
}

func (lm *ListMap) Delete(args ...Any) bool {
	if !lm.IdDelete(args[0]) { return false }
	if len(args) == 1 {
		if !lm.UrlDelete(args[0]) { return false }
	} else if len(args) >= 2 {
		if !lm.UrlDelete(args[0], args[1]) { return false }
	}

	return true
}