package server

import (
	"fmt"
	"github.com/EdmundMartin/boselecta/internal/flag"
	"log"
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
			_, cErr := c.SendAll([]byte(errorResp))
			if cErr != nil {
				log.Printf("encountered error sending ERROR response: %s", cErr)
				return
			}
		} else {
			cErr := sendSuccessResponse(c, f)
			if cErr != nil {
				log.Printf("encounted error sending response: %s", cErr)
				return
			}
		}
	}
	return
}
