package api

import "testing"

func FuzzBuildResultQuery(f *testing.F) {
	f.Add("2+2", "", "", "", "", "", "", "")
	f.Add("integrate x", "assume", "plaintext", "metric", "Boston", "40,-71", "30", "800")

	f.Fuzz(func(t *testing.T, input, assumption, format, units, location, latlong, timeout, maxwidth string) {
		opts := ResultOptions{
			Assumption: assumption,
			Format:     format,
			Units:      units,
			Location:   location,
			LatLong:    latlong,
			Timeout:    timeout,
			MaxWidth:   maxwidth,
		}

		q, err := BuildResultQuery(input, opts)
		if input == "" {
			if err == nil {
				t.Fatalf("expected error for empty input")
			}
			return
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := q.Get("input"); got != input {
			t.Fatalf("input mismatch: got %q want %q", got, input)
		}
	})
}
