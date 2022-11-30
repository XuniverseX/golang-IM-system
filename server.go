package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// create new server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}

	return server
}

func (server *Server) handler(conn net.Conn) {
	// 当前链接的业务..
	fmt.Println("链接建立成功!")
}

// 启用服务器的接口
func (server *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("net.Listener err:", err)
		return
	}
	// close listen socket
	defer listener.Close()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err!")
			return
		}

		// do Handler
		go server.handler(conn)
	}
}
