package issuer

import (
	_ "embed"
	"encoding/csv"
	"strings"
	"testing"
)

//go:embed valid-cards.csv
var validCardNumberList string

func TestIdentify(t *testing.T) {
	r := csv.NewReader(strings.NewReader(validCardNumberList))
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("error parsing credit card list: %s", err)
	}

	for _, record := range records[1:] {
		want, cardNumber := record[0], record[1]
		have := Identify(cardNumber)
		if have.String() != want {
			t.Errorf("incorrect issuer: want %s have %s (%s)", want, have, cardNumber)
		}
	}
}
