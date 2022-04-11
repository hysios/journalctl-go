package journalctl

import "github.com/hysios/utils/convert"

func (e Entry) Message() string {
	switch x := e.MESSAGE.(type) {
	case string:
		return x
	case []byte:
		return string(x)
	case []interface{}:
		return string(sliceToBytes(x))
	default:
		return ""
	}
}

func sliceToBytes(s []interface{}) []byte {
	b := make([]byte, len(s))
	for i, v := range s {
		b[i], _ = convert.Byte(v)
	}
	return b
}
