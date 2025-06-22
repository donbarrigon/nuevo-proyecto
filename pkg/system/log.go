package system

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"
)

const (
	LOG_DEBUG = iota
	LOG_INFO
	LOG_WARN
	LOG_ERROR
	LOG_FATAL
)

const (
	LOG_FLAG_TIMESTAMP       = 1 << iota // 1     // Agrega la fecha y hora formateada según LOG_DATE_FORMAT
	LOG_FLAG_LONGFILE                    // 2     // Ruta completa del archivo y número de línea: /a/b/c/d.go:23
	LOG_FLAG_SHORTFILE                   // 4     // Solo el nombre del archivo y línea: d.go:23
	LOG_FLAG_RELATIVEFILE                // 8     // Ruta relativa al directorio del proyecto
	LOG_FLAG_FUNCTION                    // 16    // Nombre de la función desde donde se llamó
	LOG_FLAG_LINE                        // 32    // Solo el número de línea (sin ruta de archivo)
	LOG_FLAG_PREFIX                      // 64    // Agrega un prefijo antes del mensaje (por ejemplo: [DEBUG])
	LOG_FLAG_CONSOLE_AS_JSON             // 128   // Salida en formato JSON en la consola
	LOG_FLAG_FILE_AS_JSON                // 256   // Salida en formato JSON en el archivo
	LOG_FLAG_CONTEXT                     // 512   // Agrega el contexto de la petición al log
	LOG_FLAG_DETAIL                      // 1024  // Las variables las se imprimen de forma detallada

	// Combinación de todos los flags
	LOG_FLAG_ALL = LOG_FLAG_TIMESTAMP |
		// LOG_FLAG_LONGFILE |
		// LOG_FLAG_SHORTFILE |
		LOG_FLAG_RELATIVEFILE |
		LOG_FLAG_FUNCTION |
		LOG_FLAG_LINE |
		LOG_FLAG_PREFIX |
		LOG_FLAG_CONSOLE_AS_JSON |
		LOG_FLAG_FILE_AS_JSON |
		LOG_FLAG_CONTEXT |
		LOG_FLAG_DETAIL
)

const (
	LOG_OUTPUT_CONSOLE  = 1 << iota // 1 - salida por consola estándar
	LOG_OUTPUT_FILE                 // 2 - salida a archivo
	LOG_OUTPUT_DATABASE             // 4 - guardar logs en base de datos
	LOG_OUTPUT_REMOTE               // 8 - enviar a un servidor remoto (opcional)
)

type LogLevel int

type Logger struct {
	Level     LogLevel
	Flags     int
	Output    int
	RemoteURL string
	Path      string
	Context   map[string]any
}

var Log = Logger{
	Level:     LOG_DEBUG,
	Flags:     LOG_FLAG_ALL,
	Output:    LOG_OUTPUT_CONSOLE | LOG_OUTPUT_DATABASE,
	RemoteURL: "http://127.0.0.1/debug/log",
	Path:      "log.json",
}

