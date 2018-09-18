package handlers

import "sync"

type SocketManager struct {
	name     string
	msg      *Message
	newMsg   chan struct{}
	shutdown chan struct{}
	mu       *sync.Mutex
	sendChan chan struct{}
}

func newSocketManager(name string) *SocketManager {
	return &SocketManager{
		name:     name,
		newMsg:   make(chan struct{}, 1),
		shutdown: make(chan struct{}),
		mu:       &sync.Mutex{},
		sendChan: make(chan struct{}, 1),
	}
}

func (sm *SocketManager) SendMsg(msg *Message) {
	sm.sendChan <- struct{}{}
	defer func() {
		<-sm.sendChan
	}()
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
