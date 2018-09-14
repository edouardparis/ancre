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
		Commands: []*ucli.Command{
			{
				Name:  "info",
				Usage: "display timestamp information of a given .ots file",
				Action: func(c *ucli.Context) error {
					return cmd.Info(c.Args().First())
				},
			},
			{
				Name:  "stamp",
				Usage: "stamp the given file",
				Action: func(c *ucli.Context) error {
					return cmd.Stamp(c.Args().First())
				},
			},
		},
	}

}
