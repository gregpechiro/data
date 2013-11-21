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

func (self *Store) Set(m *M) {
	self.Lock()
	self.data[m.Key] = m.Val
	_, ok := self.data[m.Key]
	self.Unlock()
	m.Cmd, m.Key, m.Val = "", "", ok

}

func (self *Store) Get(m *M) {
	self.RLock()
	m.Cmd, m.Key, m.Val = "", "", self.data[m.Key]
	self.RUnlock()
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
			self.Set(&m)
		case "get":
			self.Get(&m)
		default:
			m.Cmd, m.Key, m.Val = "", "", nil
		}
		enc.Encode(m)
	}
}
