package main

import (
	"fmt"
	"github.com/EdmundMartin/boselecta/pkg/server"
	"github.com/EdmundMartin/boselecta/pkg/storage/mongoStorage"
	"net"
)

func main()  {
	l, err := net.Listen("tcp4", ":11300")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	storage, err := mongoStorage.NewMongo("mongodb://localhost:27017", "featureFlag", "flags")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		conn := server.NewClientConn(c, storage)
		go conn.HandleConn()
	}
}
