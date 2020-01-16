package config

import (
	"time"

	"okumoto.net/scheduler"
	"okumoto.net/valvebox"
)

type Config struct {
	Stations    []valvebox.Station
	NoWaterDays map[time.Weekday]bool
}

func New(vb *valvebox.ValveBox) *Config {

	start0 := time.Date(0, 0, 0, 0, 0, 0, 0, time.Local)
	start1 := time.Date(0, 0, 0, 0, 15, 0, 0, time.Local)
	start2 := time.Date(0, 0, 0, 12, 00, 0, 0, time.Local)
	duration := 1000 * time.Millisecond

	stations := []valvebox.Station{
		{
			Name: "GrassHouse",
			Span: []scheduler.Span{
				scheduler.Span{
					Name:     "start0",
					Start:    start0,
					Duration: duration},
				scheduler.Span{
					Name:     "start1",
					Start:    start1,
					Duration: duration},
			},
			ValveBox: vb,
		},
		{
			Name: "GrassFence",
			Span: []scheduler.Span{
				scheduler.Span{
					Name:     "start0",
					Start:    start0,
					Duration: duration},
				scheduler.Span{
					Name:     "start1",
					Start:    start1,
					Duration: duration},
			},
			ValveBox: vb,
		},
		{
			Name: "Drip",
			Span: []scheduler.Span{
				scheduler.Span{
					Name:     "start0",
					Start:    start0,
					Duration: duration},
				scheduler.Span{
					Name:     "start1",
					Start:    start1,
					Duration: duration},
				scheduler.Span{
					Name:     "start2",
					Start:    start2,
					Duration: duration},
			},
			ValveBox: vb,
		},
		{
			Name: "PlanterBoxes",
			Span: []scheduler.Span{
				scheduler.Span{
					Name:     "start0",
					Start:    start0,
					Duration: duration},
				scheduler.Span{
					Name:     "start1",
					Start:    start1,
					Duration: duration},
				scheduler.Span{
					Name:     "start2",
					Start:    start2,
					Duration: duration},
			},
			ValveBox: vb,
		},
	}

	noWaterDays := make(map[time.Weekday]bool)
	noWaterDays[time.Sunday] = false
	noWaterDays[time.Monday] = false
	noWaterDays[time.Tuesday] = false
	noWaterDays[time.Wednesday] = true
	noWaterDays[time.Thursday] = false
	noWaterDays[time.Friday] = true
	noWaterDays[time.Saturday] = false

	return &Config{
		Stations:    stations,
		NoWaterDays: noWaterDays,
	}
}
