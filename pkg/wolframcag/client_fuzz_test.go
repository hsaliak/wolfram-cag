package wolframcag

import "testing"

func FuzzDecodeJSON(f *testing.F) {
	f.Add([]byte(`{"result":"ok"}`))
	f.Add([]byte(`{"result":123}`))
	f.Add([]byte(`not-json`))

	f.Fuzz(func(t *testing.T, data []byte) {
		var out struct {
			Result string `json:"result"`
		}

		err := DecodeJSON(data, &out)
		if err == nil {
			if len(data) == 0 {
				t.Fatalf("empty payload should not decode successfully")
			}
		}
	})
}
