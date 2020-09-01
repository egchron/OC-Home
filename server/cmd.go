package server

import (
	"log"
	"time"

	"okumoto.net/cliutil"
	"okumoto.net/controller"
	"okumoto.net/valvebox"

	"github.com/go-co-op/gocron"
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
		log.Fatalf("cli: %+v", err)
	}

	//==================================================
	// Setup Hardware
	//==================================================
	vb := valvebox.New("backyard")

	rb := &valvebox.Numato{}
	if err := rb.Open(cli.DevName); err != nil {
		log.Fatal(err)
	}

	// The stuff here should be read in from some sort of
	// configuration file.
	relayNames := []string{
		"GrassHouse",
		"GrassFence",
		"PlanterBoxes",
		"Drip",
	}
	for i, name := range relayNames {
		vb.AddRelay(name, rb, i)
	}

	stationNames := []string{
		"GrassHouse",
		"GrassFence",
		"PlanterBoxes",
		"Drip",
	}
	for _, name := range stationNames {
		relayName := name
		vb.AddStation(name, []string{relayName}, time.Second)
	}

	//==================================================
	// Setup Controller
	//==================================================
	c := controller.Controller{
		ValveBox:  vb,
		Scheduler: gocron.NewScheduler(time.Local),
	}

	c.MainLoop()

	return cliutil.ExitOk
}
