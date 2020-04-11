package server

import (
	"errors"
	"github.com/EdmundMartin/boselecta/pkg/flag"
	"strings"
)

func (c *ClientConn) handleGet(cmd string) (*flag.FeatureFlag, error) {
	args := strings.Split(cmd, " ")
	if len(args) != 2 {
		return nil, errors.New("invalid command")
	}
	namespace, flagName := args[0], args[1]
	fl, err := c.dataStore.GetFlag(namespace, flagName)
	if err != nil {
		return nil, err
	}
	return fl, err
}
