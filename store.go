package data

import (
	"sync"
)

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
