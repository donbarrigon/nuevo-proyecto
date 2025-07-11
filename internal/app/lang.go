package app

import "fmt"

func Translate(lang string, format string, fields ...F) string {
	if words, ok := TranslateMap[lang]; ok {

		if msg, found := words[format]; found {
			format = msg
		}

		// Traducir también los valores dinámicos si están en el diccionario
		for i, f := range fields {
			if translated, found := words[fmt.Sprint(f.Value)]; found {
				fields[i].Value = translated
			}
		}
	}
	return interpolatePlaceholders(format, fields)
}

var TranslateMap = map[string]map[string]string{
	"en": {
		"No encontrado": "Not found",
	},
}
