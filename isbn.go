// The isbn package inspects and verifies International Standard Book Numbers
// (ISBN) according to the ISO standard
package isbn

import (
	"fmt"
	"strings"
)

type parsed struct {
	body []int32
}

const integerRuneStart = 0x30
const urnPrefix string = "urn:isbn:"
const isbn10Prefix13 string = "\x09\x07\x08"

// sanitizeRune is a map function to be used with strings.Map that strips all
// non-ISBN characters from a string and returns int32 values of digits
func sanitizeRune(r rune) rune {
	if r >= '0' && r <= '9' {
		// return the integer value of the rune by subtracting the number of utf-8
		// runes before the first integer rune
		return r - integerRuneStart
	}
	if r == 'X' || r == 'x' {
		// in base11 checksums used for SBNs and ISB-10s 'X' is used to substitute
		// 10.
		return 10
	}
	return -1
}

// convertDigitsToString reverts sanitizeRune by converting int32 values to
// their utf-8 representation as string
func convertDigitsToString(i []int32) (s string) {
	s = strings.Map(func(r rune) rune {
		return r + integerRuneStart
	}, string(i))
	return s
}

// parse parses a string to an isbn
func parse(s string) (parsed, error) {
	s = strings.TrimPrefix(s, urnPrefix)

	// maximum length of an isbn13 is 13 characters + 4 hyphens
	if len(s) > 13+4 {
		return parsed{}, fmt.Errorf("isbn: parse error: too long")
	}

	runes := strings.Map(sanitizeRune, s)
	switch len(runes) {
	case 9:
		// this is an SBN
		return parseSbn(runes)
	case 10:
		// this is an isbn10
		return parse10(runes)
	case 13:
		// this is an isbn13
		return parse13(runes)
	}

	return parsed{}, fmt.Errorf("isbn: parse: invalid length %v", len(runes))
}

// parseSbn parses a slice of 9 integers by interpreting them as ISBN-10
func parseSbn(s string) (parsed, error) {
	return parse10("\x00" + s)
}

// parse10 parses a slice of 10 integers by calculating the check digit.
func parse10(s string) (parsed, error) {
	runes := []rune(s)
	if check10(runes) != runes[len(runes)-1] {
		return parsed{}, fmt.Errorf("isbn: invalid isbn-10 checksum")
	}
	return parsed{
		body: append([]int32{9, 7, 8}, runes...),
	}, nil
}

// parse13 parses a slice of 13 integers by verifying they begin with a valid
// isbn prefix (978 or 979) and calculating the check digit.
//
// currently this function does only verify if an isbn adheres to the isbn
// format, but not if it is actually allocated by the international isbn agency
func parse13(s string) (parsed, error) {
	runes := []rune(s)
	if string(runes[:3]) != "\x09\x07\x08" && string(runes[:3]) != "\x09\x07\x09" {
		return parsed{}, fmt.Errorf("isbn: invalid isbn-13 gs1")
	}
	if check13(runes) != runes[len(runes)-1] {
		return parsed{}, fmt.Errorf("isbn: invalid isbn-13 checksum")
	}
	return parsed{
		body: runes,
	}, nil
}

// check13 calculates the check digit for an ISBN-13 by multiplying every digit
// with a weight, adding them together so that the sum of all digits including
// the check is a multiple of 10.
// If 13 digits are passed in the input slice, the last digit will be discarded
// in favour of the new check digit.
func check13(i []int32) int32 {
	if len(i) == 13 {
		i = i[:12]
	}
	var check int32
	for index, number := range i {
		if index%2 == 0 {
			check = check + number
			continue
		}
		check = check + number*3
	}
	return (10 - check%10) % 10
}

// check10 calculates the check digit for an ISBN-10 by multiplying every digit
// with its weight and adding all digits together so that the sum of all digits
// including the check is a multiple of eleven.
func check10(i []int32) int32 {
	if len(i) == 10 {
		i = i[:9]
	}
	var check int32
	for index, number := range i {
		check += number * int32(10-index)
	}
	return (11 - check%11) % 11
}

func isbn13(p parsed) string {
	body := convertDigitsToString(p.body)
	check := rune(check13(p.body))
	return fmt.Sprintf("%s%c", body, check)
}

func isbn10(p parsed) string {
	body := convertDigitsToString(p.body[3:])
	check := rune(check13(p.body))
	return fmt.Sprintf("%s%c", body, check)
}

// SBN takes a valid ISBN-13 or ISBN-10 and returns the corresponding British
// Standard Book Number (SBN) which is nine digits and two hyphens long. An SBN
// only exists, if the ISBN group element is zero.
func SBN(s string) (string, error) {
	p, err := parse(s)
	if err != nil {
		return "", err
	}
	if string(p.body[3:4]) != "\x00" {
		// cannot interpret ISBN as SBN because SBN depends on having the same
		// checksum as the equivalent ISBN-10 - which is only possible if the ISBN
		// group part is '0'
		return "", fmt.Errorf("isbn: sbn: group is not 0")
	}
	body := convertDigitsToString(p.body[4:])
	check := rune(check13(p.body))
	return fmt.Sprintf("%s%c", body, check), nil
}

// ISBN10 takes a valid ISBN-13 or ISBN-10 and returns the corresponding
// ISBN-10 which is ten runes and three hyphens long.
func ISBN10(s string) (string, error) {
	p, err := parse(s)
	if err != nil {
		return "", err
	}
	if string(p.body[:3]) != "\x09\x07\x08" {
		// cannot convert ISBN-13 to ISBN-10 because only ISBNs with 978 prefix can
		// be interpreted as ISBN-10
		return "", fmt.Errorf("isbn: isbn-10: gs1 is not 978")
	}
	return isbn10(p), nil
}

// ISBN13 takes a valid ISBN-13 or ISBN-10 and returns the corresponding
// ISBN-13 which is thirteen runes and four hyphens long.
func ISBN13(s string) (string, error) {
	p, err := parse(s)
	if err != nil {
		return "", err
	}
	return isbn13(p), nil
}
