package appno

import (
	"fmt"
	"regexp"
)

const Prefix = "WC"

var pattern = regexp.MustCompile(`^WC-\d{4}-\d{7,}$`)

func Format(fiscalYear int, seq int64) string {
	return fmt.Sprintf("%s-%04d-%07d", Prefix, fiscalYear, seq)
}

func IsValid(s string) bool { return pattern.MatchString(s) }
