package data

import (
	"encoding/json"
	"io"
	"net"
	"sync"
	"time"
)

type M struct {
	Cmd string      `json:"cmd,omitempty"`
	Key string      `json:"key,omitempty"`
	Val interface{} `json:"val,omitempty"`
}

type Store struct {
	data map[string]interface{}
	sync.RWMutex
}

func InitStore() *Store {
	return &Store{
		data: make(map[string]interface{}),
	}
}

func (self *Store) Set(m *M) M {
	self.Lock()
	self.data[m.Key] = m.Val
	_, ok := self.data[m.Key]
	self.Unlock()
	return M{Val: ok}
}

func (self *Store) Get(m *M) M {
	self.RLock()
	if v, ok := self.data[m.Key]; ok {
		self.RUnlock()
		return M{Val: v}
	}
	self.RUnlock()
	return M{Val: false}
}

func (self *Store) Handle(conn net.Conn) {
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
			enc.Encode(self.Set(&m))
		case "get":
			enc.Encode(self.Get(&m))
		default:
			enc.Encode(M{Val: "false"})
		}
	}
}
