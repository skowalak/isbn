package isbn

import (
	"reflect"
	"testing"
)

func TestIsbn10Checksum(t *testing.T) {
	// example from wikipedia
	isbnData := []int32{0, 3, 0, 6, 4, 0, 6, 1, 5}
	var want int32 = 2

	got := check10(isbnData)

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestIsbn13Checksum(t *testing.T) {
	// example from wikipedia
	isbnData := []int32{9, 7, 8, 0, 3, 0, 6, 4, 0, 6, 1, 5}
	var want int32 = 7

	got := check13(isbnData)

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestIsbn10Parser(t *testing.T) {
	isbn10s := []struct {
		input string
		data  []int32
		valid bool
	}{
		{
			input: "0672323567",
			data:  []int32{9, 7, 8, 0, 6, 7, 2, 3, 2, 3, 5, 6, 7},
			// normal isbn10 without hyphens
			valid: true,
		},
		{
			input: "1-316-87371-4",
			data:  []int32{9, 7, 8, 1, 3, 1, 6, 8, 7, 3, 7, 1, 4},
			// isbn10 with hyphens
			valid: true,
		},
		{
			input: "urn:isbn:0789702363",
			data:  []int32{9, 7, 8, 0, 7, 8, 9, 7, 0, 2, 3, 6, 3},
			// isbn10 URN
			valid: true,
		},
		{
			input: "059610183X",
			data:  []int32{9, 7, 8, 0, 5, 9, 6, 1, 0, 1, 8, 3, 10},
			// with an x as checksum
			valid: true,
		},
		{
			input: "05961018310",
			data:  []int32{9, 7, 8, 0, 5, 9, 6, 1, 0, 1, 8, 3, 10},
			// too long
			valid: false,
		},
		{
			input: "067232356",
			data:  []int32{9, 7, 8, 0, 6, 7, 2, 3, 2, 3, 5, 6, 7},
			// too short
			valid: false,
		},
		{
			input: "067232357",
			data:  []int32{9, 7, 8, 0, 6, 7, 2, 3, 2, 3, 5, 6, 7},
			// invalid checksum
			valid: false,
		},
	}
	for _, tt := range isbn10s {
		got, err := parse(tt.input)
		if tt.valid {
			if err != nil {
				t.Errorf("got error for valid isbn %v: %v", tt.input, err)
			}
			if !reflect.DeepEqual(got.body, tt.data) {
				t.Errorf("got %v want %v", got.body, tt.data)
			}
		} else {
			// invalid
			if err == nil {
				t.Errorf("got no error for invalid isbn %v", tt.input)
			}
		}
	}
}
func TestIsbn13Parser(t *testing.T) {
	// for anyone reading this, i just pulled the newest books from nyt
	// bestsellers
	isbn13s := []struct {
		input string
		data  []int32
		valid bool
	}{
		{
			input: "9781250289551",
			data:  []int32{9, 7, 8, 1, 2, 5, 0, 2, 8, 9, 5, 5, 1},
			// isbn13 without hyphens
			valid: true,
		},
		{
			input: "978-1-64937-404-2",
			data:  []int32{9, 7, 8, 1, 6, 4, 9, 3, 7, 4, 0, 4, 2},
			// isbn13 with hyphens
			valid: true,
		},
		{
			input: "urn:isbn:9780593422946",
			data:  []int32{9, 7, 8, 0, 5, 9, 3, 4, 2, 2, 9, 4, 6},
			// isbn10 URN
			valid: true,
		},
		// TODO: get a valid 979 prefix isbn somewhere
		// {
		// 	input: "059610183X",
		// 	data:  []int32{9, 7, 8, 0, 5, 9, 6, 1, 0, 1, 8, 3, 10},
		// 	// with an x as checksum
		// 	valid: true,
		// },
		{
			input: "97805934929189",
			data:  []int32{9, 7, 8, 0, 5, 9, 3, 4, 9, 2, 9, 1, 8, 9},
			// too long
			valid: false,
		},
		{
			input: "978059349291",
			data:  []int32{9, 7, 8, 0, 5, 9, 3, 4, 9, 2, 9, 1},
			// too short
			valid: false,
		},
		{
			input: "9780593492910",
			data:  []int32{9, 7, 8, 0, 5, 9, 3, 4, 9, 2, 9, 1, 0},
			// invalid checksum
			valid: false,
		},
		{
			input: "978O593236598",
			data:  []int32{9, 7, 8, 0, 5, 9, 3, 2, 3, 6, 5, 9, 8},
			// invalid characters
			valid: false,
		},
	}
	for _, tt := range isbn13s {
		got, err := parse(tt.input)
		if tt.valid {
			if err != nil {
				t.Errorf("got error for valid isbn %v: %v", tt.input, err)
			}
			if !reflect.DeepEqual(got.body, tt.data) {
				t.Errorf("got %v want %v", got.body, tt.data)
			}
		} else {
			// invalid
			if err == nil {
				t.Errorf("got no error for invalid isbn %v", tt.input)
			}
		}
	}
}
