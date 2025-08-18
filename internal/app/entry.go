package app

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Entry struct {
	Key   string
	Value any
}

type EntryList []Entry

func (e EntryList) MarshalJSON() ([]byte, error) {
	m := make(map[string]any, len(e))
	for _, field := range e {
		m[field.Key] = field.Value
	}
	return json.Marshal(m)
}

// InterpolatePlaceholders reemplaza placeholders en el mensaje con valores del contexto
// Soporta formatos: {placeholder} y :placeholder
func InterpolatePlaceholders(msg string, ctx ...Entry) string {
	if len(ctx) == 0 {
		return msg
	}

	for _, field := range ctx {
		// Crear ambos formatos de placeholder
		placeholder1 := fmt.Sprintf("{%s}", field.Key) // Formato {key}
		placeholder2 := fmt.Sprintf(":%s", field.Key)  // Formato :key
		valueStr := fmt.Sprint(field.Value)

		// Reemplazar ambos formatos
		msg = strings.ReplaceAll(msg, placeholder1, valueStr)
		msg = strings.ReplaceAll(msg, placeholder2, valueStr)
	}

	return msg
}
