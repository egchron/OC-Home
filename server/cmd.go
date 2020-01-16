package server

import (
	"log"
	"time"

	"okumoto.net/cliutil"
	"okumoto.net/config"
	"okumoto.net/controller"
	"okumoto.net/valvebox"

	"github.com/beevik/timerqueue"
)

// Goals:
//	Water lawn on a regular schedule.
//	Adapt to weather conditions.
//		When it rains don't water lawn.
//		curl https://api.weather.gov/points/37.3479,-121.9695
//		curl https://api.weather.gov/gridpoints/MTR/100,104/forecast
//		curl https://api.weather.gov/gridpoints/MTR/100,104/forecast/hourly
//	Record water usage
//	Avoid watering on specified days
//	Limit watering to one station at a time
//	Avoid watering during the day.
//
func CmdMain(args []string) int {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cli := Cli{}
	if err := cli.ParseArgsEnv(args); err != nil {
		log.Fatal(err)
	}

	n := &valvebox.Numato{
		DevName: cli.DevName,
	}

	vb := valvebox.New(n)

	// The stuff here should be read in from some sort of
	// configuration file.
	relayNames := []string{
		"GrassHouse",
		"GrassFence",
		"PlanterBoxes",
		"Drip",
	}
	for i, name := range relayNames {
		vb.Add(i, name)
	}

	config := config.New(vb)

	c := controller.Controller{
		Queue: timerqueue.New(),
	}
	c.Setup(config, time.Now())

	/*
	 * Work starts here.
	 */
	c.MainLoop()

	return cliutil.ExitOk
}
