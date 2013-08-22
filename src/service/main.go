package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"engine/user"
)

type UserService interface {
	Query(*user.UserArg, *user.UserRet) error
}

func Run(clientName, clientUri string) (run func(), err error) {
	service := NewUserService()

	server := rpc.NewServer()
	if err := server.RegisterName(clientName, service); err != nil {
		return nil, err
	}

	run = func() {
		l, err := net.Listen("tcp", clientUri)
		if err != nil {
			log.Fatal("net.Listen: ", err)
		}
		defer l.Close()

		for {
			conn, err := l.Accept()
			if err != nil {
				log.Fatal("l.Accept: ", err)
			}
			defer conn.Close()

			go server.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}

	return run, err
}

type userService struct {
}

func NewUserService() UserService {
	return new(userService)
}

func (u *userService) Query(in *user.UserArg, ret *user.UserRet) error {
	var users = map[string][]string{
		"1": []string{"viney", "viney.chow@gmail.com"},
	}

	if v, ok := users[in.Uid]; ok {
		ret.Name = v[0]
		ret.Email = v[1]
	} else {
		ret.Name = "Null"
		ret.Email = "Null"
	}

	return nil
}

func main() {
	run, err := Run(user.ClientName, user.ClientUri)
	if err != nil {
		log.Fatal("Run: ", err)
		return
	}

	run()
}
