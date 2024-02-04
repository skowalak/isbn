package isbn

import "testing"

func TestIsbn10Checksum(t *testing.T) {
	// example from
	isbnData := []int32{0, 3, 0, 6, 4, 0, 6, 1, 5}
	var want int32 = 2

	got := check10(isbnData)

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestIsbn13Checksum(t *testing.T) {
	// example from
	isbnData := []int32{9, 7, 8, 0, 3, 0, 6, 4, 0, 6, 1, 5}
	var want int32 = 7

	got := check13(isbnData)

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
