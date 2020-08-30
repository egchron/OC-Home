package valvebox

import (
	"sync"
)

type Relay struct {
	Id   int
	Name string
}

type ValveBox struct {
	mux    sync.Mutex
	rb     RelayBoard
	relays map[string]int
}

func New(rb RelayBoard) *ValveBox {

	m := ValveBox{
		rb:     rb,
		relays: make(map[string]int),
	}

	return &m
}

func (v *ValveBox) Add(id int, name string) {
	v.relays[name] = id
}
