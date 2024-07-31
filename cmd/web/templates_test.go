package main

import (
	"testing"
	"time"

	"github.com/rynhndrcksn/portfolio/internal/assert"
)

func TestHumanDate(t *testing.T) {
	// Use an anonymous struct to create all of our test cases.
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 01, 1, 12, 0, 0, 0, time.UTC),
			want: "01 Jan 2024 at 12:00",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "PST",
			tm:   time.Date(2024, 01, 1, 12, 0, 0, 0, time.FixedZone("UTC-8", 8*60*60)),
			want: "01 Jan 2024 at 04:00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			assert.Equal(t, hd, tt.want)
		})
	}
}
