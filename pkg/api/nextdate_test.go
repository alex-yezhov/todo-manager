package api

import (
	"testing"
	"time"
)

func TestNextDate_Days(t *testing.T) {
	now := mustParseDate(t, "20240126")

	tests := []struct {
		name    string
		start   string
		repeat  string
		want    string
		wantErr bool
	}{
		{
			name:   "every 7 days",
			start:  "20240113",
			repeat: "d 7",
			want:   "20240127",
		},
		{
			name:   "every 20 days",
			start:  "20240120",
			repeat: "d 20",
			want:   "20240209",
		},
		{
			name:   "every 1 day leap year",
			start:  "20240228",
			repeat: "d 1",
			want:   "20240229",
		},
		{
			name:    "interval missing",
			start:   "20240113",
			repeat:  "d",
			wantErr: true,
		},
		{
			name:    "interval too large",
			start:   "20240113",
			repeat:  "d 401",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NextDate(now, tt.start, tt.repeat)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNextDate_Years(t *testing.T) {
	now := mustParseDate(t, "20240126")

	tests := []struct {
		name    string
		start   string
		repeat  string
		want    string
		wantErr bool
	}{
		{
			name:   "simple yearly",
			start:  "20240101",
			repeat: "y",
			want:   "20250101",
		},
		{
			name:   "leap day",
			start:  "20240229",
			repeat: "y",
			want:   "20250301",
		},
		{
			name:   "future date still moves one year forward",
			start:  "20250701",
			repeat: "y",
			want:   "20260701",
		},
		{
			name:    "invalid yearly format",
			start:   "20240101",
			repeat:  "y 2",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NextDate(now, tt.start, tt.repeat)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNextDate_Week(t *testing.T) {
	now := mustParseDate(t, "20240126") // friday

	tests := []struct {
		name    string
		start   string
		repeat  string
		want    string
		wantErr bool
	}{
		{
			name:   "nearest sunday",
			start:  "20240120",
			repeat: "w 7",
			want:   "20240128",
		},
		{
			name:   "monday thursday friday",
			start:  "20240120",
			repeat: "w 1,4,5",
			want:   "20240129",
		},
		{
			name:    "invalid weekday",
			start:   "20240120",
			repeat:  "w 8",
			wantErr: true,
		},
		{
			name:    "missing weekday list",
			start:   "20240120",
			repeat:  "w",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NextDate(now, tt.start, tt.repeat)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNextDate_Month(t *testing.T) {
	now := mustParseDate(t, "20240126")

	tests := []struct {
		name    string
		start   string
		repeat  string
		want    string
		wantErr bool
	}{
		{
			name:   "days in every month",
			start:  "20240116",
			repeat: "m 16,5",
			want:   "20240205",
		},
		{
			name:   "last or 18th day",
			start:  "20240201",
			repeat: "m -1,18",
			want:   "20240218",
		},
		{
			name:   "specific months",
			start:  "20240102",
			repeat: "m 3 1,3,6",
			want:   "20240303",
		},
		{
			name:    "invalid day of month",
			start:   "20240101",
			repeat:  "m 40",
			wantErr: true,
		},
		{
			name:    "invalid month",
			start:   "20240101",
			repeat:  "m 5 13",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NextDate(now, tt.start, tt.repeat)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNextDate_InvalidInput(t *testing.T) {
	now := mustParseDate(t, "20240126")

	tests := []struct {
		name    string
		start   string
		repeat  string
		wantErr bool
	}{
		{
			name:    "empty repeat",
			start:   "20240126",
			repeat:  "",
			wantErr: true,
		},
		{
			name:    "bad start date",
			start:   "20240199",
			repeat:  "y",
			wantErr: true,
		},
		{
			name:    "unsupported format",
			start:   "20240126",
			repeat:  "ooops",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NextDate(now, tt.start, tt.repeat)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
	}
}

func mustParseDate(t *testing.T, s string) time.Time {
	t.Helper()

	d, err := parseDate(s)
	if err != nil {
		t.Fatalf("failed to parse date %q: %v", s, err)
	}

	return d
}
