package ws_upgrader

import (
	"testing"
)

func TestHash(t *testing.T) {
	data := map[string]string{
		"HSohHIffGH1RFigAVUNDYw==": "GZS2YkUYBCu6eW6qTtiqD2bqEFE=",
	}
	for raw, ret := range data {
		if hashAccKey(raw) != ret {
			t.Error("opps!")
		}
	}
}
