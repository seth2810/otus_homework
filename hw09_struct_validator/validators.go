package hw09structvalidator

import (
	"fmt"
	"regexp"
)

func validateStrLen(str string, length int) error {
	if len(str) != length {
		return fmt.Errorf("must be %d characters long", length)
	}

	return nil
}

func validateStrIn(str string, values []string) error {
	for _, v := range values {
		if str == v {
			return nil
		}
	}

	return fmt.Errorf("must be one of %q", values)
}

func validateStrRegExp(str string, exp *regexp.Regexp) error {
	if !exp.MatchString(str) {
		return fmt.Errorf("must satisfy regular expression '%s'", exp.String())
	}

	return nil
}

func validateIntMin(val int64, min int) error {
	if val < int64(min) {
		return fmt.Errorf("must be greater than %d", min)
	}

	return nil
}

func validateIntMax(val int64, max int) error {
	if val > int64(max) {
		return fmt.Errorf("must be less than %d", max)
	}

	return nil
}

func validateIntIn(val int64, values []int) error {
	for _, v := range values {
		if val == int64(v) {
			return nil
		}
	}

	return fmt.Errorf("must be one of %v", values)
}
