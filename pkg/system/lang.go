package system

import "fmt"

func Translate(lang string, format string, a ...any) string {

	if message, ok := TraslateMap[lang][format]; ok {
		return fmt.Sprintf(message, a...)
	}

	return fmt.Sprintf(format, a...)
}

var TraslateMap = map[string]map[string]string{
	"en": {
		"No encontrado": "Not found",
	},
}
