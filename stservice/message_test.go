package stservice

import (
	"fmt"
	"testing"
)

func TestEncodeDecodeInt(t *testing.T) {
	var m Message

	type payload int

	foo := payload(1)
	m.SetData(foo)

	bar := m.GetData().(payload)

	if foo != bar {
		t.Fatalf("Correct values not recovered: foo:%v bar:%v", foo, bar)
	}
}

func TestEncodeDecodeStruct(t *testing.T) {
	var m Message

	type payload struct {
		A int
		b int
	}

	foo := payload{1, 2}
	m.SetData(foo)

	bar := m.GetData().(payload)

	if foo != bar {
		t.Fatalf("Correct values not recovered: foo:%v bar:%v", foo, bar)
	}
}

func TestSendReceive(t *testing.T) {
	ch := make(chan Message, 1)

	type payload int

	// Sending data to server
	foo := payload(1)
	m1 := NewMessage(foo)
	ch <- m1

	// Server received data
	m2 := <-ch
	bar := m2.GetData().(payload)

	// Server processes data and replies
	baz := fmt.Sprint(bar)
	m2.SetData(baz)
	m2.Reply()

	// Client receives response
	m3 := m1.Receive()
	qux := m3.GetData().(string)

	if foo != bar {
		t.Fatalf("Correct values not recovered: foo:%v bar:%v", foo, bar)
	}

	if baz != qux {
		t.Fatalf("Correct values not recovered: baz:%v qux:%v", baz, qux)
	}

}
