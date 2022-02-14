package organisationsnummer

import (
	"errors"
	"fmt"

	personnummer "github.com/personnummer/go/v3"
)

var (
	errInvalidSecurityNumber = errors.New("Invalid Swedish organization number")
	rule3                    = [...]int{0, 2, 4, 6, 8, 1, 3, 5, 7, 9}
	unkown                   = "Okänt"
	types                    = map[byte]string{
		'1': "Dödsbon",
		'2': "Stat, landsting, kommun eller församling",
		'3': "Utländska företag som bedriver näringsverksamhet eller äger fastigheter i Sverige",
		'5': "Aktiebolag",
		'6': "Enkelt bolag",
		'7': "Ekonomisk förening eller bostadsrättsförening",
		'8': "Ideella förening och stiftelse",
		'9': "Handelsbolag, kommanditbolag och enkelt bolag",
	}
)

// charsToDigit converts char bytes to a digit
// example: ['1", '1'] => 11
func charsToDigit(chars []byte) int {
	l := len(chars)
	r := 0
	for i, c := range chars {
		p := int((c - '0'))
		for j := 0; j < l-i-1; j++ {
			p *= 10
		}
		r += p
	}
	return r
}

// getCleanNumber will return clean numbers.
func getCleanNumber(in string) []byte {
	cleanNumber := make([]byte, 0, len(in))

	for _, c := range in {
		if c == '+' {
			continue
		}
		if c == '-' {
			continue
		}

		if c > '9' {
			return nil
		}
		if c < '0' {
			return nil
		}

		cleanNumber = append(cleanNumber, byte(c))
	}

	return cleanNumber
}

// luhn will test if the given string is a valid luhn string.
func luhn(s []byte) bool {
	odd := len(s) & 1

	var sum int

	for i, c := range s {
		if i&1 == odd {
			sum += rule3[c-'0']
		} else {
			sum += int(c - '0')
		}
	}

	return sum%10 == 0
}

// Organisationsnummer represents the organisationsnummer struct.
type Organisationsnummer struct {
	number       string
	personnummer *personnummer.Personnummer
}

// New parse a Swedish organization numbers and returns a new struct or a error.
func New(input string) (*Organisationsnummer, error) {
	o := &Organisationsnummer{}

	if err := o.parse(input); err != nil {
		return nil, err
	}

	return o, nil
}

// parse Swedish organization numbers and set struct properpties or return a error.
func (o *Organisationsnummer) parse(input string) error {
	number := getCleanNumber(input)
	p, err := personnummer.Parse(input)

	if err == nil {
		o.personnummer = p
		o.number = string(number)
	} else if len(number) == 12 {
		// May only be prefixed with 16.
		if charsToDigit(number[0:2]) != 16 {
			return errInvalidSecurityNumber
		}

		number = number[2:]
	}

	if len(number) == 10 {
		// Third digit bust be more than 20.
		if charsToDigit(number[2:4]) < 20 {
			return errInvalidSecurityNumber
		}

		// May not start with leading 0.
		if charsToDigit(number[0:2]) < 10 {
			return errInvalidSecurityNumber
		}

		if !luhn(number) {
			return errInvalidSecurityNumber
		}

		o.number = string(getCleanNumber(input))
	}

	return nil
}

// Get Personnummer instance
func (o *Organisationsnummer) Personnummer() personnummer.Personnummer {
	return *o.personnummer
}

// Determine if personnummer or not.
func (o *Organisationsnummer) IsPersonnummer() bool {
	return o.personnummer != nil
}

// Format a Swedish organization number as one of the official formats,
// a long format or a short format.
func (o *Organisationsnummer) Format(separator ...bool) string {
	var number = o.number

	if o.IsPersonnummer() {
		f, _ := o.personnummer.Format(true)
		number = f[2:]
	}

	if len(separator) > 0 && separator[0] {
		return fmt.Sprintf("%s-%s", number[0:6], number[6:])
	}

	return number
}

// Get the organization type.
func (o *Organisationsnummer) GetType() string {
	if o.IsPersonnummer() {
		return "Enskild firma"
	}

	if types[o.number[0]] != "" {
		return types[o.number[0]]
	}

	return unkown
}

// Get the organization type.
func (o *Organisationsnummer) String() string {
	return o.GetType()
}

// Get vat number for a organization number.
func (o *Organisationsnummer) vatNumber() string {
	return fmt.Sprintf("SE%s01", o.Format(false))
}

// Valid will validate Swedish organization numbers
func Valid(input string) bool {
	_, err := Parse(input)
	return err == nil
}

// Parse Swedish organization numbers and return a new struct.
func Parse(input string) (*Organisationsnummer, error) {
	return New(input)
}
