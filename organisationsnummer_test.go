package organisationsnummer

import (
	"os"
	"strings"
	"testing"

	"github.com/frozzare/go-assert"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestValidOrganisationsnummer(t *testing.T) {
	var numbers = []string{"556016-0680", "556103-4249", "5561034249", "165561034249"}

	for _, n := range numbers {
		assert.True(t, Valid(n))
	}
}
func TestInvalidOrganisationsnummer(t *testing.T) {
	var numbers = []string{"556016-0681", "556103-4250", "5561034250", "165561034250"}

	for _, n := range numbers {
		assert.False(t, Valid(n))
	}
}

func TestValidOrganisationsnummerFormatShort(t *testing.T) {
	var numbers = map[string]string{
		"556016-0680": "5560160680",
		"556103-4249": "5561034249",
		"5561034249":  "5561034249",
	}

	for k, v := range numbers {
		o, _ := Parse(k)
		assert.Equal(t, o.Format(false), v)
	}
}

func TestValidOrganisationsnummerFormatLong(t *testing.T) {
	var numbers = map[string]string{
		"5560160680":  "556016-0680",
		"5561034249":  "556103-4249",
		"556103-4249": "556103-4249",
	}

	for k, v := range numbers {
		o, _ := Parse(k)
		assert.Equal(t, o.Format(true), v)
	}
}

func TestValidOrganisationsnummerType(t *testing.T) {
	var numbers = map[string]string{
		"5560160680":  "Aktiebolag",
		"5561034249":  "Aktiebolag",
		"556103-4249": "Aktiebolag",
	}

	for k, v := range numbers {
		o, _ := Parse(k)
		assert.Equal(t, o.GetType(), v)
	}
}

func TestValidPersonnummer(t *testing.T) {
	var _type = "Enskild firma"
	var number = "121212121212"

	org, _ := Parse(number)
	assert.Equal(t, org.GetType(), _type)
	assert.Equal(t, org.String(), _type)
	assert.True(t, org.IsPersonnummer())
	assert.Equal(t, org.Personnummer().FullYear, "1212")
}

func TestValidPersonnummerFormat(t *testing.T) {
	var numbers = map[string]string{
		"121212121212":  "1212121212",
		"12121212-1212": "121212-1212",
	}

	for k, v := range numbers {
		o, _ := Parse(k)
		assert.Equal(t, o.Format(strings.Contains(k, "-")), v)
	}
}
