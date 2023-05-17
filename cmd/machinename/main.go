package main

import (
	"fmt"
	"io/fs"
	"log"
	"math/bits"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/PlanktoScope/machine-name/internal/wordlists/generated"
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
	Version: "v0.1.0",
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
		first, second, err := loadWordlists(c.String("lang"))
		if err != nil {
			return errors.Wrapf(err, "couldn't load wordlists for language '%s", c.String("lang"))
		}
		fmt.Println(computeName(parsed, first, second))
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

func loadWordlists(lang string) (first []string, second []string, err error) {
	wfs := generated.WordlistsFS
	lfs, err := fs.Sub(wfs, lang)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "couldn't open subdirectory '%s'", lang)
	}
	if first, err = loadWordlist(lfs, "first.txt"); err != nil {
		return nil, nil, errors.Wrapf(err, "couldn't load '%s/first.txt'", lang)
	}
	if second, err = loadWordlist(lfs, "second.txt"); err != nil {
		return nil, nil, errors.Wrapf(err, "couldn't load '%s/second.txt'", lang)
	}
	return first, second, nil
}

func loadWordlist(lang fs.FS, file string) ([]string, error) {
	contents, err := fs.ReadFile(lang, file)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't read file '%s'", file)
	}
	fileWords := strings.Split(string(contents), "\n")
	words := make([]string, 0, len(fileWords))
	for _, word := range fileWords {
		if word == "" {
			continue
		}
		if strings.HasPrefix(word, "# ") {
			continue
		}
		words = append(words, word)
	}
	return words, nil
}

func computeName(sn uint32, first []string, second []string) string {
	firstSize := uint32(len(first))
	secondSize := uint32(len(second))
	quotient := shuffle(sn)
	firstIndex := quotient % firstSize
	quotient /= firstSize
	secondIndex := quotient % secondSize
	quotient /= secondSize
	return fmt.Sprintf("%s-%s-%d\n", first[firstIndex], second[secondIndex], quotient)
}

const (
	shuffleShift8 = 8
	shuffleShift4 = 4
	shuffleShift2 = 2
	shuffleShift1 = 1
	shuffleMask8  = 0x0000ff00
	shuffleMask4  = 0x00f000f0
	shuffleMask2  = 0x0c0c0c0c
	shuffleMask1  = 0x22222222
)

// shuffle performs a one-to-one mapping of the serial number so that consecutive numbers are no
// longer close to each other.
func shuffle(x uint32) uint32 {
	// This code was copied from the Hacker's Delight website at
	// https://web.archive.org/web/20160405214331/http://hackersdelight.org/hdcodetxt/shuffle.c.txt
	// which is licensed released to the public domain - for details, refer to
	// https://web.archive.org/web/20160309224818/http://www.hackersdelight.org/permissions.htm
	t := (x ^ (x >> shuffleShift8)) & shuffleMask8
	x = x ^ t ^ (t << shuffleShift8)
	t = (x ^ (x >> shuffleShift4)) & shuffleMask4
	x = x ^ t ^ (t << shuffleShift4)
	t = (x ^ (x >> shuffleShift2)) & shuffleMask2
	x = x ^ t ^ (t << shuffleShift2)
	t = (x ^ (x >> shuffleShift1)) & shuffleMask1
	x = x ^ t ^ (t << shuffleShift1)
	x = bits.Reverse32(x)
	return x
}

/*
// unshuffle inverts the shuffle operation.
func unshuffle(x uint32) uint32 {
	// This code was copied from the Hacker's Delight website at
	// https://web.archive.org/web/20160405214331/http://hackersdelight.org/hdcodetxt/shuffle.c.txt
	// which is licensed released to the public domain - for details, refer to
	// https://web.archive.org/web/20160309224818/http://www.hackersdelight.org/permissions.htm
	x = bits.Reverse32(x)
	t := (x ^ (x >> shuffleShift1)) & shuffleMask1
	x = x ^ t ^ (t << shuffleShift1)
	t = (x ^ (x >> shuffleShift2)) & shuffleMask2
	x = x ^ t ^ (t << shuffleShift2)
	t = (x ^ (x >> shuffleShift4)) & shuffleMask4
	x = x ^ t ^ (t << shuffleShift4)
	t = (x ^ (x >> shuffleShift8)) & shuffleMask8
	x = x ^ t ^ (t << shuffleShift8)
	return x
}
*/
