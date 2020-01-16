package main

import (
	"os"

	"okumoto.net/server"
)

func main() {
	os.Exit(server.CmdMain(os.Args))
}
