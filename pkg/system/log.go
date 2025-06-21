package system

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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
	LOG_FLAG_DATE         = 1 << iota // 1      // la fecha: 2006/01/02
	LOG_FLAG_TIME                     // 2      // la hora: 15:04:05
	LOG_FLAG_MICROSECONDS             // 4      // microsegundos: 15:04:05.000000
	LOG_FLAG_LONGFILE                 // 8      // ruta completa del archivo + línea: /a/b/c/d.go:23
	LOG_FLAG_SHORTFILE                // 16     // solo el archivo + línea: d.go:23
	LOG_FLAG_RELATIVEFILE             // 32     // ruta relativa al directorio del proyecto
	LOG_FLAG_FUNCTION                 // 64     // nombre de la función desde donde se llamó
	LOG_FLAG_LINE                     // 128    // número de línea (sin ruta)
	LOG_FLAG_PREFIX                   // 256    // añade el prefijo antes del mensaje [DEBUG]
	LOG_FLAG_JSON                     // 512    // hace que el formato de salida sea json si no esta la salida del arhivo es csvy en la consola es inline
	LOG_FLAG_CONTEXT                  // 1024   // agrega el context de la peticion al log

	LOG_FLAG_STD       = LOG_FLAG_DATE | LOG_FLAG_TIME                         // Combinación estándar por defecto: fecha + hora
	LOG_FLAG_TIMESTAMP = LOG_FLAG_DATE | LOG_FLAG_TIME | LOG_FLAG_MICROSECONDS // Combinación estándar por defecto: fecha + hora + microsegundos

	// Combinación de todos los flags
	LOG_FLAG_ALL = LOG_FLAG_DATE |
		LOG_FLAG_TIME |
		LOG_FLAG_MICROSECONDS |
		// LOG_FLAG_LONGFILE |
		// LOG_FLAG_SHORTFILE |
		LOG_FLAG_RELATIVEFILE |
		LOG_FLAG_FUNCTION |
		LOG_FLAG_LINE |
		LOG_FLAG_PREFIX |
		LOG_FLAG_JSON |
		LOG_FLAG_CONTEXT
)

const (
	LOG_OUTPUT_CONSOLE  = iota // 0 - salida por consola estándar
	LOG_OUTPUT_FILE            // 1 - salida a archivo
	LOG_OUTPUT_DATABASE        // 2 - guardar logs en base de datos
	LOG_OUTPUT_REMOTE          // 3 - enviar a un servidor remoto (opcional)
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
		go l.outputEntry(l.Level.String(), format, a...)
	}
}

func (l *Logger) Info(format string, a ...any) {
	if l.Level <= LOG_INFO {
		go l.outputEntry(l.Level.String(), format, a...)
	}
}

func (l *Logger) Warning(format string, a ...any) {
	if l.Level <= LOG_WARN {
		go l.outputEntry(l.Level.String(), format, a...)
	}
}

func (l *Logger) Error(format string, a ...any) {
	if l.Level <= LOG_ERROR {
		go l.outputEntry(l.Level.String(), format, a...)
	}
}

func (l *Logger) Fatal(format string, a ...any) {
	if l.Level <= LOG_FATAL {
		go l.outputEntry(l.Level.String(), format, a...)
	}
}

func (l *Logger) Print(format string, a ...any) {
	go l.outputEntry("PRINT", format, a...)
}

func (l *Logger) outputEntry(level string, format string, a ...any) {
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
	msg := Translate(Env.APP_LOCALE, format, a...)

	// Crear estructura de log
	entry := map[string]any{
		"level":   level,
		"message": msg,
		"time":    time.Now().Format("2006-01-02 15:04:05"),
	}

	if l.Flags&LOG_FLAG_DATE != 0 || l.Flags&LOG_FLAG_TIME != 0 {
		entry["timestamp"] = time.Now().Format(time.RFC3339)
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

	// Formato JSON
	if l.Flags&LOG_FLAG_JSON != 0 {
		data, _ := json.Marshal(entry)
		l.output(string(data))
		return
	}

	// Formato inline tipo: [LEVEL] mensaje (file:line)
	var buf bytes.Buffer
	if l.Flags&LOG_FLAG_PREFIX != 0 {
		buf.WriteString("[" + level + "] ")
	}
	buf.WriteString(msg)

	if l.Flags&(LOG_FLAG_FUNCTION|LOG_FLAG_LINE|LOG_FLAG_SHORTFILE|LOG_FLAG_LONGFILE|LOG_FLAG_RELATIVEFILE) != 0 {
		buf.WriteString(" [")
		if l.Flags&LOG_FLAG_FUNCTION != 0 {
			buf.WriteString(funcName)
		}
		if l.Flags&LOG_FLAG_LINE != 0 {
			buf.WriteString(fmt.Sprintf(":%d", line))
		}
		if l.Flags&(LOG_FLAG_SHORTFILE|LOG_FLAG_LONGFILE|LOG_FLAG_RELATIVEFILE) != 0 {
			buf.WriteString(" " + file)
		}
		buf.WriteString("]")
	}

	l.output(buf.String())
}

func (l *Logger) output(msg string) {
	if l.Output&LOG_OUTPUT_CONSOLE != 0 {
		fmt.Println(msg)
	}
	// A futuro salida a archivo, DB o API remota
}

func (l *Logger) WithContext(level LogLevel, ctx map[string]any, format string, a ...any) {
	if l.Level > level {
		return
	}

	copy := *l
	copy.Context = ctx
	go copy.outputEntry(level.String(), format, a...)
}
