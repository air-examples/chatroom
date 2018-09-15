package handlers

import "sync"

type SocketManager struct {
	name     string
	msg      *Message
	newMsg   chan struct{}
	shutdown chan struct{}
	mu       *sync.Mutex
}

func newSocketManager(name string) *SocketManager {
	return &SocketManager{
		name:     name,
		newMsg:   make(chan struct{}, 1),
		shutdown: make(chan struct{}),
		mu:       &sync.Mutex{},
	}
}

func (sm *SocketManager) SendMsg(msg *Message) {
	for _, i := range users.Keys() {
		value, _ := users.Get(i)
		v, _ := value.(*SocketManager)
		v.mu.Lock()
		v.msg = msg
		v.newMsg <- struct{}{}
	}
}

func (sm *SocketManager) Close() {
	sm.shutdown <- struct{}{}
	users.Remove(sm.name)
}

func CloseSocket() {
	for _, i := range users.Keys() {
		value, _ := users.Get(i)
		v, _ := value.(*SocketManager)
		v.Close()
	}
}
