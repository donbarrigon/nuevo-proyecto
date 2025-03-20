package lang

import "fmt"

func M(lang string, msg string, v ...any) string {

	if _, ok := m[lang]; !ok {
		lang = "es"
	}

	if message, ok := m[lang][msg]; ok {
		return fmt.Sprintf(message, v...)
	}

	return msg
}

var m = map[string]map[string]string{
	"es": {
		"app.unautorized":     "No autorizado",
		"app.not-found":       "No Existe",
		"app.service.store":   "No se guardo",
		"app.service.update":  "No se guardo",
		"app.service.delete":  "No se elimino",
		"app.service.restore": "No se restauro",
		"app.service.destroy": "no se elimino permanentemente",
	},
}
