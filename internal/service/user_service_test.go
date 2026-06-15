package service

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		dob     string
		now     string
		wantAge int
	}{
		{
			name:    "exact birthday",
			dob:     "10-05-1990",
			now:     "10-05-2025",
			wantAge: 35,
		},
		{
			name:    "one day before birthday",
			dob:     "10-05-1990",
			now:     "09-05-2025",
			wantAge: 34,
		},
		{
			name:    "one day after birthday",
			dob:     "10-05-1990",
			now:     "11-05-2025",
			wantAge: 35,
		},
		{
			name:    "leap year dob (Feb 29) — non-leap year now (Feb 28)",
			dob:     "29-02-2000",
			now:     "28-02-2025",
			wantAge: 24,
		},
		{
			name:    "leap year dob (Feb 29) — non-leap year now (Mar 01)",
			dob:     "29-02-2000",
			now:     "01-03-2025",
			wantAge: 25,
		},
		{
			name:    "newborn (same day)",
			dob:     "01-01-2025",
			now:     "01-01-2025",
			wantAge: 0,
		},
		{
			name:    "year boundary — December 31 dob, January 1 now",
			dob:     "31-12-1990",
			now:     "01-01-2025",
			wantAge: 34,
		},
		{
			name:    "age 100",
			dob:     "15-06-1924",
			now:     "15-06-2024",
			wantAge: 100,
		},
	}

	layout := "02-01-2006"

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dob, err := time.Parse(layout, tc.dob)
			if err != nil {
				t.Fatalf("failed to parse dob %q: %v", tc.dob, err)
			}
			now, err := time.Parse(layout, tc.now)
			if err != nil {
				t.Fatalf("failed to parse now %q: %v", tc.now, err)
			}

			got := CalculateAge(dob, now)
			if got != tc.wantAge {
				t.Errorf("CalculateAge(%q, %q) = %d, want %d",
					tc.dob, tc.now, got, tc.wantAge)
			}
		})
	}
}
