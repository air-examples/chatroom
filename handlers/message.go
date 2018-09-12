package handlers

import (
	"encoding/json"

	"github.com/aofei/air"
)

type Message struct {
	From    string                   `json:"from"`
	Mtype   air.WebSocketMessageType `json:"type"`
	Content string                   `json:"content"`
}

func newMsg(from string, t air.WebSocketMessageType, b []byte) *Message {
	return &Message{
		From:    from,
		Mtype:   t,
		Content: string(b),
	}
}

func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(*m)
}
