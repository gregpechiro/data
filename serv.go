package data

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"time"
)

type M struct {
	Cmd string      `json:"cmd,omitempty"`
	Key string      `json:"key,omitempty"`
	Val interface{} `json:"val,omitempty"`
}

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
		go self.Handle(conn)
	}
}

func (self *Server) Handle(conn net.Conn) {
	enc, dec := json.NewEncoder(conn), json.NewDecoder(conn)
	for {
		var m M
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			conn.Close()
			return
		} else {
			conn.SetDeadline(time.Now().Add(time.Second * 300))
		}
		switch m.Cmd {
		case "set":
			enc.Encode(self.st.Set(&m))
		case "get":
			enc.Encode(self.st.Get(&m))
		default:
			enc.Encode(M{Val: "false"})
		}
	}
}
