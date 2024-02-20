package personnummer

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	lengthWithoutCentury = 10
	lengthWithCentury    = 12
)

var (
	errInvalidSecurityNumber = errors.New("Invalid swedish personal identity number")
	monthDays                = map[int]int{
		1:  31,
		3:  31,
		4:  30,
		5:  31,
		6:  30,
		7:  31,
		8:  31,
		9:  30,
		10: 31,
		11: 30,
		12: 31,
	}
	now            = time.Now
	rule3          = [...]int{0, 2, 4, 6, 8, 1, 3, 5, 7, 9}
	interimLetters = []rune("TRSUWXJKLMN")
)

// charsToDigit converts char bytes to a digit
// example: ['1', '1'] => 11
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

func runeInSlice(a rune, list []rune) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
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

		if !runeInSlice(c, interimLetters) {
			if c > '9' {
				return nil
			}
			if c < '0' {
				return nil
			}
		}

		cleanNumber = append(cleanNumber, byte(c))
	}

	return cleanNumber
}

// getCoOrdinationDay will return co-ordination day.
func getCoOrdinationDay(day []byte) []byte {
	d := charsToDigit(day)
	if d < 60 {
		return day
	}

	d -= 60

	if d < 10 {
		return []byte{'0', byte(d) + '0'}
	}

	return []byte{
		byte(d)/10 + '0',
		byte(d)%10 + '0',
	}
}

