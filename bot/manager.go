package bot

import (
	"fmt"

	"github.com/bgreen/space-traders-go/st"
	"github.com/bgreen/space-traders-go/stapi"
)

type msg any

type shipIdleMsg stapi.Ship

type errorMsg error

type quitMsg bool

type batchMsg []msg

type Manager struct {
	queue  chan msg
	client *st.Client
}

func (m *Manager) Run() error {

	ships, err := m.client.GetMyShips()
	if err != nil {
		return err
	}

	go func() {
		for _, s := range ships {
			m.queue <- shipIdleMsg(s)
		}
	}()

	for {

		switch msg := (<-m.queue).(type) {
		case errorMsg:
			fmt.Println(msg)
			return msg
		case quitMsg:
			return nil
		case shipIdleMsg:
			b := m.AllocateShipBot(stapi.Ship(msg))
			m.RunBot(b)
		case botDoneMsg:
			m.RunBot(msg)
		case batchMsg:
			go func() {
				for _, v := range msg {
					m.queue <- v
				}
			}()
		}
	}
}

func NewManager(c *st.Client) *Manager {
	return &Manager{
		queue:  make(chan msg, 10),
		client: c,
	}
}

func (m Manager) AllocateShipBot(ship stapi.Ship) Bot {
	var b Bot
	if isShipMiner(ship) {
		b = MinerBot{
			Ship:   ship,
			Client: m.client,
		}
	}
	return b
}

func (m Manager) RunBot(b Bot) {
	if b != nil {
		go func() {
			m.queue <- b.RunOnce()
		}()
	}
}
