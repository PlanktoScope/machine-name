package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/PlanktoScope/machine-name/pkg/haikunator"
	"github.com/PlanktoScope/machine-name/pkg/wordlists"
)

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var app = &cli.App{
	Name: "machine-name",
	// TODO: see if there's a way to get the version from a build tag, so that we don't have to update
	// this manually
	Version: "v0.1.2",
	Usage:   "Generates localized Heroku-style names from 32-bit serial numbers",
	Commands: []*cli.Command{
		nameCmd,
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "lang",
			Aliases: []string{"language"},
			Usage:   "Locale for names",
			EnvVars: []string{"LANG"},
		},
	},
	Suggest: true,
}

var nameCmd = &cli.Command{
	Name:  "name",
	Usage: "Generates a name from a serial number",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "sn",
			Aliases: []string{"serial-number"},
			Usage:   "32-bit serial number for generating the machine name",
		},
		&cli.StringFlag{
			Name:    "format",
			Aliases: []string{"f"},
			Usage:   "Serial number for generating the machine name",
			Value:   "hex",
		},
	},
	Action: func(c *cli.Context) error {
		raw := c.String("sn")
		format := c.String("format")
		parsed, err := parseSerialNumber(raw, format)
		if err != nil {
			return errors.Wrapf(err, "couldn't parse '%s'-formatted serial number '%s'", format, raw)
		}
		first, second, err := wordlists.Load(wordlists.FS, c.String("lang"))
		if err != nil {
			return errors.Wrapf(err, "couldn't load wordlists for language '%s", c.String("lang"))
		}
		fmt.Println(haikunator.SelectName(parsed, first, second))
		return nil
	},
}

func parseSerialNumber(raw string, format string) (uint32, error) {
	if format != "hex" {
		return 0, errors.Errorf("unrecognized serial number format '%s'", format)
	}
	const base = 16
	const parsedWidth = 32
	parsed64, err := strconv.ParseUint(strings.TrimPrefix(raw, "0x"), base, parsedWidth)
	return uint32(parsed64), err
}
