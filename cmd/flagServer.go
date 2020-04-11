package main

import (
	"fmt"
	"github.com/EdmundMartin/boselecta/pkg/flag"
	"github.com/EdmundMartin/boselecta/pkg/server"
	"github.com/EdmundMartin/boselecta/pkg/storage/simpleDisk"
	"net"
)

func main()  {
	l, err := net.Listen("tcp4", ":11300")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	f := &flag.FeatureFlag{
		Namespace: "example-namespace",
		FlagName:  "hello",
		Value:     100,
		Type:      flag.IntegerFlag,
		Refresh:   100,
	}
	storage := simpleDisk.NewDiskStore(true)
	storage.Create("example-namespace", f)
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		conn := server.NewClientConn(c, storage)
		//go conn.handleConnection()
		go conn.HandleConn()
	}
}
