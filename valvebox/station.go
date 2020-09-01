package valvebox

import (
	"time"

	"github.com/pkg/errors"
)

// Station represents a single valve in the irrigation system.
//
// One or more relays maybe activated in sequence to pass water.
// For example, the first may activate a pump or a master valve, while
// the last might be the actual solenoid.
type Station struct {
	Name   string
	relays []*Relay

	// It takes time for valves to shut off.  Prevent other stations
	// from opening and lowering the pressure so much that it will
	// not close.
	offLatancy time.Duration
}

func (s Station) state(on bool) error {
	if on {
		for _, r := range s.relays {
			if err := r.state(on); err != nil {
				return errors.Wrap(err, "Station")
			}
		}
	} else {
		// Deactivate the relays in reverse order.
		for i := len(s.relays) - 1; i >= 0; i-- {
			r := s.relays[i]
			if err := r.state(on); err != nil {
				return errors.Wrap(err, "Station")
			}
		}
		time.Sleep(s.offLatancy)
	}
	return nil
}
