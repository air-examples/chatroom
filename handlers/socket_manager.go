package handlers

type SocketManager struct {
	name      string
	msg       *Message
	newMsg    chan struct{}
	shutdown  chan struct{}
	writeChan chan struct{}
	sendChan  chan struct{}
}

func newSocketManager(name string) *SocketManager {
	return &SocketManager{
		name:      name,
		newMsg:    make(chan struct{}, 1),
		shutdown:  make(chan struct{}),
		writeChan: make(chan struct{}, 1),
		sendChan:  make(chan struct{}, 1),
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
		v.writeChan <- struct{}{}
		v.msg = msg
		v.newMsg <- struct{}{}
	}
}

func (sm *SocketManager) Close() {
	users.Remove(sm.name)
	sm.shutdown <- struct{}{}
}

func CloseSocket() {
	for _, i := range users.Keys() {
		value, _ := users.Get(i)
		v, _ := value.(*SocketManager)
		v.Close()
	}
}
