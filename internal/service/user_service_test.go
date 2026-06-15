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
			dob:     "1990-05-10",
			now:     "2025-05-10",
			wantAge: 35,
		},
		{
			name:    "one day before birthday",
			dob:     "1990-05-10",
			now:     "2025-05-09",
			wantAge: 34,
		},
		{
			name:    "one day after birthday",
			dob:     "1990-05-10",
			now:     "2025-05-11",
			wantAge: 35,
		},
		{
			name:    "leap year dob (Feb 29) — non-leap year now (Feb 28)",
			dob:     "2000-02-29",
			now:     "2025-02-28",
			wantAge: 24,
		},
		{
			name:    "leap year dob (Feb 29) — non-leap year now (Mar 01)",
			dob:     "2000-02-29",
			now:     "2025-03-01",
			wantAge: 25,
		},
		{
			name:    "newborn (same day)",
			dob:     "2025-01-01",
			now:     "2025-01-01",
			wantAge: 0,
		},
		{
			name:    "year boundary — December 31 dob, January 1 now",
			dob:     "1990-12-31",
			now:     "2025-01-01",
			wantAge: 34,
		},
		{
			name:    "age 100",
			dob:     "1924-06-15",
			now:     "2024-06-15",
			wantAge: 100,
		},
	}

	layout := "2006-01-02"

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
