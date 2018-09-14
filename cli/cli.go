package cli

import (
	ucli "gopkg.in/urfave/cli.v2"

	"github.com/ulule/ancre/cli/cmd"
)

// New creates a new cli app.
func New() *ucli.App {
	return &ucli.App{
		Name:  "ancre",
		Usage: "anchoring data in time",
		Flags: []ucli.Flag{
			&ucli.StringFlag{
				Name:  "D",
				Usage: "calendar host",
			},
		},
		Commands: []*ucli.Command{
			{
				Name:   "info",
				Usage:  "display timestamp information of a given .ots file",
				Action: info,
			},
			{
				Name:   "stamp",
				Usage:  "stamp the given file",
				Action: stamp,
			},
		},
	}

}

func info(c *ucli.Context) error {
	return cmd.Info(c.Args().First())
}

func stamp(c *ucli.Context) error {
	return cmd.Stamp(c.Args().First(), c.String("D"))
}
