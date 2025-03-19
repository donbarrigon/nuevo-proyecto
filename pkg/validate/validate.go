package validate

import (
	"regexp"
)

func Email(value string) bool {
	// Expresión regular para validar el formato del email
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(regex, value)
	return matched
}
