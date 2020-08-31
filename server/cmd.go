package server

import (
	"log"

	"okumoto.net/cliutil"
	"okumoto.net/controller"
	"okumoto.net/valvebox"
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

	//==================================================
	// Setup Controller
	//==================================================
	c := controller.Controller{}

	c.MainLoop()

	return cliutil.ExitOk
}
