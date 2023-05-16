package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
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
	Version: "v0.0.1",
	Usage:   "Generates localized Heroku-style names from serial numbers",
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

const uint32max = 4294967295

var nameCmd = &cli.Command{
	Name:  "name",
	Usage: "Generates a name",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "sn",
			Aliases: []string{"serial-number"},
			Usage:   "Serial number for generating the machine name",
		},
		&cli.StringFlag{
			Name:    "format",
			Aliases: []string{"f"},
			Usage:   "Serial number for generating the machine name",
			Value:   "hex",
		},
		&cli.Uint64Flag{
			Name:        "max",
			Usage:       "Upper limit (inclusive) of the range of valid serial numbers",
			Value:       uint32max,
			DefaultText: "2^32 - 1",
		},
	},
	Action: func(c *cli.Context) error {
		fmt.Println(c.String("lang"))
		raw := c.String("sn")
		format := c.String("format")
		parsed, err := parseSerialNumber(raw, format)
		if err != nil {
			return errors.Wrapf(err, "couldn't parse '%s'-formatted serial number '%s'", format, raw)
		}
		max := c.Uint64("max")
		if parsed > max {
			return errors.Errorf(
				"serial number '%x' is greater than maximum allowed number '%x'", parsed, max,
			)
		}
		fmt.Printf("0x%x (%d) 0x%x\n", parsed, parsed, max)
		return nil
	},
}

func parseSerialNumber(raw string, format string) (uint64, error) {
	if format != "hex" {
		return 0, errors.Errorf("unrecognized serial number format '%s'", format)
	}
	const base = 16
	const parsedWidth = 64
	return strconv.ParseUint(strings.TrimPrefix(raw, "0x"), base, parsedWidth)
}
