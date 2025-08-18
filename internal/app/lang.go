package app

import "fmt"

func Translate(lang string, format string, fields ...Item) string {
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
	return InterpolatePlaceholders(format, fields...)
}

var TranslateMap = map[string]map[string]string{
	"es": {},
}
