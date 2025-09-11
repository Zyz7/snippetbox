package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {

	// Initialize and pass it to the hnumanDate function
	/*tm := time.Date(2024, 3, 17, 10, 15, 0, 0, time.UTC)
	hd := humanDate(tm)

	if hd != "17 Mar 2024 at 10:15" {
		t.Errorf("go %q; want %q", hd, "17 Mar 2024 at 10:15")
	}*/

	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2024 at 10:15",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			hd := humanDate(tt.tm)
			if hd != tt.want {
				t.Errorf("go %q; want %q", hd, tt.want)
			}
		})
	}
}
