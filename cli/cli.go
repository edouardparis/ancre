package cli

import (
	ucli "gopkg.in/urfave/cli.v2"

	"github.com/edouardparis/ancre/cli/cmd"
	"github.com/edouardparis/ancre/logging"
)

// New creates a new cli app.
func New() *ucli.App {
	ucli.VersionFlag = &ucli.BoolFlag{
		Name: "version", Aliases: []string{},
		Usage: "print the version",
	}
	return &ucli.App{
		Name:                  "ancre",
		Usage:                 "anchoring data in time",
		Copyright:             "MIT License Copyright (c) 2018 EdouardParis",
		EnableShellCompletion: true,
		Flags: []ucli.Flag{
			&ucli.BoolFlag{
				Name: "verbose", Aliases: []string{"v"},
				Usage: "print log output",
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
				Flags: []ucli.Flag{
					&ucli.StringSliceFlag{
						Name:  "c",
						Usage: "list of calendar",
					},
					&ucli.StringFlag{
						Name:  "o",
						Value: "",
						Usage: "output file, default is ./file.ots",
					},
				},
			},
		},
	}

}

func setUpLogger(c *ucli.Context) (logging.Logger, error) {
	return logging.NewCliLogger(&logging.Config{
		Verbose: c.Bool("verbose"),
	})
}

func info(c *ucli.Context) error {
	logger, err := setUpLogger(c)
	if err != nil {
		return err
	}

	err = cmd.Info(logger, c.Args().First())
	if err != nil {
		return ucli.Exit(err, 1)
	}
	return nil
}

func stamp(c *ucli.Context) error {
	logger, err := setUpLogger(c)
	if err != nil {
		return err
	}

	calendars := c.StringSlice("c")
	if len(calendars) == 0 {
		calendars = []string{
			"https://a.pool.opentimestamps.org",
			"https://b.pool.opentimestamps.org",
			"https://a.pool.eternitywall.com",
			"https://ots.btc.catallaxy.com",
		}
	}

	err = cmd.Stamp(logger, c.Args().First(),
		c.String("o"), calendars)
	if err != nil {
		return ucli.Exit(err, 1)
	}
	return nil
}
