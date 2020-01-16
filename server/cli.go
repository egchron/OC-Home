package server

import (
	"flag"

	"github.com/pkg/errors"
)

type Cli struct {
	DevName string
}

func (cli *Cli) ParseArgsEnv(args []string) error {

	flag.StringVar(&cli.DevName, "device", "/dev/tty.usbmodem2301",
		"Relay board device.")

	if err := flag.CommandLine.Parse(args[1:]); err != nil {
		errors.Wrap(err, "ParseArgsEnv")
	}

	return nil
}
