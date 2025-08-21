package app

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type V map[string]any

type Environment struct {
	APP_NAME   string
	APP_KEY    string
	APP_URL    string
	APP_LOCALE string

	SERVER_PORT             string
	SERVER_HTTPS_ENABLED    bool
	SERVER_HTTPS_CERT_PATH  string
	SERVER_HTTPS_KEY_PATH   string
	SERVER_TIMEOUT          int
	SERVER_MIGRATION_ENABLE bool

	DB_DATABASE          string
	DB_CONNECTION_STRING string

	LOG_LEVEL       LogLevel
	LOG_FLAGS       int
	LOG_OUTPUT      int
	LOG_URL         string
	LOG_URL_TOKEN   string
	LOG_PATH        string
	LOG_CHANNEL     string
	LOG_FILE_FORMAT LogFileFormat
	LOG_DAYS        int
	LOG_DATE_FORMAT string

	MAIL_MAILER     string
	MAIL_HOST       string
	MAIL_PORT       string
	MAIL_USERNAME   string
	MAIL_PASSWORD   string
	MAIL_ENCRYPTION string
	MAIL_FROM_NAME  string
}

// si se proporciona una ruta, se usa la primera; de lo contrario, se carga el archivo .env por defecto
var Env = Environment{
	APP_NAME:   "MiAppGo",
	APP_KEY:    "base64:AlgunacadenacodificadaenBase64aleatoria==",
	APP_URL:    "http://localhost",
	APP_LOCALE: "es",

	SERVER_PORT:             "8080",
	SERVER_HTTPS_ENABLED:    false,
	SERVER_HTTPS_CERT_PATH:  "certs/server.crt",
	SERVER_HTTPS_KEY_PATH:   "certs/server.key",
	SERVER_TIMEOUT:          60,
	SERVER_MIGRATION_ENABLE: false,

	DB_DATABASE:          "sample_mflix",
	DB_CONNECTION_STRING: "mongodb://localhost:27017",

	LOG_LEVEL:       LOG_DEBUG,
	LOG_FLAGS:       LOG_FLAG_ALL,
	LOG_OUTPUT:      LOG_OUTPUT_CONSOLE | LOG_OUTPUT_FILE | LOG_OUTPUT_DATABASE | LOG_OUTPUT_REMOTE,
	LOG_URL:         "http://127.0.0.1/debug/log",
	LOG_URL_TOKEN:   "",
	LOG_PATH:        "log",
	LOG_CHANNEL:     "daily",
	LOG_DAYS:        14,
	LOG_DATE_FORMAT: "2006-01-02 15:04:05.000000",

	MAIL_MAILER:     "smtp",
	MAIL_HOST:       "smtp.gmail.com",
	MAIL_PORT:       "587",
	MAIL_USERNAME:   "tuemail@gmail.com",
	MAIL_PASSWORD:   "tu_contraseña_o_app_password",
	MAIL_ENCRYPTION: "tls",
	MAIL_FROM_NAME:  "MiAppGo",
}

