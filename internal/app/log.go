package app

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"gopkg.in/yaml.v3"
)

type LogLevel int
type LogFileFormat int

type Logger struct {
	ID       string   `json:"id,omEntrypty" yaml:"id,omEntrypty" id:"time,omEntrypty"`
	Time     string   `json:"time,omEntrypty" yaml:"time,omEntrypty" xml:"time,omEntrypty"`
	Level    LogLevel `json:"level,omEntrypty" yaml:"level,omEntrypty" xml:"level,omEntrypty"`
	Message  string   `json:"message" yaml:"message" xml:"message"`
	Function string   `json:"function,omEntrypty" yaml:"function,omEntrypty" xml:"function,omEntrypty"`
	Line     string   `json:"line,omEntrypty" yaml:"line,omEntrypty" xml:"line,omEntrypty"`
	File     string   `json:"file,omEntrypty" yaml:"file,omEntrypty" xml:"file,omEntrypty"`
	Context  Object   `json:"context,omEntrypty" yaml:"context,omEntrypty" xml:"context,omEntrypty"`
}

const (
	LOG_OFF LogLevel = iota // 0 - Desactiva todos los logs

	LOG_EMERGENCY // 1 - El sistema está inutilizable
	LOG_ALERT     // 2 - Se necesita acción inmediata
	LOG_CRITICAL  // 3 - Fallo crítico del sistema
	LOG_ERROR     // 4 - Errores de ejecución
	LOG_WARNING   // 5 - Algo inesperado pasó
	LOG_NOTICE    // 6 - Eventos normales, pero significativos
	LOG_INFO      // 7 - Información general
	LOG_DEBUG     // 8 - Información detallada para depuración
	LOG_PRINT     // 9 - Solo imprime en consola
)

const (
	LOG_FLAG_TIMESTAMP       = 1 << iota // 1     - Agrega la fecha y hora formateada según LOG_DATE_FORMAT
	LOG_FLAG_LONGFILE                    // 2     - Ruta completa del archivo y número de línea: /a/b/c/d.go:23
	LOG_FLAG_SHORTFILE                   // 4     - Solo el nombre del archivo y línea: d.go:23
	LOG_FLAG_FUNCTION                    // 8    - Nombre de la función desde donde se llamó
	LOG_FLAG_LINE                        // 16    - Solo el número de línea (sin ruta de archivo)
	LOG_FLAG_PREFIX                      // 32    - Agrega un prefijo antes del mensaje (por ejemplo: [DEBUG])
	LOG_FLAG_CONSOLE_AS_JSON             // 64   - Salida en formato JSON en la consola
	LOG_FLAG_CONSOLE_COLOR               // 128   - salida en consola con solor segun el lv
	LOG_FLAG_CONTEXT                     // 256   - Agrega el contexto de la petición al log
	LOG_FLAG_ID                          // 512  - Genera un ID único en formato hexadecimal string (bson.ObjectID.Hex())

	// Combinación de todos los flags
	LOG_FLAG_ALL = LOG_FLAG_TIMESTAMP |
		LOG_FLAG_LONGFILE |
		LOG_FLAG_SHORTFILE |
		LOG_FLAG_FUNCTION |
		LOG_FLAG_LINE |
		LOG_FLAG_PREFIX |
		// LOG_FLAG_CONSOLE_AS_JSON |
		LOG_FLAG_CONTEXT
	// LOG_FLAG_ID
)

const (
	LOG_OUTPUT_CONSOLE  = 1 << iota // 1 - salida por consola estándar
	LOG_OUTPUT_FILE                 // 2 - salida a archivo
	LOG_OUTPUT_DATABASE             // 4 - guardar logs en base de datos
	LOG_OUTPUT_REMOTE               // 8 - enviar a un servidor remoto (opcional)
)

const (
	LOG_FILE_FORMAT_NDJSON LogFileFormat = iota // 0 - NDJSON (JSON por línea)
	LOG_FILE_FORMAT_CSV                         // 1 - CSV (valores separados por coma)
	LOG_FILE_FORMAT_PLAIN                       // 2 - Texto plano
	LOG_FILE_FORMAT_XML                         // 3 - XML estructurado
	LOG_FILE_FORMAT_YAML                        // 4 - YAML legible para humanos
	LOG_FILE_FORMAT_LTSV                        // 5 - LTSV (Labelled Tab-separated Values)
)

