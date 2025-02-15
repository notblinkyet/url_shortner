package generate

import (
	"fmt"
	"testing"
)

func TestGenerate(t *testing.T) {
	tests := []int{5, 10, 15, 20}
	var avail string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_"
	set := make(map[byte]struct{}, 63)
	for i := range avail {
		set[avail[i]] = struct{}{}
	}

	for _, length := range tests {
		t.Run(fmt.Sprintf("Test generate wirh length: %d", length), func(t *testing.T) {
			s, err := Generate(length)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(s) != length {
				t.Errorf("expected length %d, got %d", length, len(s))
			}
			for i := range s {
				if _, ok := set[s[i]]; !ok {
					t.Errorf("invalid character %c in result", s[i])
				}
			}

		})
	}
}
