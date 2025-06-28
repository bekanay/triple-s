package validate

import (
	"errors"
	"strings"
	"unicode"
)

func Name(name string) error {
	if len(name) < 3 || len(name) > 63 {
		return errors.New("name must be 3–63 characters")
	}
	if name[0] == '-' || name[len(name)-1] == '-' {
		return errors.New("must not begin or end with hyphen")
	}
	if strings.Contains(name, "--") {
		return errors.New("must not contain consecutive hyphens")
	}
	if strings.Contains(name, "..") {
		return errors.New("must not contain consecutive periods")
	}
	for _, r := range name {
		if unicode.IsLower(r) || unicode.IsDigit(r) || r == '-' || r == '.' {
			continue
		}
		return errors.New("invalid character in name")
	}
	return nil
}
