package tools

func Lengths(buf []byte) []byte {
	for k, v := range buf {
		if v == 0 {
			return buf[:k]
		}
	}
	return nil
}