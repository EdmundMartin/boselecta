package main

import (
	"fmt"
	"github.com/EdmundMartin/boselecta/internal/server"
	"github.com/EdmundMartin/boselecta/internal/storage/mongoStorage"
	"log"
	"net"
	"flag"
)

func tcpServer(port int) (net.Listener, error) {
	l, err := net.Listen("tcp4", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return l, err
}

func main()  {
	var storageHost string
	var storageDB string
	var storageCol string
	var port int

	flag.StringVar(&storageHost, "host", "mongodb://localhost:27017", "MongoDB URI")
	flag.StringVar(&storageDB, "db", "featureFlag", "MongoDB DB name")
	flag.StringVar(&storageCol, "collection", "flags", "MongoDB Collection name")
	flag.IntVar(&port, "port", 11300, "TCP Server Port")

	l, err := tcpServer(port)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	storage, err := mongoStorage.NewMongo(storageHost, storageDB, storageCol)
	if err != nil {
		log.Fatal(err)
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
