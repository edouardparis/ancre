package main

import (
	"os"

	"github.com/edouardparis/ancre/cli"
)

func main() {
	app := cli.New()
	app.Run(os.Args)
}
