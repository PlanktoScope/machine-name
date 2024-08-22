// Package haikunator provides language-localized Heroku-style name generation from 32-bit serial
// numbers, with a one-to-one mapping between names and serial numbers and no obvious
// consecutiveness between names for consecutive serial numbers
package haikunator

import (
	"fmt"
	"math/bits"
)

func SelectName(sn uint32, first []string, second []string) string {
	firstSize := uint32(len(first))   //nolint:gosec // we assume string length is less than 2^32
	secondSize := uint32(len(second)) //nolint:gosec // we assume string length is less than 2^32
	quotient := shuffle(sn)
	firstIndex := quotient % firstSize
	quotient /= firstSize
	secondIndex := quotient % secondSize
	quotient /= secondSize
	return fmt.Sprintf("%s-%s-%d", first[firstIndex], second[secondIndex], quotient)
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
