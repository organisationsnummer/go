package organisationsnummer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
)

type TestListItem struct {
	LongFormat  string `json:"long_format"`
	ShortFormat string `json:"short_format"`
	Valid       bool   `json:"valid"`
	Type        string `json:"type"`
	Input       string `json:"input"`
	VatNumber   string `json:"vat_number"`
}

func getJSON(url string, v any) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s: %s", url, res.Status)
	}
	return json.NewDecoder(res.Body).Decode(v)
}

func assertEqual[T comparable](t *testing.T, expected, actual T) {
	t.Helper()
	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

var testList []*TestListItem

func TestMain(m *testing.M) {
	if err := getJSON("https://raw.githubusercontent.com/organisationsnummer/meta/main/testdata/list.json", &testList); err != nil {
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

		assertEqual(t, true, Valid(item.Input))
	}
}
func TestInvalidOrganisationsnummer(t *testing.T) {
	for _, item := range testList {
		if item.Valid {
			continue
		}

		assertEqual(t, false, Valid(item.Input))
	}
}

func TestValidOrganisationsnummerFormatShort(t *testing.T) {
	for _, item := range testList {
		if !item.Valid {
			continue
		}

		o, _ := Parse(item.Input)
		assertEqual(t, o.Format(false), item.ShortFormat)
	}
}

func TestValidOrganisationsnummerFormatLong(t *testing.T) {
	for _, item := range testList {
		if !item.Valid {
			continue
		}

		o, _ := Parse(item.Input)
		assertEqual(t, o.Format(true), item.LongFormat)
	}
}

func TestValidOrganisationsnummerType(t *testing.T) {
	for _, item := range testList {
		if !item.Valid {
			continue
		}

		o, _ := Parse(item.Input)
		assertEqual(t, o.GetType(), item.Type)
	}
}

func TestValidOrganisationsnummerVatNumber(t *testing.T) {
	for _, item := range testList {
		if !item.Valid {
			continue
		}

		o, _ := Parse(item.Input)
		assertEqual(t, o.VatNumber(), item.VatNumber)
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

		assertEqual(t, true, Valid(item.LongFormat))
		org, _ := Parse(item.Input)
		assertEqual(t, org.GetType(), item.Type)
		assertEqual(t, org.String(), item.Type)
		assertEqual(t, true, org.IsPersonnummer())
		assertEqual(t, org.VatNumber(), item.VatNumber)
		assertEqual(t, org.Format(true), item.LongFormat)
		assertEqual(t, org.Format(false), item.ShortFormat)
	}
}
