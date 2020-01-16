package controller

import (
	"log"
	"time"

	"okumoto.net/config"

	"github.com/beevik/timerqueue"
)

type Daily struct {
	queue  *timerqueue.Queue
	config *config.Config
}

func (d Daily) OnTimer(now time.Time) {
	log.Printf("Daily Hit\n")

	midnight := time.Date(now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, time.Local)
	tomorrow := midnight.Add(24 * time.Hour)

	log.Printf("Daily Next: %v\n", tomorrow)
	d.queue.Schedule(d, tomorrow)

	d.run(midnight, now) // Do daily stuff.
}

func (d Daily) run(midnight, now time.Time) {
	if d.config.NoWaterDays[now.Weekday()] == false {
		// Schedule today's watering sessions.
		for _, station := range d.config.Stations {
			for _, span := range station.Span {
				go station.Run(midnight, now, span)
			}
		}
	}
}
