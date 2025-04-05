package lang

import "fmt"

func TT(lang string, msg string, v ...any) string {

	if message, ok := TraslateMap[lang][msg]; ok {
		return fmt.Sprintf(message, v...)
	}

	return fmt.Sprintf(msg, v...)
}

func Traslate(lang string, msg string, v ...any) string {
	return TT(lang, msg, v...)
}

var TraslateMap = map[string]map[string]string{
	"en": {
		"No encontrado": "Not found",
	},
}
