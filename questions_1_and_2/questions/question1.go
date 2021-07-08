package questions

import (
	"fmt"
	"strconv"
	"unicode"
)

const (
	idSize = 13
)

// SAIDDetails respresents fields of an id string
type SAIDDetails struct {
	DateOfBirth struct {
		Year  int
		Month int
		Day   int
	}
	Gender      string
	CitizenShip bool
}

// ErrorYear could not be parsed
var ErrorYear = fmt.Errorf("could not parse the year")

// ErrorMonth could not be parsed
var ErrorMonth = fmt.Errorf("could not parse the month")

// ErrorDay could not be parsed
var ErrorDay = fmt.Errorf("could not parse the day")

// ErrorSize when the size is not equal to 13
var ErrorSize = fmt.Errorf("wrong ID length size")

// ErrorGender when the gender code is not between 0 and 9999
var ErrorGender = fmt.Errorf("unknown gender code")

// ErrorCitizenship when the citizenship code is not 0 or 1
var ErrorCitizenship = fmt.Errorf("unknown Citizenship code")

// ErrorCheckDigit when this value is not a digit
var ErrorCheckDigit = fmt.Errorf("bad check digit")

func luhnCheck(id string) bool {
	// you can safely and confidently do direct conversion because all
	// the indexes has been traversed at this stage without any errors.
	var digits byte = 0
	for i := 0; i < len(id)-1; i++ {
		val := id[i] - '0'
		if i%2 == 0 {
			if val > 4 {
				val = 1 + (5-val)*2
			} else {
				val *= 2
			}
		}
		digits += val
	}

	checkDigit := id[12] - '0'
	return (digits+checkDigit)%10 == 0
}

// Parse populates fields of a id string into a SAIDDetails struct
func Parse(id string) (*SAIDDetails, error) {
	if size := len(id); size != idSize {
		return nil, ErrorSize
	}
	details := &SAIDDetails{}

	// use strconv to defend against none digit entries
	year, err := strconv.Atoi(id[:2])
	if err != nil {
		return nil, ErrorYear
	}
	details.DateOfBirth.Year = year
	month, err := strconv.Atoi(id[2:4])
	if err != nil {
		return nil, ErrorMonth
	}
	details.DateOfBirth.Month = month
	day, err := strconv.Atoi(id[4:6])
	if err != nil {
		return nil, ErrorDay
	}
	details.DateOfBirth.Day = day
	gender, err := strconv.Atoi(id[6:10])
	if err != nil {
		return nil, ErrorGender
	}
	if gender >= 0 && gender <= 4999 {
		details.Gender = "MALE"
	} else if gender >= 5000 && gender <= 9999 {
		details.Gender = "FEMALE"
	}

	details.DateOfBirth.Day = day
	citizenship, err := strconv.Atoi(id[10:11])
	if err != nil {
		return nil, err
	}
	if citizenship == 1 {
		details.CitizenShip = false
	} else if citizenship == 0 {
		details.CitizenShip = true
	} else {
		return nil, ErrorCitizenship
	}

	if !unicode.IsDigit(rune(id[12])) {
		return nil, ErrorCheckDigit
	}

	return details, nil
}
