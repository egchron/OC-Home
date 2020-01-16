package valvebox

import (
	"log"
	"time"

	"okumoto.net/scheduler"

	"github.com/pkg/errors"
)

type Station struct {
	Name     string
	Span     []scheduler.Span
	ValveBox *ValveBox
}

func (s Station) Run(midnight, now time.Time, span scheduler.Span) {

	log.Printf("Scheduling %s-%s\n", span.Name, s.Name)

	when := time.Date(
		midnight.Year(), midnight.Month(), midnight.Day(),
		span.Start.Hour(), span.Start.Minute(), span.Start.Second(), 0,
		time.Local)
	time.Sleep(when.Sub(now)) // Sleep until scheduled time.

	log.Printf("Wakeup %s-%s\n", span.Name, s.Name)

	err := s.ValveBox.Cycle(s.Name, span)

	switch e0 := err.(type) {
	case nil:
		// No errors.

	default:
		log.Printf("Failed cycling %s %v %T %v T%\n",
			s.Name,
			e0, e0,
			errors.Cause(e0), errors.Cause(e0),
		)
	}
}
