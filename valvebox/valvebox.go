package valvebox

import (
	"fmt"
	"log"
	"sync"
	"time"

	"okumoto.net/scheduler"

	"github.com/pkg/errors"
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

// Cycle opens the specified valve for the specified duration.  We
// acquire a lock to make sure only one valve is open at a time.
func (m *ValveBox) Cycle(name string, span scheduler.Span) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	id, ok := m.relays[name]
	if !ok {
		return errors.Errorf("%s is not a known valve.", name)
	}

	// Open named valve.
	msg := fmt.Sprintf("%s-Open(%s)\n", span.Name, name)
	log.Printf(msg)
	if err := m.rb.Relay(id, true); err != nil {
		return errors.Wrap(err, msg)
	}

	time.Sleep(span.Duration) // Sleep until duration passes.

	// Close named valve.
	msg = fmt.Sprintf("%s-Close(%s)\n", span.Name, name)
	log.Printf(msg)
	if err := m.rb.Relay(id, false); err != nil {
		return errors.Wrap(err, msg)
	}

	// Wait for valve to close.
	time.Sleep(5 * time.Second)
	return nil
}