func (lv LogLevel) String() string {
	switch lv {
	case LOG_OFF:
		return "OFF"
	case LOG_EMERGENCY:
		return "EMERGENCY"
	case LOG_ALERT:
		return "ALERT"
	case LOG_CRITICAL:
		return "CRITICAL"
	case LOG_ERROR:
		return "ERROR"
	case LOG_WARNING:
		return "WARNING"
	case LOG_NOTICE:
		return "NOTICE"
	case LOG_INFO:
		return "INFO"
	case LOG_DEBUG:
		return "DEBUG"
	case LOG_PRINT:
		return "PRINT"
	default:
		return "UNKNOWN"
	}
}

func (lv LogLevel) Color() string {
	switch lv {
	case LOG_EMERGENCY:
		return "\033[91m" // rojo brillante
	case LOG_ALERT:
		return "\033[95m" // magenta
	case LOG_CRITICAL:
		return "\033[35m" // fucsia
	case LOG_ERROR:
		return "\033[31m" // rojo
	case LOG_WARNING:
		return "\033[33m" // amarillo
	case LOG_NOTICE:
		return "\033[92m" // verde claro
	case LOG_INFO:
		return "\033[34m" // azul
	case LOG_DEBUG:
		return "\033[36m" // cian
	case LOG_PRINT:
		return "\033[90m" // gris claro
	default:
		return "\033[0m"
	}
}

func (lv LogLevel) DefaultColor() string {
	return "\033[0m"
}

func (l LogLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.String())
}

func (l LogLevel) MarshalYAML() (interface{}, error) {
	return l.String(), nil
}

func (l LogLevel) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(l.String(), start)
}

func (f LogFileFormat) String() string {
	switch f {
	case LOG_FILE_FORMAT_NDJSON:
		return "ndjson"
	case LOG_FILE_FORMAT_CSV:
		return "csv"
	case LOG_FILE_FORMAT_PLAIN:
		return "plain"
	case LOG_FILE_FORMAT_XML:
		return "xml"
	case LOG_FILE_FORMAT_YAML:
		return "yaml"
	case LOG_FILE_FORMAT_LTSV:
		return "ltsv"
	default:
		return "unknown"
	}
}

func PrintEmergency(msg string, ctx ...Entry) {
	if Env.LOG_LEVEL >= LOG_EMERGENCY {
		l := &Logger{
			Level:   LOG_EMERGENCY,
			Message: msg,
			Context: ctx,
		}
		go l.output()
	}
}

func PrintAlert(msg string, ctx ...Entry) {
	if Env.LOG_LEVEL >= LOG_ALERT {
		l := &Logger{
			Level:   LOG_ALERT,
			Message: msg,
			Context: ctx,
		}
		go l.output()
	}
}

func PrintCritical(msg string, ctx ...Entry) {
	if Env.LOG_LEVEL >= LOG_CRITICAL {
		l := &Logger{
			Level:   LOG_CRITICAL,
			Message: msg,
			Context: ctx,
		}
		go l.output()
	}
}

func PrintError(msg string, ctx ...Entry) {
	if Env.LOG_LEVEL >= LOG_ERROR {
		l := &Logger{
			Level:   LOG_ERROR,
			Message: msg,
			Context: ctx,
		}
		go l.output()
	}
}

func PrintWarning(msg string, ctx ...Entry) {
	if Env.LOG_LEVEL >= LOG_WARNING {
		l := &Logger{
			Level:   LOG_WARNING,
			Message: msg,
			Context: ctx,
		}
		go l.output()
	}
}

func PrintNotice(msg string, ctx ...Entry) {
	if Env.LOG_LEVEL >= LOG_NOTICE {
		l := &Logger{
			Level:   LOG_NOTICE,
			Message: msg,
			Context: ctx,
		}
		go l.output()
	}
}

func PrintInfo(msg string, ctx ...Entry) {
	if Env.LOG_LEVEL >= LOG_INFO {
		l := &Logger{
			Level:   LOG_INFO,
			Message: msg,
			Context: ctx,
		}
		go l.output()
	}
}

func PrintDebug(msg string, ctx ...Entry) {
	if Env.LOG_LEVEL >= LOG_DEBUG {
		l := &Logger{
			Level:   LOG_DEBUG,
			Message: msg,
			Context: ctx,
		}
		go l.output()
	}
}

func PrintLog(level LogLevel, msg string, ctx ...Entry) {
	if Env.LOG_LEVEL >= level {
		l := &Logger{
			Level:   LOG_PRINT,
			Message: msg,
			Context: ctx,
		}
		go l.output()
	}
}

