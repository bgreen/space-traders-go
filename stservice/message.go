package stservice

import (
	"reflect"
)

type Message struct {
	Data      any
	ReplyChan chan Message
	Err       error
}

func NewMessage(d any) Message {
	return Message{Data: d, ReplyChan: make(chan Message)}
}

func (m *Message) SetData(r any) {
	m.Data = r
}

func (m Message) GetData() any {
	return m.Data
}

func (m Message) GetType() reflect.Type {
	return reflect.TypeOf(m.Data)
}

func (m Message) Reply() {
	m.ReplyChan <- m
}

func (m Message) Receive() Message {
	return <-m.ReplyChan
}