func LoadEnv(filepath ...string) {
	f := ".env"
	if len(filepath) > 0 {
		f = filepath[0]
	}

	file, err := os.Open(f)
	if err != nil {
		PrintWarning("Environment file (.env) not found at location {file}", Entry{"file", f})
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0

	for scanner.Scan() {
		i++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			PrintWarning("Syntax error in environment variables at line {lineNumber}: {rawLine}",
				Entry{"lineNumber", i},
				Entry{"rawLine", line})
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if idx := strings.Index(value, "#"); idx != -1 && !strings.HasPrefix(value, `"`) && !strings.HasPrefix(value, `'`) {
			value = strings.TrimSpace(value[:idx])
		}
		value = strings.Trim(value, `"'`)

		if key == "" {
			PrintWarning("Empty key detected while loading environment variables at line {lineNumber}: {rawLine}",
				Entry{"lineNumber", i},
				Entry{"rawLine", line},
			)
			continue
		}

		os.Setenv(key, value)

		switch key {
		case "APP_NAME":
			Env.APP_NAME = value
		case "APP_KEY":
			Env.APP_KEY = value
		case "APP_URL":
			Env.APP_URL = value
		case "APP_LOCALE":
			Env.APP_LOCALE = value

		case "SERVER_PORT":
			Env.SERVER_PORT = value
		case "SERVER_HTTPS_ENABLED":
			Env.SERVER_HTTPS_ENABLED = false
			if strings.ToLower(value) == "true" {
				Env.SERVER_HTTPS_ENABLED = true
			}
		case "SERVER_HTTPS_CERT_PATH":
			Env.SERVER_HTTPS_CERT_PATH = value
		case "SERVER_HTTPS_KEY_PATH":
			Env.SERVER_HTTPS_KEY_PATH = value
		case "SERVER_TIMEOUT":
			timeout, e := strconv.Atoi(value)
			if e != nil {
				Env.SERVER_TIMEOUT = timeout
			}
		case "SERVER_MIGRATION_ENABLE":
			Env.SERVER_MIGRATION_ENABLE = false
			if strings.ToLower(value) == "true" {
				Env.SERVER_MIGRATION_ENABLE = true
			}

		case "DB_DATABASE":
			Env.DB_DATABASE = value
		case "DB_CONNECTION_STRING":
			Env.DB_CONNECTION_STRING = value

		case "LOG_LEVEL":
			switch strings.ToUpper(strings.TrimSpace(value)) {
			case "OFF":
				Env.LOG_LEVEL = LOG_OFF
			case "EMERGENCY":
				Env.LOG_LEVEL = LOG_EMERGENCY
			case "ALERT":
				Env.LOG_LEVEL = LOG_ALERT
			case "CRITICAL":
				Env.LOG_LEVEL = LOG_CRITICAL
			case "ERROR":
				Env.LOG_LEVEL = LOG_ERROR
			case "WARNING":
				Env.LOG_LEVEL = LOG_WARNING
			case "NOTICE":
				Env.LOG_LEVEL = LOG_NOTICE
			case "INFO":
				Env.LOG_LEVEL = LOG_INFO
			case "DEBUG":
				Env.LOG_LEVEL = LOG_DEBUG
			default:
				Env.LOG_LEVEL = LOG_DEBUG
			}
		case "LOG_FLAGS":
			flags := 0
			parts := strings.Split(value, ",")
			for _, part := range parts {
				switch strings.ToUpper(strings.TrimSpace(part)) {
				case "TIMESTAMP":
					flags |= LOG_FLAG_TIMESTAMP
				case "LONGFILE":
					flags |= LOG_FLAG_LONGFILE
				case "SHORTFILE":
					flags |= LOG_FLAG_SHORTFILE
				case "FUNCTION":
					flags |= LOG_FLAG_FUNCTION
				case "LINE":
					flags |= LOG_FLAG_LINE
				case "PREFIX":
					flags |= LOG_FLAG_PREFIX
				case "CONSOLE_AS_JSON":
					flags |= LOG_FLAG_CONSOLE_AS_JSON
				case "CONTEXT":
					flags |= LOG_FLAG_CONTEXT
				}
			}
			Env.LOG_FLAGS = flags
		case "LOG_OUTPUT":
			outputs := 0
			parts := strings.Split(value, ",")
			for _, part := range parts {
				switch strings.ToUpper(strings.TrimSpace(part)) {
				case "CONSOLE":
					outputs |= LOG_OUTPUT_CONSOLE
				case "FILE":
					outputs |= LOG_OUTPUT_FILE
				case "DATABASE":
					outputs |= LOG_OUTPUT_DATABASE
				case "REMOTE":
					outputs |= LOG_OUTPUT_REMOTE
				}
			}
			Env.LOG_OUTPUT = outputs
		case "LOG_URL":
			Env.LOG_URL = value
		case "LOG_URL_TOKEN":
			Env.LOG_URL_TOKEN = value
		case "LOG_PATH":
			Env.LOG_PATH = value
		case "LOG_CHANNEL":
			value = strings.ToLower(value)
			if value == "monthly" || value == "weekly" || value == "single" {
				Env.LOG_CHANNEL = value
			} else {
				Env.LOG_CHANNEL = "daily"
			}
		case "LOG_DAYS":
			days, err := strconv.Atoi(value)
			if err != nil {
				PrintWarning("Invalid LOG_DAYS value at line {lineNumber}: {value}",
					Entry{"lineNumber", i},
					Entry{"value", value},
				)
				continue
			}
			Env.LOG_DAYS = days
		case "LOG_DATE_FORMAT":
			Env.LOG_DATE_FORMAT = value
		case "MAIL_MAILER":
			Env.MAIL_MAILER = value
		case "MAIL_HOST":
			Env.MAIL_HOST = value
		case "MAIL_PORT":
			Env.MAIL_PORT = value
		case "MAIL_USERNAME":
			Env.MAIL_USERNAME = value
		case "MAIL_PASSWORD":
			Env.MAIL_PASSWORD = value
		case "MAIL_ENCRYPTION":
			Env.MAIL_ENCRYPTION = value
		case "MAIL_FROM_NAME":
			Env.MAIL_FROM_NAME = value
		default:
			PrintWarning("{envKey} is not a valid environment variable name",
				Entry{"envKey", key},
			)
		}
	}

	if scanner.Err() != nil {
		PrintWarning("Critical failure loading environment variables from file {file}\nerror: {error}",
			Entry{"file", f},
			Entry{"error", scanner.Err().Error()},
			Entry{"env", Env},
		)
		return
	}

	PrintInfo("Environment variables loaded successfully from file {file}",
		Entry{"file", f},
		Entry{"env", Env},
	)
}