func (lv LogLevel) String() string {
	switch lv {
	case LOG_DEBUG:
		return "DEBUG"
	case LOG_INFO:
		return "INFO"
	case LOG_WARN:
		return "WARN"
	case LOG_ERROR:
		return "ERROR"
	case LOG_FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

func (l *Logger) Debug(format string, a ...any) {
	if l.Level <= LOG_DEBUG {
		go l.output(l.Level.String(), format, a...)
	}
}

func (l *Logger) Info(format string, a ...any) {
	if l.Level <= LOG_INFO {
		go l.output(l.Level.String(), format, a...)
	}
}

func (l *Logger) Warning(format string, a ...any) {
	if l.Level <= LOG_WARN {
		go l.output(l.Level.String(), format, a...)
	}
}

func (l *Logger) Error(format string, a ...any) {
	if l.Level <= LOG_ERROR {
		go l.output(l.Level.String(), format, a...)
	}
}

func (l *Logger) Fatal(format string, a ...any) {
	if l.Level <= LOG_FATAL {
		go l.output(l.Level.String(), format, a...)
	}
}

func (l *Logger) Print(format string, a ...any) {
	go l.output("PRINT", format, a...)
}

func (l *Logger) WithContext(level LogLevel, ctx map[string]any, format string, a ...any) {
	if l.Level > level {
		return
	}

	copy := *l
	copy.Context = ctx
	go copy.output(level.String(), format, a...)
}

func (l *Logger) output(level string, format string, a ...any) {
	// Obtener información del runtime
	pc, file, line, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()

	// Si RELATIVEFILE está activado
	if l.Flags&LOG_FLAG_RELATIVEFILE != 0 {
		if wd, err := os.Getwd(); err == nil {
			if rel, err := filepath.Rel(wd, file); err == nil {
				file = rel
			}
		}
	}

	// Si SHORTFILE está activado
	if l.Flags&LOG_FLAG_SHORTFILE != 0 {
		file = filepath.Base(file)
	}

	// Preparar mensaje
	if l.Flags&LOG_FLAG_DETAIL != 0 {
		// poner el tipo de dato de la de cada variable del slice a. si es string(len) "value", array[type(len)]{"value1","value2"}, map[type(len)]type{(len_del_key)key:(len_del_value)"value"}, struct igual que el map
	}
	msg := Translate(Env.APP_LOCALE, format, a...)

	// Crear estructura de log
	entry := map[string]any{
		"level":   level,
		"message": msg,
	}

	if l.Flags&LOG_FLAG_TIMESTAMP != 0 {
		now := time.Now().Format(Env.LOG_DATE_FORMAT)
		entry["time"] = now
	}

	if l.Flags&LOG_FLAG_FUNCTION != 0 {
		entry["function"] = funcName
	}

	if l.Flags&LOG_FLAG_LINE != 0 {
		entry["line"] = line
	}

	if l.Flags&(LOG_FLAG_LONGFILE|LOG_FLAG_SHORTFILE|LOG_FLAG_RELATIVEFILE) != 0 {
		entry["file"] = file
	}

	if l.Flags&LOG_FLAG_CONTEXT != 0 && l.Context != nil {
		entry["context"] = l.Context
	}

	// salida en consola
	if l.Output&LOG_OUTPUT_CONSOLE != 0 {
		if l.Flags&LOG_FLAG_CONSOLE_AS_JSON != 0 {
			data, _ := json.MarshalIndent(entry, "", "  ")
			fmt.Println(string(data))
		} else {
			var b strings.Builder

			if l.Flags&LOG_FLAG_TIMESTAMP != 0 {
				b.WriteString(fmt.Sprintf("%s ", entry["time"]))
			}
			if l.Flags&LOG_FLAG_PREFIX != 0 {
				b.WriteString(fmt.Sprintf("[%s] ", entry["level"]))
			}
			b.WriteString(fmt.Sprintf("%s", entry["message"]))

			if l.Flags&LOG_FLAG_FUNCTION != 0 {
				b.WriteString(fmt.Sprintf(" [%s]", entry["function"]))
			}
			if l.Flags&(LOG_FLAG_LONGFILE|LOG_FLAG_SHORTFILE|LOG_FLAG_RELATIVEFILE|LOG_FLAG_LINE) != 0 {
				b.WriteString(fmt.Sprintf(" (%s:%d)", entry["file"], entry["line"]))
			}

			fmt.Println(b.String())

			// Detalle de argumentos
			if l.Flags&LOG_FLAG_DETAIL != 0 && len(a) > 0 {
				for i, arg := range a {
					detail := l.formatDump(arg)
					fmt.Printf("  arg[%d]: %s\n", i, detail)
				}
			}
		}
	}

}

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
		b.WriteString(fmt.Sprintf("array(%d){\n", length))

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
		b.WriteString(fmt.Sprintf("slice(%d){\n", length))

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
		b.WriteString(fmt.Sprintf("map(%d){\n", len(keys)))

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

		b.WriteString(fmt.Sprintf("struct(%s){\n", name))

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
