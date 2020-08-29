package controller

import (
	"fmt"
	"log"
	"sync"
	"time"

	"okumoto.net/db"
	"okumoto.net/valvebox"

	"github.com/go-co-op/gocron"
	"github.com/pkg/errors"
)

type Controller struct {
	mux sync.Mutex // Used to limit station activation.

	ValveBox  *valvebox.ValveBox
	Scheduler *gocron.Scheduler
	Db        db.DB
}

func (c *Controller) MainLoop() {

	c.Scheduler.StartAsync()

	//c.Scheduler.Every(1).Day().Do(c.DayLoop)
	c.Scheduler.Every(20).Second().Do(c.DayLoop)

	c.Scheduler.StartBlocking()
}

func (c *Controller) DayLoop() error {

	sList, err := c.Db.GetStations()
	if err != nil {
		return errors.Wrap(err, "DayLoop GetStations")
	}
	for _, s := range sList {
		dur := s.Budget / 7
		fmt.Printf("Station: %s\n", s.Name)
		c.CycleStation(s.Name, dur*time.Second)
	}

	return nil
}

func (c *Controller) CycleStation(stationName string, dur time.Duration) {

	// Only allow one station to be active at a time.
	c.mux.Lock()
	defer c.mux.Unlock()

	if err := c.ValveBox.State(stationName, true); err != nil {
		log.Fatal(err)
	}

	time.Sleep(dur)

	if err := c.ValveBox.State(stationName, false); err != nil {
		log.Fatal(err)
	}
}
