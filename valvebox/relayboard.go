package valvebox

type RelayBoard interface {
	Open(devName string) error
	Relay(id int, on bool) error
}
