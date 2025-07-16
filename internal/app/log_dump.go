package app

import (
	"fmt"
	"reflect"
	"strings"
)

func (l *Logger) formatDump(val any) string {
	v := reflect.ValueOf(val)
	t := reflect.TypeOf(val)

	// Si es puntero, desreferenciar
	isPtr := false
	if v.Kind() == reflect.Ptr {
		isPtr = true
		if v.IsNil() {
			return fmt.Sprintf("*%s(nil)", t.Elem().Kind())
		}
		v = v.Elem()
		t = t.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		s := v.String()
		if isPtr {
			return fmt.Sprintf("*string(%d):\"%s\"", len(s), s)
		}
		return fmt.Sprintf("string(%d):\"%s\"", len(s), s)

	//----------------------------------------------------------------
	case reflect.Int32: // rune
		// Confirmar que es rune explícito, no solo int32 genérico
		if t.Name() == "rune" || t.String() == "rune" {
			val := v.Int()
			char := rune(val)
			if isPtr {
				return fmt.Sprintf("*rune(%d)'%c'", val, char)
			}
			return fmt.Sprintf("rune(%d)'%c'", val, char)
		}
		// si no es rune, seguir como int32 normal
		if isPtr {
			return fmt.Sprintf("*int32(%d)", v.Int())
		}
		return fmt.Sprintf("int32(%d)", v.Int())

	//----------------------------------------------------------------
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int64:
		if isPtr {
			return fmt.Sprintf("*%s(%d)", t.Kind(), v.Int())
		}
		return fmt.Sprintf("%s(%d)", t.Kind(), v.Int())

	//----------------------------------------------------------------
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if isPtr {
			return fmt.Sprintf("*%s(%d)", t.Kind(), v.Uint())
		}
		return fmt.Sprintf("%s(%d)", t.Kind(), v.Uint())

	//----------------------------------------------------------------
	case reflect.Float32, reflect.Float64:
		if isPtr {
			return fmt.Sprintf("*%s(%f)", t.Kind(), v.Float())
		}
		return fmt.Sprintf("%s(%f)", t.Kind(), v.Float())

	//----------------------------------------------------------------
	case reflect.Array:
		var b strings.Builder
		length := v.Len()
		if isPtr {
			b.WriteString(fmt.Sprintf("*array(%d){\n", length))
		} else {
			b.WriteString(fmt.Sprintf("array(%d){\n", length))
		}

		for i := 0; i < length; i++ {
			val := v.Index(i)
			valStr := l.formatDump(val.Interface())
			b.WriteString(fmt.Sprintf("  [%d] => %s,\n", i, valStr))
		}

		b.WriteString("}")
		return b.String()

	//----------------------------------------------------------------
	case reflect.Slice:
		if v.IsNil() {
			return "slice(nil){}"
		}

		var b strings.Builder
		length := v.Len()
		if isPtr {
			b.WriteString(fmt.Sprintf("*slice(%d){\n", length))
		} else {
			b.WriteString(fmt.Sprintf("slice(%d){\n", length))
		}

		for i := 0; i < length; i++ {
			val := v.Index(i)
			valStr := l.formatDump(val.Interface())
			b.WriteString(fmt.Sprintf("  [%d] => %s,\n", i, valStr))
		}

		b.WriteString("}")
		return b.String()

	//----------------------------------------------------------------
	case reflect.Map:
		if v.IsNil() {
			return "map(nil){}"
		}

		var b strings.Builder
		keys := v.MapKeys()
		if isPtr {
			b.WriteString(fmt.Sprintf("*map(%d){\n", len(keys)))
		} else {
			b.WriteString(fmt.Sprintf("map(%d){\n", len(keys)))
		}

		for _, k := range keys {
			keyStr := l.formatDump(k.Interface())
			valStr := l.formatDump(v.MapIndex(k).Interface())
			b.WriteString(fmt.Sprintf("  [%s] => %s,\n", keyStr, valStr))
		}

		b.WriteString("}")
		return b.String()

	//----------------------------------------------------------------
	case reflect.Struct:
		t := v.Type()
		var b strings.Builder

		name := t.Name()
		if name == "" {
			name = "anonymous"
		}

		if isPtr {
			b.WriteString(fmt.Sprintf("*struct(%s){\n", name))
		} else {
			b.WriteString(fmt.Sprintf("struct(%s){\n", name))
		}

		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			if field.PkgPath != "" {
				continue // campo no exportado
			}

			fieldName := field.Name
			val := l.formatDump(v.Field(i).Interface())

			rawTag := string(field.Tag)
			if rawTag != "" {
				// Parsear todas las etiquetas: key:"value"
				var tagParts []string
				pairs := strings.Split(rawTag, " ")
				for _, pair := range pairs {
					pair = strings.TrimSpace(pair)
					if pair != "" {
						tagParts = append(tagParts, pair)
					}
				}
				tagStr := strings.Join(tagParts, ", ")
				b.WriteString(fmt.Sprintf("  %s ((%d) %s) => %s,\n", fieldName, len(tagParts), tagStr, val))
			} else {
				b.WriteString(fmt.Sprintf("  %s => %s,\n", fieldName, val))
			}
		}

		b.WriteString("}")
		return b.String()

	//----------------------------------------------------------------
	case reflect.Bool:
		val := v.Bool()
		if isPtr {
			return fmt.Sprintf("*bool(%v)", val)
		}
		return fmt.Sprintf("bool(%v)", val)

	//----------------------------------------------------------------
	case reflect.Interface:
		if v.IsNil() {
			return "any(nil)"
		}

		inner := v.Elem()
		innerType := inner.Type().String()
		valStr := l.formatDump(inner.Interface())

		return fmt.Sprintf("any(%s) => %s", innerType, valStr)

	//----------------------------------------------------------------
	case reflect.Func:
		var inputTypes []string
		var outputTypes []string

		// Obtener tipos de entrada
		for i := 0; t.NumIn() > i; i++ {
			in := t.In(i)
			inputTypes = append(inputTypes, in.String())
		}

		// Obtener tipos de salida
		for i := 0; t.NumOut() > i; i++ {
			out := t.Out(i)
			outputTypes = append(outputTypes, out.String())
		}

		signature := fmt.Sprintf("func(%s)", strings.Join(inputTypes, ", "))
		if len(outputTypes) > 0 {
			signature += fmt.Sprintf(" -> %s", strings.Join(outputTypes, ", "))
		}

		if isPtr {
			return "*" + signature
		}
		return signature

	//----------------------------------------------------------------
	case reflect.Chan:
		elemType := t.Elem().String()

		if v.IsNil() {
			return fmt.Sprintf("chan(%s)[nil]", elemType)
		}

		switch t.ChanDir() {
		case reflect.SendDir:
			// Canal de solo envío
			return fmt.Sprintf("chan<-(%s)[send-only]", elemType)

		case reflect.RecvDir, reflect.BothDir:
			// Intentar recibir sin bloquear
			recv, ok := v.TryRecv()

			if !ok {
				// Canal cerrado
				if t.ChanDir() == reflect.RecvDir {
					return fmt.Sprintf("<-chan(%s)[closed]", elemType)
				}
				return fmt.Sprintf("chan(%s)[closed]", elemType)
			}

			// Canal abierto, valor recibido
			dumped := l.formatDump(recv.Interface())
			if t.ChanDir() == reflect.RecvDir {
				return fmt.Sprintf("<-chan(%s)[open: %s]", elemType, dumped)
			}
			return fmt.Sprintf("chan(%s)[open: %s]", elemType, dumped)

		default:
			// esto nunca se ejecutara pero el compilador me obliga
			return fmt.Sprintf("chan(%s)[unknown direction]", elemType)
		}

	//----------------------------------------------------------------
	case reflect.Invalid:
		return "invalid[nil]"

	//----------------------------------------------------------------
	case reflect.Complex64, reflect.Complex128:
		val := v.Complex()
		kind := t.Kind().String() // "complex64" o "complex128"

		if isPtr {
			return fmt.Sprintf("*%s(%f%+fi)", kind, real(val), imag(val))
		}
		return fmt.Sprintf("%s(%f%+fi)", kind, real(val), imag(val))

	default:
		if isPtr {
			return fmt.Sprintf("*%s: %v", t.Kind(), v.Interface())
		}
		return fmt.Sprintf("%s: %v", t.Kind(), v.Interface())
	}
}

// InterpolatePlaceholders reemplaza placeholders en el mensaje con valores del contexto
// Soporta formatos: {placeholder} y :placeholder
func InterpolatePlaceholders(msg string, ctx ...F) string {
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
