package birthdate

import "time"

type Precision string

const (
	Full      Precision = "FULL"
	YearMonth Precision = "YEAR_MONTH"
	YearOnly  Precision = "YEAR_ONLY"
)

const MinYear = 1900

const (
	FieldYear      = "birth_year"
	FieldMonth     = "birth_month"
	FieldDay       = "birth_day"
	FieldPrecision = "birth_precision"
)

type BirthDate struct {
	Year      int
	Month     *int
	Day       *int
	Precision Precision
}

func (b BirthDate) Validate(today time.Time) []string {
	var bad []string

	if b.Year < MinYear || b.Year > today.Year() {
		bad = append(bad, FieldYear)
	}

	switch b.Precision {
	case Full:
		if b.Month == nil {
			bad = append(bad, FieldMonth)
		}
		if b.Day == nil {
			bad = append(bad, FieldDay)
		}
	case YearMonth:
		if b.Month == nil {
			bad = append(bad, FieldMonth)
		}
		if b.Day != nil {
			bad = append(bad, FieldDay)
		}
	case YearOnly:
		if b.Month != nil {
			bad = append(bad, FieldMonth)
		}
		if b.Day != nil {
			bad = append(bad, FieldDay)
		}
	default:
		bad = append(bad, FieldPrecision)
	}

	if b.Month != nil && (*b.Month < 1 || *b.Month > 12) && !has(bad, FieldMonth) {
		bad = append(bad, FieldMonth)
	}
	if b.Month != nil && b.Day != nil && !has(bad, FieldMonth) && !has(bad, FieldDay) {
		d := time.Date(b.Year, time.Month(*b.Month), *b.Day, 0, 0, 0, 0, time.UTC)
		if d.Day() != *b.Day || int(d.Month()) != *b.Month {
			bad = append(bad, FieldDay)
		}
	}

	if len(bad) == 0 {
		y, m, d := today.Date()
		if b.Earliest().After(time.Date(y, m, d, 0, 0, 0, 0, time.UTC)) {
			bad = append(bad, b.lastKnownField())
		}
	}
	return bad
}

func (b BirthDate) Earliest() time.Time {
	m, d := 1, 1
	if b.Month != nil {
		m = *b.Month
	}
	if b.Day != nil {
		d = *b.Day
	}
	return time.Date(b.Year, time.Month(m), d, 0, 0, 0, 0, time.UTC)
}

func (b BirthDate) AgeAt(t time.Time) int {
	e := b.Earliest()
	age := t.Year() - e.Year()
	if t.Month() < e.Month() || (t.Month() == e.Month() && t.Day() < e.Day()) {
		age--
	}
	return age
}

func (b BirthDate) lastKnownField() string {
	switch {
	case b.Day != nil:
		return FieldDay
	case b.Month != nil:
		return FieldMonth
	default:
		return FieldYear
	}
}

func has(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}
