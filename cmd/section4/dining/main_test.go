package main

import (
	"testing"
	"time"
)

func Test_dine(t *testing.T) {
	var tests = []struct {
		name  string
		delay time.Duration
	}{
		{"zero delay", 0 * time.Second},
		{"quarter second delay", 250 * time.Millisecond},
		{"half second delay", 500 * time.Millisecond},
	}

	for _, test := range tests {
		eatingTime = test.delay
		thinkingTime = test.delay
		dine()

		if len(results.finished) != len(results.finished) {
			t.Errorf(
				"%q: incorrect length of slice; expected %d, but got %d",
				test.name,
				len(philosophers),
				len(results.finished),
			)
		}
	}
}
