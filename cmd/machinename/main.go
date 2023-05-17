package main

import (
	"fmt"
	"io/fs"
	"log"
	"math/bits"
	"math/rand"
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
	Version: "v0.0.1",
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
		shuffleWordlist(first, 0)
		shuffleWordlist(second, 1)
		firstSize := uint32(len(first))
		secondSize := uint32(len(second))
		quotient := shuffle(parsed)
		firstIndex := quotient % firstSize
		quotient /= firstSize
		secondIndex := quotient % secondSize
		quotient /= secondSize
		fmt.Printf("%s-%s-%d\n", first[firstIndex], second[secondIndex], quotient)
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

func shuffleWordlist(wordlist []string, randSeed int64) {
	rng := rand.New(rand.NewSource(randSeed)) //nolint:gosec // We need a PRNG to shuffle reproducibly
	rng.Shuffle(len(wordlist), func(i, j int) {
		atI := wordlist[i]
		atJ := wordlist[j]
		wordlist[i] = atJ
		wordlist[j] = atI
	})
}

// shuffle performs a one-to-one mapping of the serial number so that consecutive numbers are no
// longer close to each other.
//
//nolint:gomnd // We need a bunch of well-defined bit shifts and masks for this algorithm
func shuffle(x uint32) uint32 {
	// This code was copied from the Hacker's Delight website at
	// https://web.archive.org/web/20160405214331/http://hackersdelight.org/hdcodetxt/shuffle.c.txt
	// which is licensed released to the public domain - for details, refer to
	// https://web.archive.org/web/20160309224818/http://www.hackersdelight.org/permissions.htm
	t := (x ^ (x >> 8)) & 0x0000ff00
	x = x ^ t ^ (t << 8)
	t = (x ^ (x >> 4)) & 0x00f000f0
	x = x ^ t ^ (t << 4)
	t = (x ^ (x >> 2)) & 0x0c0c0c0c
	x = x ^ t ^ (t << 2)
	t = (x ^ (x >> 1)) & 0x22222222
	x = x ^ t ^ (t << 1)
	x = bits.Reverse32(x)
	return x
}

/*
// unshuffle inverts the shuffle operation.
//
//nolint:gomnd // We need a bunch of well-defined bit shifts and masks for this algorithm
func unshuffle(x uint32) uint32 {
	// This code was copied from the Hacker's Delight website at
	// https://web.archive.org/web/20160405214331/http://hackersdelight.org/hdcodetxt/shuffle.c.txt
	// which is licensed released to the public domain - for details, refer to
	// https://web.archive.org/web/20160309224818/http://www.hackersdelight.org/permissions.htm
	x = bits.Reverse32(x)
	t := (x ^ (x >> 1)) & 0x22222222
	x = x ^ t ^ (t << 1)
	t = (x ^ (x >> 2)) & 0x0c0c0c0c
	x = x ^ t ^ (t << 2)
	t = (x ^ (x >> 4)) & 0x00f000f0
	x = x ^ t ^ (t << 4)
	t = (x ^ (x >> 8)) & 0x0000ff00
	x = x ^ t ^ (t << 8)
	return x
}
*/
