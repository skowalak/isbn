package isbn

import (
	"fmt"
	"strings"
)

type parsed struct {
	gs1 []int32
	// registration group element is a one to five digit number that is valid
	// within a single gs1 element
	group []int32
	// registrant is a variable-length code of up to seven digits
	registrant []int32
	// publication is a variable-length code of up to six digits
	publication []int32

	// if group, registrant and publication could not be parsed, data will be
	// contained in body.
	body []int
}

const integerRuneStart = 0x30
const urnPrefix string = "urn:isbn:"
const isbn10Prefix13 string = "\x09\x07\x08"

func Parse(s string) (string, error) {
	p, err := parse(s)
	if err != nil {
		return "", err
	}

	return isbn13(p), nil
}

func sanitizeRune(r rune) rune {
	if r >= 0 && r < 10 {
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

func convertDigitsToString(i []int32) string {
	var out strings.Builder
	for _, rune := range i {
		out.WriteRune(rune + integerRuneStart)
	}
	return out.String()
}

func convertString(r rune) rune {
	return r + integerRuneStart
}

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

	return parsed{}, fmt.Errorf("isbn: parse: invalid length")
}

func parseSbn(s string) (parsed, error) {
	return parse10("\x00" + s)
}

func parse10(s string) (parsed, error) {
	runes := []rune(s)
	if check10(runes) != runes[len(runes)-1] {
		return parsed{}, fmt.Errorf("isbn: invalid isbn-10 checksum")
	}
	return parsed{
		//TODO
	}, nil
}

func parse13(s string) (parsed, error) {
	runes := []rune(s)
	// if string(runes[:3]) == "\x09\x07\x08" {
	// // get isbn registry data for 978-prefixed isbns
	// } else if string(runes[:3]) == \"x09\x07\x09" {
	// // get isbn registry data for 979-prefixed isbns
	// }

	// currently this lib does only verify if an isbn adheres to the isbn format,
	// but not if it is actually allocated by the international isbn agency
	// https://www.isbn-international.org/
	if string(runes[:3]) != "\x09\x07\x08" && string(runes[:3]) != "\x09\x07\x09" {
		return parsed{}, fmt.Errorf("isbn: invalid isbn-13 gs1")
	}
	return parsed{}, nil
}

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
	gs1 := convertDigitsToString(p.gs1)
	group := convertDigitsToString(p.group)
	reg := convertDigitsToString(p.registrant)
	pub := convertDigitsToString(p.publication)
	check := rune(check13(nil))
	return fmt.Sprintf("%s-%s-%s-%s-%c", gs1, group, reg, pub, check)
}

func isbn10(p parsed) string {
	group := convertDigitsToString(p.group)
	reg := convertDigitsToString(p.registrant)
	pub := convertDigitsToString(p.publication)
	//return fmt.Sprintf("%s-%s-%s-%c", group, reg, pub, check10(p))
	return group + reg + pub
}

// SBN takes a valid ISBN-13 or ISBN-10 and returns the corresponding British
// Standard Book Number (SBN) which is nine digits and two hyphens long. An SBN
// only exists, if the ISBN group element is zero.
func SBN(s string) (string, error) {
	p, err := parse(s)
	if err != nil {
		return "", err
	}
	if string(p.group) != "\x00" {
		return "", fmt.Errorf("isbn: sbn: group is not 0")
	}
	reg := convertDigitsToString(p.registrant)
	pub := convertDigitsToString(p.publication)
	//return fmt.Sprintf("%s-%s-%c", reg, pub, rune(check10(p))), nil
	return reg + pub, nil
}

// ISBN10 takes a valid ISBN-13 or ISBN-10 and returns the corresponding
// ISBN-10 which is ten runes and three hyphens long.
func ISBN10(s string) (string, error) {
	p, err := parse(s)
	if err != nil {
		return "", err
	}
	if string(p.gs1) != "\x09\x07\x08" {
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
