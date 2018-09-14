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
	for _, v := range users {
		v.mu.Lock()
		v.msg = msg
		v.newMsg <- struct{}{}
	}
}

func (sm *SocketManager) Close() {
	sm.shutdown <- struct{}{}
	delete(users, sm.name)
}

func CloseSocket() {
	for _, v := range users {
		v := v
		go func() {
			v.Close()
		}()
	}
}
