package valvebox

import (
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/pkg/errors"
	"github.com/pkg/term"
)

// Numato represents a USB serial relay board made by "Numato Lab".
type Numato struct {
	devName string

	term *term.Term
}

// Open sets up the device for reading and writing.
func (n *Numato) Open(devName string) error {
	if err := n.setup(devName); err != nil {
		return errors.Wrap(err, "Numato setup")
	}
	n.devName = devName
	return nil
}

// Relay changes the state of the relay to on or off.
func (n *Numato) Relay(id int, on bool) error {
	if err := n.setup(n.devName); err != nil {
		return errors.Wrap(err, "Numato re-setup")
	}

	var state string
	if on {
		state = "on"
	} else {
		state = "off"
	}

	cmd := fmt.Sprintf("relay %s %x\r", state, id)
	if err := n.write(cmd); err != nil {
		return errors.Wrapf(err,
			"Numato relay(%d, %s)", id, state)
	}
	return nil
}

func (n *Numato) setup(devName string) error {
	if n.term != nil {
		return nil // Already setup.
	}

	term, err := term.Open(devName,
		term.Speed(9600),
		term.FlowControl(term.NONE),
		term.RawMode,
	)
	if err != nil {
		return errors.Wrap(err, "Numato open")
	}

	if err := term.Flush(); err != nil {
		return errors.Wrap(err, "Numato flush")
	}

	n.term = term
	return nil
}

func (n *Numato) write(cmd string) error {
	b, err := n.term.Write([]byte(cmd))
	switch e := err.(type) {
	case nil:
		return nil // Write successful.

	case *os.PathError:
		if e.Err == syscall.ENXIO {
			// The write failed with a "No such device
			// or address", close it so it will be
			// initialized later.
			n.term.Close()
			n.term = nil
			return errors.Wrap(err,
				"Numato write")
		}
		return errors.Wrap(err,
			"Numato unexpected PathError")

	default:
		if err == io.ErrShortWrite {
			// The term library returned short write, assume
			// that something happened to the device.
			// Close it so it will be re-initialized later.
			n.term.Close()
			n.term = nil
			return errors.Wrapf(err,
				"Numato wrote %d bytes", b)
		}
		return errors.Wrapf(err,
			"Numato unexpected error type '%T'", err)
	}
}
