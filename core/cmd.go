package core

import (
	"errors"
	"strings"
)

type Command struct {
	Cmd  string
	Args []string
}

// func DecodeArrayCmd(buf []byte) (*Command, error) {

// 	tokens := make([]string, len(ts))

// 	for i := range len(ts) {
// 		tokens[i] = ts[i].(string)
// 	}

// 	return &Command{
// 		Cmd:  strings.ToUpper(tokens[0]),
// 		Args: tokens[1:],
// 	}, nil
// }

func DecodeRESP(buf []byte) (*Command, error) {
	if len(buf) == 0 {
		return nil, errors.New("empty buffer")
	}

	switch buf[0] {
	case 'p', 'P', '-', ':', '$', '+', '*':
		return decode(buf)
	default:
		return nil, errors.New("unknown RESP type")
	}

}

func decode(buf []byte) (*Command, error) {
	s := string(buf)
	tokens := strings.Split(s, "\r\n")

	if strings.EqualFold(tokens[0], "ping") {
		return &Command{
			Cmd:  strings.ToUpper(tokens[0]),
			Args: tokens[1:],
		}, nil
	}

	return &Command{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}
