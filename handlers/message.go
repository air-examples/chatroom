package handlers

import (
	"encoding/json"
	"time"
)

type Message struct {
	From    string `json:"from"`
	Type    string `json:"type"`
	Content string `json:"content"`
	Time    string `json:"time"`
}

func newMsg(from string, b string) *Message {
	m := Map{}
	_ = json.Unmarshal([]byte(b), &m)
	content, _ := m["content"].(string)
	msgType, _ := m["type"].(string)
	return &Message{
		From:    from,
		Type:    msgType,
		Content: content,
		Time:    time.Now().Format("15:04:05"),
	}
}

func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(*m)
}
