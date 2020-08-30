package valvebox

import (
	"okumoto.net/scheduler"
)

type Station struct {
	Name     string
	Span     []scheduler.Span
	ValveBox *ValveBox
}