// luhn will test if the given string is a valid luhn string.
func luhn(s []byte) bool {
	for _, c := range s {
		if !unicode.IsNumber(rune(c)) {
			return false
		}
	}

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

// toString converts int to string.
func toString(in interface{}) string {
	switch v := in.(type) {
	case int, int32, int64, uint, uint32, uint64:
		return fmt.Sprint(v)
	case string:
		return v
	default:
		return ""
	}
}

// input time without centry.
func validateTime(time []byte) bool {
	length := len(time)

	date := charsToDigit(time[length-2 : length])
	month := charsToDigit(time[length-4 : length-2])

	if month != 2 {
		days, ok := monthDays[month]
		if !ok {
			return false
		}
		return date <= days
	}

	year := charsToDigit(time[:length-4])

	leapYear := year%4 == 0 && year%100 != 0 || year%400 == 0

	if leapYear {
		return date <= 29
	}
	return date <= 28
}

// Personnummer represents the personnummer struct.
type Personnummer struct {
	Century            string
	FullYear           string
	Year               string
	Month              string
	Day                string
	Sep                string
	Num                string
	Check              string
	leapYear           bool
	coordinationNumber bool
}

// Options represents the personnummer options.
type Options struct {
	AllowInterimNumber        bool
	DisableCoordinationNumber bool
}

// New parse a Swedish personal identity numbers and returns a new struct or a error.
func New(pin string, options ...*Options) (*Personnummer, error) {
	p := &Personnummer{}
	o := &Options{}

	if len(options) > 0 {
		o = options[0]
	}

	if err := p.parse(pin, o); err != nil {
		return nil, err
	}

	return p, nil
}

// parse Swedish personal identity numbers and set struct properpties or return a error.
func (p *Personnummer) parse(pin string, options *Options) error {
	var century, year, num, check string

	if pin == "" {
		return errInvalidSecurityNumber
	}

	dateBytes := getCleanNumber(pin)

	switch len(dateBytes) {
	case lengthWithCentury:
		century = string(dateBytes[0:2])
		year = string(dateBytes[2:4])
		num = string(dateBytes[8:11])
		check = string(dateBytes[11:])
		dateBytes = dateBytes[2:8]
		break
	case lengthWithoutCentury:
		year = string(dateBytes[0:2])
		num = string(dateBytes[6:9])
		check = string(dateBytes[9:])
		dateBytes = dateBytes[0:6]
		break
	default:
		return errInvalidSecurityNumber
	}

	if num == "000" {
		return errInvalidSecurityNumber
	}

	length := len(dateBytes)
	day := charsToDigit(dateBytes[length-2 : length])
	month := charsToDigit(dateBytes[length-4 : length-2])

	if month != 2 {
		if _, ok := monthDays[month]; !ok {
			return errInvalidSecurityNumber
		}
	}

	plus := strings.Contains(pin, "+")

	p.Century = century
	p.Year = year
	p.FullYear = toString(century + year)
	p.Check = check
	p.Num = num
	p.Sep = "-"
	p.Day = toString(fmt.Sprintf("%02d", day))
	p.Month = toString(fmt.Sprintf("%02d", month))

	if p.Century == "" {
		year := charsToDigit(dateBytes[:length-4])

		var baseYear int
		if plus {
			baseYear = now().Year() - 100
			p.Sep = "+"
		} else {
			baseYear = now().Year()
		}

		centuryStr := strconv.Itoa((baseYear - ((baseYear - year) % 100)))
		century, err := strconv.Atoi(centuryStr[0:2])
		if err != nil {
			return err
		}

		p.Century = toString(century)
		p.FullYear = toString(p.Century + p.Year)
	} else {
		fullYear, err := strconv.Atoi(century + year)
		if err != nil {
			return err
		}

		if now().Year()-fullYear < 100 {
			p.Sep = "-"
		} else {
			p.Sep = "+"
		}
	}

	if !p.valid() {
		return errInvalidSecurityNumber
	}

	if p.IsCoordinationNumber() && options.DisableCoordinationNumber {
		return errInvalidSecurityNumber
	}

	if p.IsInterimNumber() && !options.AllowInterimNumber {
		return errInvalidSecurityNumber
	}

	return nil
}

// Valid will validate Swedish personal identity numbers.
func (p *Personnummer) valid() bool {
	num := p.Num
	if p.IsInterimNumber() {
		num = "1" + p.Num[1:]
	}

	pin := fmt.Sprintf("%s%s%s%s%s%s", p.Century, p.Year, p.Month, p.Day, num, p.Check)

	bytes := []byte(pin)
	if !luhn(bytes[2:]) {
		return false
	}

	var dateBytes = append(bytes[:6], getCoOrdinationDay(bytes[6:8])...)

	return validateTime(dateBytes)
}

// Format a Swedish personal identity number as one of the official formats,
// a long format or a short format.
func (p *Personnummer) Format(longFormat ...bool) (string, error) {
	if len(longFormat) > 0 && longFormat[0] {
		return fmt.Sprintf("%s%s%s%s%s%s", p.Century, p.Year, p.Month, p.Day, p.Num, p.Check), nil
	}

	return fmt.Sprintf("%s%s%s%s%s%s", p.Year, p.Month, p.Day, p.Sep, p.Num, p.Check), nil
}

// GetDate returns the date from a Swedish personal identity number.
func (p *Personnummer) GetDate() time.Time {
	ageDay := charsToDigit([]byte(p.Day))

	if p.IsCoordinationNumber() {
		ageDay = ageDay - 60
	}

	fullYear := charsToDigit([]byte(p.FullYear))
	month := charsToDigit([]byte(p.Month))

	return time.Date(fullYear, time.Month(month), ageDay, 0, 0, 0, 0, time.UTC)
}

// GetAge returns the age from a Swedish personal identity number.
func (p *Personnummer) GetAge() int {
	t := p.GetDate()
	a := math.Floor(float64(now().Sub(t)/1e6) / 3.15576e+10)

	return int(a)
}

// IsCoordinationNumber determine if a Swedish personal identity number is a coordination number or not.
// Returns true if it's a coordination number.
func (p *Personnummer) IsCoordinationNumber() bool {
	day := charsToDigit([]byte(p.Day)) - 60
	str := fmt.Sprintf("%s%s%s", p.Century, p.Year, p.Month)
	if day < 10 {
		str += fmt.Sprintf("0%d", day)
	} else {
		str += fmt.Sprintf("%d", day)
	}
	return validateTime([]byte(str))
}

// IsInterimNumber determine if a Swedish personal identity number is a interim number or not.
// Returns true if it's a interim number.
func (p *Personnummer) IsInterimNumber() bool {
	return runeInSlice(rune(p.Num[0]), interimLetters)
}

// IsFemale checks if a Swedish personal identity number is for a female.
func (p *Personnummer) IsFemale() bool {
	return !p.IsMale()
}

// IsMale checks if a Swedish personal identity number is for a male.
// The second argument should be a boolean
func (p *Personnummer) IsMale() bool {
	sexDigit := int(p.Num[2])

	return sexDigit%2 == 1
}

// Valid will validate Swedish personal identity numbers
func Valid(pin string, options ...*Options) bool {
	_, err := Parse(pin, options...)
	return err == nil
}

// Parse Swedish personal identity numbers and return a new struct.
func Parse(pin string, options ...*Options) (*Personnummer, error) {
	return New(pin, options...)
}
