package main

import (
	"testing"
	"time"

	"github.com/Divyue30597/snippetbox-lets-go/internal/assert"
)

func TestHumanDate(t *testing.T) {
	// tm := time.Date(2023, 2, 18, 12, 9, 0, 0, time.UTC)
	// hd := humanDate(tm)

	// if hd != "18 Feb 2023 at 12:09" {
	// 	t.Errorf("got %q, want %q", hd, "18 Feb 2023 at 12:09")
	// }

	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2023, 2, 18, 12, 9, 0, 0, time.UTC),
			want: "18 Feb 2023 at 12:09",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2023, 2, 18, 12, 9, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "18 Feb 2023 at 11:09",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			// if hd != tt.want {
			// 	t.Errorf("got %q, want %q", hd, tt.want)
			// }

			assert.Equal(t, hd, tt.want)
		})
	}
}
