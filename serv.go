package data

import (
	"log"
	"net"
)

type Server struct {
	st *Store
}

func InitServer() *Server {
	return &Server{
		st: InitStore(),
	}
}

func (self *Server) Listen(port string) {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, _ := ln.Accept()
		go self.st.Handle(conn)
	}
}
