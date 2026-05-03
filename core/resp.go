package resp

type Command struct {
}

func Decode(data []byte) (interface{}, bool) {

	val, _, err := DecodeOne(data)
	return val, false
}

func DecodeOne(buf []byte) (interface{}, int, bool) {
	if len(buf) == 0 {
		return nil, 0, true
	}

	switch buf[0] {
	case '+':
		return readSimpleString(buf)
	case '-':
		return readError(buf)
	case ':':
		return readInt64(buf)
	case '$':
		return readBulkString(buf)
	case '*':
		return readArray(buf)
	}

	return nil, 0, false
}

func readSimpleString(buf []byte) (string, int, bool) {
	pos := 1

	for ; buf[pos] != '\r'; pos++ {
	}

	return string(buf[1:pos]), pos + 2, false
}
func readError(buf []byte) (string, int, bool) {
	return readSimpleString(buf)
}
func readInt64(buf []byte) (int64, int, bool) {
	if len(buf) < 3 { //|| buf[0] != ':' already checked
		return 0, 0, true
	}

	pos := 1
	sign := int64(1)

	if buf[pos] == '-' {
		sign = -1
		pos++
	}

	var res int64

	for pos < len(buf) && buf[pos] != '\r' {
		b := buf[pos]
		if b < '0' || b > '9' {
			return 0, pos, true
		}
		res = res*10 + int64(b-'0')
		pos++
	}

	if pos+1 >= len(buf) || buf[pos] != '\r' || buf[pos+1] != '\n' {
		return 0, pos, true
	}

	return sign * res, pos + 2, false
}
func readBulkString(buf []byte) (string, int, bool) {

	return "test", 0, true
}
