package controller

import (
	"fmt"
	"log"
	"sync"
	"time"

	"okumoto.net/valvebox"

	"github.com/go-co-op/gocron"
)

type Controller struct {
	mux sync.Mutex // Used to limit station activation.

	ValveBox  *valvebox.ValveBox
	Scheduler *gocron.Scheduler
}

func (c *Controller) MainLoop() {

	c.Scheduler.StartAsync()

	c.Scheduler.Every(1).Day().Do(c.DayLoop)

	c.Scheduler.StartBlocking()
}

func (c *Controller) DayLoop() {

	stationNames := []string{
		"GrassHouse",
		"GrassFence",
		"PlanterBoxes",
		"Drip",
	}
	for _, stationName := range stationNames {
		fmt.Printf("Station: %s\n", stationName)
		c.CycleStation(stationName, 4*time.Second)
	}
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
