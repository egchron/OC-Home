package controller

import (
	"log"
	"time"

	"okumoto.net/config"

	"github.com/beevik/timerqueue"
)

type Hourly struct {
	queue  *timerqueue.Queue
	config *config.Config
}

func (h Hourly) OnTimer(now time.Time) {
	log.Printf("Hourly Hit\n")

	topofhour := time.Date(now.Year(), now.Month(), now.Day(),
		now.Hour(), 0, 0, 0, time.Local)
	nextHour := topofhour.Add(time.Hour)

	log.Printf("Hourly Next: %v\n", nextHour)
	h.queue.Schedule(h, nextHour)

	h.run(now) // Do hourly stuff.
}

func (h Hourly) run(now time.Time) {
}
