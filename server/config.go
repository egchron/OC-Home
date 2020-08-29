package server

import (
	"log"
	"time"

	"okumoto.net/db"
)

//
// HACK.  For now hard code the configuration.
//

func maybeSetupConfig(database db.DB) {
	// For now hard code the data used to populate the DB.
	stationNames := []string{
		"GrassHouse",
		"GrassFence",
		"PlanterBoxes",
		"Drip",
	}

	for _, name := range stationNames {
		s := db.StationEntry{
			Name: name,
		}
		duplicate, err := database.AddStation(&s)
		if err != nil {
			log.Fatal(err)
		}
		if duplicate {
			log.Printf("We saw %s again\n", name)
		}
	}

	relayNames := []string{
		"GrassHouse",
		"GrassFence",
		"PlanterBoxes",
		"Drip",
	}
	for _, name := range relayNames {
		r := db.RelayEntry{
			Name: name,
		}
		duplicate, err := database.AddRelay(&r)
		if err != nil {
			log.Fatal(err)
		}
		if duplicate {
			log.Printf("We saw %s again\n", name)
		}
	}

	links := []struct {
		station string
		relay   string
	}{
		{"GrassHouse", "GrassHouse"},
		{"GrassFence", "GrassFence"},
		{"PlanterBoxes", "PlanterBoxes"},
		{"Drip", "Drip"},
	}
	for _, v := range links {
		station, err := database.GetStationByName(v.station)
		if err != nil {
			log.Fatal(err)
		}
		relay, err := database.GetRelayByName(v.relay)
		if err != nil {
			log.Fatal(err)
		}

		r := db.SRLinkEntry{
			Station: station.Id,
			Relay:   relay.Id,
		}
		duplicate, err := database.AddSRLink(&r)
		if err != nil {
			log.Fatal(err)
		}
		if duplicate {
			log.Printf("We saw %v again\n", v)
		}
	}

	sunday := time.Date(2017, 1, 1, 0, 0, 0, 0, time.Local)
	//monday := time.Date(2017, 1, 2, 0, 0, 0, 0, time.Local)
	//tuesday := time.Date(2017, 1, 3, 0, 0, 0, 0, time.Local)
	//wednesday := time.Date(2017, 1, 4, 0, 0, 0, 0, time.Local)
	//thursday := time.Date(2017, 1, 5, 0, 0, 0, 0, time.Local)
	//friday := time.Date(2017, 1, 6, 0, 0, 0, 0, time.Local)
	//saturday := time.Date(2017, 1, 7, 0, 0, 0, 0, time.Local)

	weeks := db.Weeks
	var one int32 = 1
	excludes := []struct {
		desc     string
		station  string
		start    time.Time
		end      *time.Time
		duration time.Duration
		interval *int32
		units    *db.Units
	}{
		{"Sunday",
			"GrassHouse",
			sunday, nil, 24 * time.Hour, &one, &weeks,
		},
		{"GrassFence",
			"GrassFence",
			sunday, nil, 24 * time.Hour, &one, &weeks,
		},
		{"PlanterBoxes",
			"PlanterBoxes",
			sunday, nil, 24 * time.Hour, &one, &weeks,
		},
		{"Drip",
			"Drip",
			sunday, nil, 24 * time.Hour, &one, &weeks,
		},
	}
	for _, v := range excludes {
		r := db.ExcludeEntry{
			Desc:     v.desc,
			Start:    v.start,
			End:      v.end,
			Duration: v.duration,
			Interval: v.interval,
			Units:    v.units,
		}
		duplicate, err := database.AddExclude(&r)
		if err != nil {
			log.Fatal(err)
		}
		if duplicate {
			log.Printf("We saw %v again\n", v)
		}
	}
}
