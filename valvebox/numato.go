package valvebox

import (
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/pkg/errors"
	"github.com/pkg/term"
)

type Numato struct {
	DevName string

	term *term.Term
}

func (n *Numato) Open() error {
	term, err := term.Open(n.DevName,
		term.Speed(9600),
		term.FlowControl(term.NONE),
		term.RawMode,
	)
	if err != nil {
		return errors.Wrap(err, "term.Open()")
	}

	if err := term.Flush(); err != nil {
		return errors.Wrap(err, "term.Flush()")
	}

	n.term = term
	return nil
}

func (n *Numato) write(cmd string) error {
	_, err := n.term.Write([]byte(cmd))
	if err == nil {
		return nil // Write success.
	}

	if err == io.ErrShortWrite {
		// The term library returned short write, assume
		// that something happened to the device.  Close it
		// so it will be initialized later.
		n.term.Close()
		n.term = nil
		return errors.Wrapf(err, "Write(%s): Short write", cmd)
	}

	switch e := err.(type) {
	case *os.PathError:
		if e.Err == syscall.ENXIO {
			// The write failed with a "No such device
			// or address", close it so it will be
			// initialized later.
			n.term.Close()
			n.term = nil
			return errors.Wrapf(err, "Write(%s): No such device", cmd)
		}

		fmt.Printf("GAM: <%T> %v\n", e.Err, e.Err)
		return err
	default:
		fmt.Printf("MAX: <%T> %v\n", err, err)
		return errors.Wrapf(err, "Write(%s)", cmd)
	}
}

func (n *Numato) Relay(id int, on bool) error {

	if n.term == nil {
		if err := n.Open(); err != nil {
			return errors.Wrap(err, "n.Open()")
		}
	}

	var state string
	if on {
		state = "on"
	} else {
		state = "off"
	}

	cmd := fmt.Sprintf("relay %s %x\r", state, id)
	if err := n.write(cmd); err != nil {
		return errors.Wrapf(err, "relay(%d, %s)", id, state)
	}
	return nil
}
