package db

import (
	"time"
)

// DBConfig is information used to connect to the mysql server.
type DBConfig struct {
	User   string
	Pass   string
	Server string
	Port   string
	Name   string
}

type DB interface {
	AddStation(s *StationEntry) (bool, error)
	GetStations() ([]StationEntry, error)
	GetStationByName(string) (*StationEntry, error)

	AddRelay(s *RelayEntry) (bool, error)
	GetRelays() ([]RelayEntry, error)
	GetRelayByName(string) (*RelayEntry, error)

	AddSRLink(s *SRLinkEntry) (bool, error)
	GetSRLinks() ([]SRLinkEntry, error)

	AddExclude(s *ExcludeEntry) (bool, error)
}

type StationEntry struct {
	Id     int64
	Name   string
	Budget time.Duration
}

type RelayEntry struct {
	Id   int64
	Name string
}

type SRLinkEntry struct {
	Id      int64
	Station int64
	Relay   int64
}

type Units int

const (
	Seconds Units = iota
	Minutes
	Hours
	Days
	Weeks
	Months
)

type ExcludeEntry struct {
	Id       int64
	Desc     string
	Start    time.Time
	End      *time.Time
	Duration time.Duration
	Interval *int32
	Units    *Units
}
