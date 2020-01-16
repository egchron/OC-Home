package valvebox

type RelayBoard interface {
	Open() error
	Relay(id int, on bool) error
}
