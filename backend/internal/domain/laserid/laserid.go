package laserid

import "regexp"

var pattern = regexp.MustCompile(`^[A-Z]{2}\d-?\d{7}-?\d{2}$`)

func IsValidFormat(s string) bool { return pattern.MatchString(s) }
