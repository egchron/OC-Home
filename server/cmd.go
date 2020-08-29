package server

import (
	"log"
	"time"

	"okumoto.net/cliutil"
	"okumoto.net/controller"
	"okumoto.net/db/mysqldb"
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
	// Setup DB
	//==================================================
	database, err := mysqldb.New(cli.DB)
	if err != nil {
		log.Fatal(err)
	}

	if err := database.MaybeSetupDB("db/mysqldb/schema.sql"); err != nil {
		log.Fatal(err)
	}

	maybeSetupConfig(database)

	//==================================================
	// Setup Hardware
	//==================================================
	vb := valvebox.New("backyard")

	rb := &valvebox.Numato{}
	if err := rb.Open(cli.DevName); err != nil {
		log.Fatal(err)
	}

	relayList, err := database.GetRelays()
	if err != nil {
		log.Fatal(err)
	}
	for i, r := range relayList {
		vb.NewRelay(r.Name, rb, i)
	}

	stationList, err := database.GetStations()
	if err != nil {
		log.Fatal(err)
	}
	for _, s := range stationList {
		vb.NewStation(s.Name, time.Second)
	}

	linkList, err := database.GetSRLinks()
	if err != nil {
		log.Fatal(err)
	}
	for _, sr := range linkList {
		station, err := database.GetStationById(sr.Station)
		if err != nil {
			log.Fatal(err)
		}
		relay, err := database.GetRelayById(sr.Relay)
		if err != nil {
			log.Fatal(err)
		}
		vb.AddRelay(station.Name, relay.Name)
	}

	//==================================================
	// Setup Controller
	//==================================================
	c := controller.Controller{
		ValveBox:  vb,
		Scheduler: gocron.NewScheduler(time.Local),
		Db:        database,
	}

	c.MainLoop()

	return cliutil.ExitOk
}
