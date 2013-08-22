package user

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
)

const (
	ClientName = "user"
	ClientUri  = "127.0.0.1:8080"
)

var mutex sync.Mutex

type UserArg struct {
	Uid string
}

type UserRet struct {
	Name  string
	Email string
}

func Query(uid string) (string, string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	client, release, err := dial()
	if err != nil {
		return "", "", err
	}
	defer func() {
		release()
		client.Close()
	}()

	in := &UserArg{uid}
	ret := &UserRet{}
	if err := client.Call(ClientName+".Query", in, ret); err != nil {
		return "", "", err
	}

	return ret.Name, ret.Email, nil
}

func dial() (*rpc.Client, func(), error) {
	conn, err := net.Dial("tcp", ClientUri)
	if err != nil {
		return nil, nil, err
	}

	release := func() {
		conn.Close()
	}

	return jsonrpc.NewClient(conn), release, nil
}
