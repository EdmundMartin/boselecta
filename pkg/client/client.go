package client

import (
	"bufio"
	"fmt"
	"net"
)

const minBuffer = 1500

type Boselecta struct {
	Host       string
	Port       int
	namespace  string
	connection *net.TCPConn
	bufRead    *bufio.Reader
	bufWrite   *bufio.Writer
	cache      *simpleCache
}

func (b *Boselecta) String() string {
	return fmt.Sprintf("%s:%d", b.Host, b.Port)
}

func NewBoselecta(namespace string) *Boselecta {
	return &Boselecta{
		namespace: namespace,
		cache: &simpleCache{
			make(map[string]*cacheValue),
		},
	}
}

func (b *Boselecta) DialConnection(host string, port int) (*Boselecta, error) {
	b.Host = host
	b.Port = port
	conn, err := net.Dial("tcp", b.String())
	if err != nil {
		return nil, err
	}
	b.bufRead = bufio.NewReader(conn)
	b.bufWrite = bufio.NewWriter(conn)
	b.connection = conn.(*net.TCPConn)
	return b, nil
}

func (b *Boselecta) sendAll(msg []byte) (int, error) {
	written := 0
	forWrite := msg
	var n int
	var err error
	for written < len(msg) {
		forBuff := len(forWrite) >= minBuffer
		if forBuff {
			n, err = sendAllBuffer(b, forWrite)
			if err != nil && !isNetTempErr(err) {
				return written, err
			}
		} else {
			n, err = sendAllNoBuffer(b, forWrite)
			if err != nil && !isNetTempErr(err) {
				return written, err
			}
		}
		written += n
		forWrite = forWrite[n:]
	}
	return written, nil
}

func sendAllNoBuffer(b *Boselecta, msg []byte) (int, error) {
	n, err := b.connection.Write(msg)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func sendAllBuffer(b *Boselecta, msg []byte) (int, error) {
	n, err := b.bufWrite.Write(msg)
	if err != nil {
		return n, err
	}
	err = b.bufWrite.Flush()
	if err != nil {
		return n, err
	}
	return n, nil
}

func (b *Boselecta) getResp(cmd string) (string, error) {
	_, err := b.sendAll([]byte(cmd))
	if err != nil {
		return "", err
	}
	resp, err := b.bufRead.ReadString('\n')
	if err != nil {
		return "", err
	}
	return resp, nil
}

func (b *Boselecta) RetrieveFlag(flagname string) (interface{}, error) {
	val, found := b.cache.RetrieveKey(flagname)
	if found {
		return val, nil
	}
	payload := fmt.Sprintf("%s %s\n", b.namespace, flagname)
	resp, _ := b.getResp(payload)
	msg := parseMessage(resp)
	if msg.hasError {
		return nil, msg.errMsg
	}
	b.cache.SetKey(flagname, msg.rawVal, msg.refresh)
	return msg.rawVal, nil
}

func (b *Boselecta) RetrieveStringFlag(flagname string, defaultVal string) string {
	f, err := b.RetrieveFlag(flagname)
	if err != nil {
		return defaultVal
	}
	var result string
	result = f.(string)
	if result == "" {
		return defaultVal
	}
	return result
}