func Print(msg string, ctx ...Entry) {
	l := &Logger{
		Level:   LOG_PRINT,
		Message: msg,
		Context: ctx,
	}
	go l.output()
}

func (l *Logger) Dump(a any) {
	fmt.Println(l.formatDump(a))
}

func (l *Logger) DumpMany(vars ...any) {
	sep := strings.Repeat("-", 30)

	for i, v := range vars {
		if i > 0 {
			fmt.Println(sep)
		}
		fmt.Println(l.formatDump(v))
	}
}

func (l *Logger) output() {
	// Obtener información del runtime
	pc, file, line, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()

	// Si SHORTFILE está activado
	if Env.LOG_FLAGS&LOG_FLAG_SHORTFILE != 0 {
		file = filepath.Base(file)
	}

	// Preparar mensaje
	l.Message = InterpolatePlaceholders(l.Message, l.Context...)

	if Env.LOG_FLAGS&LOG_FLAG_ID != 0 {
		l.ID = bson.NewObjectID().Hex()
	}

	if Env.LOG_FLAGS&LOG_FLAG_TIMESTAMP != 0 {
		now := time.Now().Format(Env.LOG_DATE_FORMAT)
		l.Time = now
	}

	if Env.LOG_FLAGS&LOG_FLAG_FUNCTION != 0 {
		l.Function = funcName
	}

	if Env.LOG_FLAGS&LOG_FLAG_LINE != 0 {
		l.Line = strconv.Itoa(line)
	}

	if Env.LOG_FLAGS&(LOG_FLAG_LONGFILE|LOG_FLAG_SHORTFILE) != 0 {
		l.File = file
	}

	//if Env.LOG_FLAGS&LOG_FLAG_CONTEXT != 0 && ctx != nil {
	// l.Context = ctx
	//}

	if Env.LOG_OUTPUT&LOG_OUTPUT_CONSOLE != 0 || l.Level == LOG_PRINT {
		l.outputConsole()
		if l.Level == LOG_PRINT {
			return
		}
	}

	if Env.LOG_OUTPUT&LOG_OUTPUT_FILE != 0 {
		l.outputFile()
	}

	if Env.LOG_OUTPUT&LOG_OUTPUT_DATABASE != 0 {
		l.outputDatabase()
	}

	if Env.LOG_OUTPUT&LOG_OUTPUT_REMOTE != 0 {
		l.outputRemote()
	}

}

func (l *Logger) outputConsole() {
	if Env.LOG_FLAGS&LOG_FLAG_CONSOLE_AS_JSON != 0 {
		data, _ := json.MarshalIndent(l, "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Println(l.outputPlain(true))
	}
}

func (l *Logger) outputFile() {
	file := l.openFile()
	if file == nil {
		return
	}
	defer file.Close()

	l.deleteOldFiles()

	var output string

	switch Env.LOG_FILE_FORMAT {
	case LOG_FILE_FORMAT_NDJSON:
		output = l.outputNDJSON()
	case LOG_FILE_FORMAT_CSV:
		output = l.outputCSV() // CSV: Time, Level, Message, Function, File, Line, context
	case LOG_FILE_FORMAT_PLAIN:
		output = l.outputPlain(false)
	case LOG_FILE_FORMAT_XML:
		output = l.outputXML()
	case LOG_FILE_FORMAT_YAML:
		output = l.outputYAML()
	case LOG_FILE_FORMAT_LTSV:
		output = l.outputLTSV()
	default:
		// Fallback a ndjson
		output = l.outputNDJSON()
	}

	file.WriteString(output + "\n")
}

func (l *Logger) outputDatabase() {

	// se queda sin funcionar hasta que resuelva lo del ciclo de dependencias.

	// id, e := bson.ObjectIDFromHex(l.ID)
	// if e != nil {
	// 	Log.Print("Failed to convert string [:input_id] to ObjectID :error ",
	// 		F{"error", e.Error()},
	// 		F{"input_id", l.ID},
	// 	)
	// }

	// ctx := make(map[string]string, len(l.Context))
	// for _, field := range l.Context {
	// 	ctx[field.Key] = fmt.Sprint(field.Value)
	// }

	// m := &model.Log{
	// 	ID:       id,
	// 	Time:     l.Time,
	// 	Level:    l.Level.String(),
	// 	Message:  l.Message,
	// 	Function: l.Function,
	// 	Line:     l.Line,
	// 	File:     l.File,
	// 	Context:  ctx,
	// }

	// if err := db.Create(m); err != nil {
	// 	Log.Print("Failed to create log in database: :error", F{"error", err.Error()})
	// }
}

func (l *Logger) outputRemote() {
	if Env.LOG_URL == "" || Env.LOG_URL_TOKEN == "" {
		return // No hay configuración completa
	}

	// Convertir el log a JSON
	jsonData, err := json.Marshal(l)
	if err != nil {
		PrintError("Failed to marshal log for remote output",
			Entry{"error", err.Error()},
			Entry{"log", l},
		)
		return
	}

	// Configuración de reintentos
	maxRetries := 3
	initialDelay := time.Second * 1
	maxDelay := time.Second * 10

	var lastError error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Esperar con backoff exponencial antes de reintentar
			delay := initialDelay * time.Duration(math.Pow(2, float64(attempt-1)))
			if delay > maxDelay {
				delay = maxDelay
			}
			time.Sleep(delay)
		}

		// Crear nueva solicitud para cada intento
		req, err := http.NewRequest("POST", Env.LOG_URL, bytes.NewBuffer(jsonData))
		if err != nil {
			lastError = err
			continue
		}

		// Configurar headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+Env.LOG_URL_TOKEN)

		// Configurar timeout (10 segundos)
		client := &http.Client{
			Timeout: time.Second * 10,
		}

		// Enviar la solicitud
		resp, err := client.Do(req)
		if err != nil {
			lastError = err
			continue
		}

		// Verificar respuesta
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			resp.Body.Close()
			return // Éxito, salir del bucle
		}

		// Si la respuesta no fue exitosa
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		lastError = fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))

		// No reintentar en errores 4xx (excepto 429 - Too Many Requests)
		if resp.StatusCode >= 400 && resp.StatusCode < 500 && resp.StatusCode != 429 {
			break
		}
	}

	// Si llegamos aquí, todos los intentos fallaron
	Print("Failed to send log to remote server after retries",
		Entry{"error", lastError.Error()},
		Entry{"url", Env.LOG_URL},
		Entry{"attempts", maxRetries},
	)
}

