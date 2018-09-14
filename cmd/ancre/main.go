package main

import (
	"os"

	"github.com/ulule/ancre/cli"
)

func main() {
	app := cli.New()
	app.Run(os.Args)
}
