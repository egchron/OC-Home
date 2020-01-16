package controller

import (
	"log"
	"time"

	"okumoto.net/config"

	"github.com/beevik/timerqueue"
)

type Weekly struct {
	queue  *timerqueue.Queue
	config *config.Config
}

func (w Weekly) OnTimer(now time.Time) {
	log.Printf("Weekly hit\n")

	midnight := time.Date(now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, time.Local)
	offset := 1 + int(time.Saturday-now.Weekday())
	nextweek := midnight.Add(time.Duration(offset) * 24 * time.Hour)

	log.Printf("Weekly Next: %v\n", nextweek)
	w.queue.Schedule(w, nextweek)

	w.run() // Do weekly stuff.
}

func (w *Weekly) run() {
}
