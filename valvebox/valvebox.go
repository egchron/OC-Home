// Package valvebox is used to control the irrigation valves.
//
// For simple configurations, each station is associated with a single
// relay.
//	24 VAC----------+-------+-------+-------+-------+
//			|	|	|	|	|
//			R0	R1	R2	R3	R4
//			|	|	|	|	|
//	Ground----------S0------S1------S2------S3------S4
//
// Another configuration is a matrix, which allows for more valves with
// the same number if wires.  Each station now has two relays that must
// activate for water to flow.
//
//	24 VAC----------+-------+-------+
//			|       |       |
//			R0	R1	R2
//			|       |       |
//		+--RA---S0------S1------S2
//		|	|       |       |
//	Ground--+--RB---S3------S4------S5
//		|	|       |       |
//		+--RC---S6------S7------S8
//
package valvebox

import (
	"time"
)

// ValveBox represents a logical collection of irrigation valves.
type ValveBox struct {
	name     string
	relays   map[string]*Relay
	stations map[string]*Station
}

// New creates a ValveBox.
func New(name string) *ValveBox {

	m := ValveBox{
		name:     name,
		relays:   make(map[string]*Relay),
		stations: make(map[string]*Station),
	}

	return &m
}

// NewRelay registers a relay.
func (v *ValveBox) NewRelay(name string, board RelayBoard, id int) {
	v.relays[name] = &Relay{
		Name:  name,
		board: board,
		id:    id,
	}
}

// NewStation registers a station and the relays used to activate the
// control valve.
func (v *ValveBox) NewStation(name string, offLatancy time.Duration) {

	v.stations[name] = &Station{
		Name:       name,
		relays:     []*Relay{},
		offLatancy: offLatancy,
	}
}

func (v *ValveBox) AddRelay(station, relay string) {

	v.stations[station].relays = append(
		v.stations[station].relays,
		v.relays[relay])
}

// State changes the state of a station.
func (v *ValveBox) State(stationName string, on bool) error {
	return v.stations[stationName].state(on)
}
