package server

import (
	"fmt"
	"github.com/EdmundMartin/boselecta/pkg/flag"
)

func sendSuccessResponse(c *ClientConn, f *flag.FeatureFlag) error {
	response := fmt.Sprintf("OK %s %d %s\n", f.Type.String(), f.Refresh, f.String())
	_, err := c.SendAll([]byte(response))
	return err
}

func (c *ClientConn) HandleConn() {
	for c.scanner.Scan() {
		cmd := c.scanner.Text()
		f, err := c.handleGet(cmd)
		if err != nil {
			errorResp := fmt.Sprintf("ERROR %s\n", err.Error())
			c.SendAll([]byte(errorResp))
		} else {
			sendSuccessResponse(c, f)
		}
	}
}
