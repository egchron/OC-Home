package controller

import (
	"log"
	"time"

	"okumoto.net/config"

	"github.com/beevik/timerqueue"
)

type Controller struct {
	Queue *timerqueue.Queue
}

func (c *Controller) Setup(config *config.Config, now time.Time) {
	c.Queue.Schedule(Weekly{c.Queue, config}, now)
	now = now.Add(time.Second)
	c.Queue.Schedule(Daily{c.Queue, config}, now)
	now = now.Add(time.Second)
	c.Queue.Schedule(Hourly{c.Queue, config}, now)
	now = now.Add(time.Second)
	c.Queue.Advance(now)
}

func (c *Controller) MainLoop() {
	for c.Queue.Len() > 0 {
		_, when := c.Queue.PeekFirst()
		log.Printf("Sleep for: %v\n", time.Until(when))
		time.Sleep(time.Until(when))
		c.Queue.Advance(time.Now())
	}

	log.Printf("No more events in mainLoop\n")
}
