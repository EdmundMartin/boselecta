package server

import (
	"errors"
	"strings"
)


func (c *ClientConn) handleGet(cmd string) (string, error) {
	args := strings.Split(cmd, " ")
	if len(args) == 2 {
		return "", errors.New("invalid command")
	}
	namespace, flagName := args[0], args[1]
	fl, err := c.dataStore.GetFlag(namespace, flagName)
	if err != nil {
		return "", err
	}
	return fl.String(), err
}