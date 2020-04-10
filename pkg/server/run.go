package server

import "fmt"

func (c *ClientConn) handleConn() {
	for c.scanner.Scan() {
		cmd := c.scanner.Text()
		response, err := c.handleGet(cmd)
		if err != nil {
			errorResp := fmt.Sprintf("ERROR %s\n", err.Error())
			c.SendAll([]byte(errorResp))
			return
		} else {
			response = fmt.Sprintf("OK %s\n", response)
			c.SendAll([]byte(response))
		}
	}
}