func (l *Logger) outputPlain(withColor bool) string {
	var b strings.Builder
	color := ""
	reset := ""
	if Env.LOG_FLAGS&LOG_FLAG_CONSOLE_COLOR != 0 && withColor {
		color = l.Level.Color()
		reset = l.Level.DefaultColor()
	}
	if Env.LOG_FLAGS&LOG_FLAG_ID != 0 {
		b.WriteString(fmt.Sprintf("[ID:%s%s%s] ", color, l.ID, reset))
	}

	if Env.LOG_FLAGS&LOG_FLAG_TIMESTAMP != 0 {
		b.WriteString(fmt.Sprintf("%s ", l.Time))
	}
	if Env.LOG_FLAGS&LOG_FLAG_PREFIX != 0 {
		b.WriteString(fmt.Sprintf("[%s%s%s] ", color, l.Level.String(), reset))
	}

	if Env.LOG_FLAGS&LOG_FLAG_CONSOLE_COLOR != 0 {
		b.WriteString(color + l.Message + reset)
	} else {
		b.WriteString(l.Message)
	}

	if Env.LOG_FLAGS&LOG_FLAG_FUNCTION != 0 {
		b.WriteString(fmt.Sprintf(" [%s]", l.Function))
	}
	if Env.LOG_FLAGS&(LOG_FLAG_LONGFILE|LOG_FLAG_SHORTFILE|LOG_FLAG_LINE) != 0 {
		b.WriteString(fmt.Sprintf(" (%s:%s)", l.File, l.Line))
	}

	return b.String()
}

func (l *Logger) outputNDJSON() string {
	jsonData, err := json.Marshal(l)
	var output string
	if err != nil {
		msg := Translate(Env.APP_LOCALE, "Log serialization error: {error}", Entry{"error", err.Error()})
		escapedDump := strings.ReplaceAll(l.formatDump(l), `"`, `\"`)
		escapedDump = strings.ReplaceAll(escapedDump, "\n", " ")
		escapedDump = strings.ReplaceAll(escapedDump, "\r", " ")
		output = InterpolatePlaceholders(`{"level":"ERROR","message":"{msg}","context":"{context}"}`,
			Entry{"msg", msg},
			Entry{"context", escapedDump},
		)

		Print(Translate(Env.APP_LOCALE, msg, Entry{"context", l}))
	} else {
		output = string(jsonData)
	}
	return output
}

