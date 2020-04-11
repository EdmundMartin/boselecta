package client

import (
	"errors"
	"strconv"
	"strings"
)

func castValue(val, valType string) interface{} {
	switch valType {
	case "Integer":
		val, _ := strconv.Atoi(val)
		return val
	case "String":
		return val
	case "Float":
		val, _ := strconv.ParseFloat(val, 64)
		return val
	default:
		return ""
	}
}

type parsedMessage struct {
	hasError bool
	errMsg   error
	refresh  int
	rawVal   interface{}
}

func parseMessage(msg string) *parsedMessage {
	vals := strings.Split(msg, " ")
	if len(vals) == 2 {
		err := errors.New(vals[1])
		return &parsedMessage{hasError: true, errMsg: err}
	}
	refresh, _ := strconv.Atoi(strings.TrimSpace(vals[3]))
	rawVal := castValue(vals[1], vals[2])
	return &parsedMessage{rawVal: rawVal, refresh: refresh, hasError: false}
}
