package organisationsnummer

import (
	"log"
	"os"
	"testing"

	"github.com/frozzare/go-assert"
	"github.com/frozzare/go/http2"
)

type TestListItem struct {
	LongFormat  string `json:"long_format"`
	ShortFormat string `json:"short_format"`
	Valid       bool   `json:"valid"`
	Type        string `json:"type"`
	Input       string `json:"input"`
	VatNumber   string `json:"vat_number"`
}

var testList []*TestListItem

func TestMain(m *testing.M) {
	if err := http2.GetJSON("https://raw.githubusercontent.com/organisationsnummer/meta/main/testdata/list.json", &testList); err != nil {
		log.Fatal(err)
	}

	code := m.Run()
	os.Exit(code)
}

func TestValidOrganisationsnummer(t *testing.T) {
	for _, item := range testList {
		if !item.Valid {
			continue
		}

		assert.True(t, Valid(item.Input))
	}
}
func TestInvalidOrganisationsnummer(t *testing.T) {
	for _, item := range testList {
		if item.Valid {
			continue
		}

		assert.False(t, Valid(item.Input))
	}
}

func TestValidOrganisationsnummerFormatShort(t *testing.T) {
	for _, item := range testList {
		if !item.Valid {
			continue
		}

		o, _ := Parse(item.Input)
		assert.Equal(t, o.Format(false), item.ShortFormat)
	}
}

func TestValidOrganisationsnummerFormatLong(t *testing.T) {
	for _, item := range testList {
		if !item.Valid {
			continue
		}

		o, _ := Parse(item.Input)
		assert.Equal(t, o.Format(true), item.LongFormat)
	}
}

func TestValidOrganisationsnummerType(t *testing.T) {
	for _, item := range testList {
		if !item.Valid {
			continue
		}

		o, _ := Parse(item.Input)
		assert.Equal(t, o.GetType(), item.Type)
	}
}

func TestValidOrganisationsnummerVatNumber(t *testing.T) {
	for _, item := range testList {
		if !item.Valid {
			continue
		}

		o, _ := Parse(item.Input)
		assert.Equal(t, o.VatNumber(), item.VatNumber)
	}
}

func TestValidPersonnummer(t *testing.T) {
	for _, item := range testList {
		if !item.Valid {
			continue
		}

		if item.Type != "Enskild firma" {
			continue
		}

		assert.True(t, Valid(item.LongFormat))
		org, _ := Parse(item.Input)
		assert.Equal(t, org.GetType(), item.Type)
		assert.Equal(t, org.String(), item.Type)
		assert.True(t, org.IsPersonnummer())
		assert.Equal(t, org.VatNumber(), item.VatNumber)
		assert.Equal(t, org.Format(true), item.LongFormat)
		assert.Equal(t, org.Format(false), item.ShortFormat)
	}
}