func (l *Logger) outputCSV() string {
	var record []string

	if Env.LOG_FLAGS&LOG_FLAG_ID != 0 {
		record = append(record, l.ID)
	}

	if Env.LOG_FLAGS&LOG_FLAG_TIMESTAMP != 0 {
		record = append(record, l.Time)
	}

	if Env.LOG_FLAGS&LOG_FLAG_PREFIX != 0 {
		record = append(record, l.Level.String())
	}

	// El mensaje siempre va
	record = append(record, l.Message)

	if Env.LOG_FLAGS&LOG_FLAG_FUNCTION != 0 {
		record = append(record, l.Function)
	}

	if Env.LOG_FLAGS&(LOG_FLAG_LONGFILE|LOG_FLAG_SHORTFILE) != 0 {
		record = append(record, l.File)
	}

	if Env.LOG_FLAGS&(LOG_FLAG_LINE) != 0 {
		record = append(record, l.Line)
	}

	if Env.LOG_FLAGS&LOG_FLAG_CONTEXT != 0 && len(l.Context) > 0 {
		dump := l.formatDump(l.Context)
		dump = strings.ReplaceAll(dump, "\n", " ")
		dump = strings.ReplaceAll(dump, "\r", " ")
		record = append(record, dump)
	}

	var b strings.Builder
	writer := csv.NewWriter(&b)
	writer.Write(record)
	writer.Flush()

	return strings.TrimSpace(b.String())
}

func (l *Logger) outputXML() string {
	xmlData, err := xml.MarshalIndent(l, "", "  ")
	if err != nil {
		xmlEscape := func(s string) string {
			s = strings.ReplaceAll(s, "\n", " ")
			s = strings.ReplaceAll(s, "\r", " ")
			s = strings.ReplaceAll(s, "&", "&amp;")
			s = strings.ReplaceAll(s, "<", "&lt;")
			s = strings.ReplaceAll(s, ">", "&gt;")
			s = strings.ReplaceAll(s, `"`, "&quot;")
			s = strings.ReplaceAll(s, `'`, "&apos;")
			return s
		}
		return InterpolatePlaceholders(
			`<log><level>ERROR</level><message>Log serialization error: {error}</message><context>{context}</context></log>`,
			Entry{"error", xmlEscape(err.Error())},
			Entry{"context", xmlEscape(l.formatDump(l))},
		)
	} else {
		return string(xmlData)
	}
}

func (l *Logger) outputYAML() string {
	yamlData, err := yaml.Marshal(l)
	if err != nil {
		escapedDump := strings.ReplaceAll(l.formatDump(l), `"`, `\"`)
		escapedDump = strings.ReplaceAll(escapedDump, "\n", " ")
		escapedDump = strings.ReplaceAll(escapedDump, "\r", " ")

		return InterpolatePlaceholders(
			"level: ERROR\nmessage: Log serialization error: {error}\ncontext: \"{context}\"",
			Entry{"error", err.Error()},
			Entry{"context", escapedDump},
		)
	} else {
		return string(yamlData)
	}
}

func (l *Logger) outputLTSV() string {
	escape := func(s string) string {
		s = strings.ReplaceAll(s, "\t", " ")
		s = strings.ReplaceAll(s, "\n", " ")
		s = strings.ReplaceAll(s, "\r", " ")
		return s
	}

	var b strings.Builder

	if Env.LOG_FLAGS&LOG_FLAG_ID != 0 {
		b.WriteString("id:" + escape(l.ID) + "\t")
	}

	if Env.LOG_FLAGS&LOG_FLAG_TIMESTAMP != 0 {
		b.WriteString("time:" + escape(l.Time) + "\t")
	}

	b.WriteString("level:" + escape(l.Level.String()) + "\t")
	b.WriteString("message:" + escape(l.Message) + "\t")

	if Env.LOG_FLAGS&LOG_FLAG_FUNCTION != 0 {
		b.WriteString("function:" + escape(l.Function) + "\t")
	}

	if Env.LOG_FLAGS&(LOG_FLAG_LONGFILE|LOG_FLAG_SHORTFILE) != 0 {
		b.WriteString("file:" + escape(l.File) + "\t")
	}

	if Env.LOG_FLAGS&LOG_FLAG_LINE != 0 {
		b.WriteString("line:" + escape(l.Line) + "\t")
	}

	if Env.LOG_FLAGS&LOG_FLAG_CONTEXT != 0 && len(l.Context) > 0 {
		b.WriteString("context:" + escape(l.formatDump(l.Context)) + "\t")
	}

	// Eliminar el tab final si existe
	output := b.String()
	if len(output) > 0 && output[len(output)-1] == '\t' {
		output = output[:len(output)-1]
	}

	return output
}

