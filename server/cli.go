package server

import (
	"okumoto.net/db"

	"github.com/jnovack/flag"
	"github.com/pkg/errors"
)

type Cli struct {
	DevName string
	DB      db.DBConfig
}

func (cli *Cli) ParseArgsEnv(args []string) error {

	flag.StringVar(&cli.DevName, "device", "/dev/tty.usbmodem2301",
		"Relay board device.")

	flag.StringVar(&cli.DB.User, "mysql-user", "root",
		"DB user.")

	flag.StringVar(&cli.DB.Pass, "mysql-pwd", "my-secrit-pw",
		"DB passwd.")

	flag.StringVar(&cli.DB.Server, "mysql-host", "localhost",
		"DB host.")

	flag.StringVar(&cli.DB.Port, "mysql-tcp-port", "3306",
		"DB port.")

	flag.StringVar(&cli.DB.Name, "mysql-name", "Water",
		"DB name.")

	if err := flag.CommandLine.Parse(args[1:]); err != nil {
		errors.Wrap(err, "ParseArgsEnv")
	}

	return nil
}
