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

type List []Entry

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

func (l *List) Set(key string, value any) {
	*l = append(*l, Entry{Key: key, Value: value})
}

func (l *List) Get(key string) any {
	for _, entry := range *l {
		if entry.Key == key {
			return entry.Value
		}
	}
	return nil
}

func (l *List) Has(key string) bool {
	for _, entry := range *l {
		if entry.Key == key {
			return true
		}
	}
	return false
}

func (l *List) Remove(key string) {
	for i, entry := range *l {
		if entry.Key == key {
			*l = append((*l)[:i], (*l)[i+1:]...)
			return
		}
	}
}

func (l *List) Clear() {
	*l = (*l)[:0]
}

func (l *List) Len() int {
	return len(*l)
}

func (l *List) Keys() []string {
	keys := make([]string, 0, len(*l))
	for _, entry := range *l {
		keys = append(keys, entry.Key)
	}
	return keys
}

func (l *List) Values() []any {
	values := make([]any, 0, len(*l))
	for _, entry := range *l {
		values = append(values, entry.Value)
	}
	return values
}

func (l List) MarshalJSON() ([]byte, error) {
	m := make(map[string]any, len(l))
	for _, field := range l {
		m[field.Key] = field.Value
	}
	return json.Marshal(m)
}

func E(key string, value any) Entry {
	return Entry{Key: key, Value: value}
}