func (l *Logger) openFile() *os.File {
	var filename string
	now := time.Now()

	switch strings.ToLower(Env.LOG_CHANNEL) {
	case "daily":
		filename = now.Format("2006-01-02") + ".log"
	case "monthly", "mensual":
		filename = now.Format("2006-01") + ".log"
	case "weekly":
		year, week := now.ISOWeek()
		filename = fmt.Sprintf("%d-W%02d.log", year, week)
	default:
		filename = "output.log"
	}

	if err := os.MkdirAll(Env.LOG_PATH, os.ModePerm); err != nil {
		Print("No se pudo crear el directorio de logs: {error}\n", Entry{"error", err})
		return nil
	}

	filePath := filepath.Join(Env.LOG_PATH, filename)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		Print("Failed to create log directory: {error}\n", Entry{"error", err})
		return nil
	}
	return file
}

func (l *Logger) deleteOldFiles() {
	if Env.LOG_CHANNEL != "single" {
		now := time.Now()
		// Eliminar archivos
		if Env.LOG_CHANNEL == "daily" && Env.LOG_DAYS > 0 {
			entries, _ := os.ReadDir(Env.LOG_PATH)
			cutoff := now.AddDate(0, 0, -Env.LOG_DAYS)
			for _, Entry := range entries {
				if Entry.IsDir() {
					continue
				}
				name := Entry.Name()
				if !strings.HasSuffix(name, ".log") {
					continue
				}
				datePart := strings.TrimSuffix(name, ".log")
				EntryDate, err := time.Parse("2006-01-02", datePart)
				if err == nil && EntryDate.Before(cutoff) {
					_ = os.Remove(filepath.Join(Env.LOG_PATH, name))
				}
			}
		}

		if Env.LOG_CHANNEL == "weekly" && Env.LOG_DAYS > 0 {
			entries, _ := os.ReadDir(Env.LOG_PATH)

			// Calcular semanas a conservar, redondeando hacia arriba (mínimo 1)
			weeksToKeep := (Env.LOG_DAYS + 6) / 7
			if weeksToKeep < 1 {
				weeksToKeep = 1
			}

			// Crear Objecta de semanas válidas (formato YYYY-Www)
			validWeeks := make(map[string]bool)
			for i := 0; i < weeksToKeep; i++ {
				weekTime := now.AddDate(0, 0, -7*i)
				year, week := weekTime.ISOWeek()
				weekStr := fmt.Sprintf("%d-W%02d", year, week)
				validWeeks[weekStr] = true
			}

			// Eliminar logs fuera del rango de semanas válidas
			for _, Entry := range entries {
				if Entry.IsDir() || !strings.HasSuffix(Entry.Name(), ".log") {
					continue
				}
				name := strings.TrimSuffix(Entry.Name(), ".log")

				// Formato semanal esperado: YYYY-Wxx
				if strings.Count(name, "-") == 1 && strings.Contains(name, "W") && len(name) == 8 {
					if !validWeeks[name] {
						_ = os.Remove(filepath.Join(Env.LOG_PATH, Entry.Name()))
					}
				}
			}
		}

		if strings.ToLower(Env.LOG_CHANNEL) == "monthly" && Env.LOG_DAYS > 0 {
			entries, _ := os.ReadDir(Env.LOG_PATH)

			// Redondear hacia arriba los días a meses (mínimo 1)
			monthsToKeep := (Env.LOG_DAYS + 29) / 30
			if monthsToKeep < 1 {
				monthsToKeep = 1
			}

			// Generar meses válidos
			validMonths := make(map[string]bool)
			for i := 0; i < monthsToKeep; i++ {
				month := now.AddDate(0, -i, 0).Format("2006-01")
				validMonths[month] = true
			}

			// Eliminar archivos fuera del rango permitido
			for _, Entry := range entries {
				if Entry.IsDir() || !strings.HasSuffix(Entry.Name(), ".log") {
					continue
				}
				name := strings.TrimSuffix(Entry.Name(), ".log")

				// Formato YYYY-MM
				if len(name) == 7 && strings.Count(name, "-") == 1 {
					if !validMonths[name] {
						_ = os.Remove(filepath.Join(Env.LOG_PATH, Entry.Name()))
					}
				}
			}
		}
	}
}
