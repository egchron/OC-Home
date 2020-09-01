package valvebox

import (
	"github.com/pkg/errors"
)

type Relay struct {
	Name  string
	board RelayBoard
	id    int
}

func (r *Relay) state(on bool) error {
	if err := r.board.Relay(r.id, on); err != nil {
		return errors.Wrap(err, "Relay")
	}
	return nil
}
