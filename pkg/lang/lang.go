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
		"app.unautorized":                "No autorizado",
		"app.not-found":                  "No Existe",
		"app.unprocessable-entity":       "No procesable",
		"app.bad-request":                "Solicitud incorrecta",
		"app.method-not-allowed":         "Metodo no soportado",
		"app.internal-server-error":      "Error interno del servidor",
		"app.internal-error":             "Error interno del servidor",
		"app.request.min.txt":            "Minimo %v caracteres",
		"app.request.min.num":            "Minimo %v",
		"app.request.max.txt":            "Maximo %v caracteres",
		"app.request.max.num":            "Maximo %v",
		"app.request.email":              "El formato de correo electronico es incorrecto",
		"app.request.required":           "Es requerido",
		"app.request.query-params":       "Parametros incorrect",
		"app.service.store":              "No se guardo",
		"app.service.update":             "No se guardo",
		"app.service.delete":             "No se elimino",
		"app.service.restore":            "No se restauro",
		"app.service.destroy":            "no se elimino permanentemente",
		"guard.auth.expired":             "La session ha expirado",
		"guard.auth.header":              "Encabezado de autorización no válido",
		"guard.auth.header-format":       "Formato de encabezado de autorización no válido",
		"guard.auth.invalid-token":       "El token no es válido: %v",
		"user.request.required":          "El email o telefono son requeridos",
		"user.service.generate-password": "Error al encriptar la contraseña",
		"user.service.unautorized":       "Las credenciales no son validas",
		"user.logout.destroy":            "Intentelo de nuevo mas tarde",
	},
}
