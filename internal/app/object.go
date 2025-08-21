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

type Object []Entry

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

func (l *Object) Set(key string, value any) {
	*l = append(*l, Entry{Key: key, Value: value})
}

func (l *Object) Get(key string) any {
	for _, Entry := range *l {
		if Entry.Key == key {
			return Entry.Value
		}
	}
	return nil
}

func (l *Object) Has(key string) bool {
	for _, Entry := range *l {
		if Entry.Key == key {
			return true
		}
	}
	return false
}

func (l *Object) Remove(key string) {
	for i, Entry := range *l {
		if Entry.Key == key {
			*l = append((*l)[:i], (*l)[i+1:]...)
			return
		}
	}
}

func (l *Object) Clear() {
	*l = (*l)[:0]
}

func (l *Object) Len() int {
	return len(*l)
}

func (l *Object) Keys() []string {
	keys := make([]string, 0, len(*l))
	for _, Entry := range *l {
		keys = append(keys, Entry.Key)
	}
	return keys
}

func (l *Object) Values() []any {
	values := make([]any, 0, len(*l))
	for _, Entry := range *l {
		values = append(values, Entry.Value)
	}
	return values
}

func (l Object) MarshalJSON() ([]byte, error) {
	m := make(map[string]any, len(l))
	for _, field := range l {
		m[field.Key] = field.Value
	}
	return json.Marshal(m)
}

func E(key string, value any) Entry {
	return Entry{Key: key, Value: value}
}
