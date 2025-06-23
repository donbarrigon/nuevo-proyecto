package system

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Environment struct {
	APP_KEY    string
	APP_URL    string
	APP_LOCALE string

	SERVER_PORT            string
	SERVER_HTTPS_ENABLED   bool
	SERVER_HTTPS_CERT_PATH string
	SERVER_HTTPS_KEY_PATH  string

	DB_CONNECTION string
	DB_HOST       string
	DB_PORT       string
	DB_DATABASE   string
	DB_USERNAME   string
	DB_PASSWORD   string
	DB_URI        string

	LOG_LEVEL       LogLevel
	LOG_FLAGS       int
	LOG_OUTPUT      int
	LOG_URL         string
	LOG_PATH        string
	LOG_DATE_FORMAT string

	MAIL_MAILER       string
	MAIL_SCHEME       string
	MAIL_HOST         string
	MAIL_PORT         int
	MAIL_USERNAME     string
	MAIL_PASSWORD     string
	MAIL_ENCRYPTION   string
	MAIL_FROM_ADDRESS string
	MAIL_FROM_NAME    string
}

// si se proporciona una ruta, se usa la primera; de lo contrario, se carga el archivo .env por defecto
var Env = Environment{
	APP_KEY:    "base64:AlgunacadenacodificadaenBase64aleatoria==",
	APP_URL:    "http://localhost",
	APP_LOCALE: "es",

	SERVER_PORT:            "8080",
	SERVER_HTTPS_ENABLED:   false,
	SERVER_HTTPS_CERT_PATH: "certs/server.crt",
	SERVER_HTTPS_KEY_PATH:  "certs/server.key",

	DB_CONNECTION: "mongodb",
	DB_HOST:       "127.0.0.1",
	DB_PORT:       "27017",
	DB_DATABASE:   "sample_mflix",
	DB_USERNAME:   "",
	DB_PASSWORD:   "",
	DB_URI:        "",

	LOG_LEVEL:       LOG_DEBUG,
	LOG_FLAGS:       LOG_FLAG_ALL,
	LOG_OUTPUT:      LOG_OUTPUT_CONSOLE | LOG_OUTPUT_DATABASE,
	LOG_URL:         "http://127.0.0.1/debug/log",
	LOG_PATH:        "log.json",
	LOG_DATE_FORMAT: "2006-01-02 15:04:05.000000",

	MAIL_MAILER:       "log",
	MAIL_SCHEME:       "smtp",
	MAIL_HOST:         "smtp.gmail.com",
	MAIL_PORT:         587,
	MAIL_USERNAME:     "tuemail@gmail.com",
	MAIL_PASSWORD:     "tu_contraseña_o_app_password",
	MAIL_ENCRYPTION:   "tls",
	MAIL_FROM_ADDRESS: "tuemail@gmail.com",
	MAIL_FROM_NAME:    "MiAppGo",
}

func LoadEnv(filepath ...string) error {
	f := ".env"
	if len(filepath) > 0 {
		f = filepath[0]
	}

	file, err := os.Open(f)
	if err != nil {
		return err
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
			return fmt.Errorf("error de sintaxis en variables de entorno en la línea %d: %v", i, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if idx := strings.Index(value, "#"); idx != -1 && !strings.HasPrefix(value, `"`) && !strings.HasPrefix(value, `'`) {
			value = strings.TrimSpace(value[:idx])
		}
		value = strings.Trim(value, `"'`)

		if key == "" {
			return fmt.Errorf("clave vacía al cargar variables de entorno en la línea %d: %v", i, line)
		}

		os.Setenv(key, value)

		switch key {
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

		case "DB_CONNECTION":
			Env.DB_CONNECTION = value
		case "DB_HOST":
			Env.DB_HOST = value
		case "DB_PORT":
			Env.DB_PORT = value
		case "DB_DATABASE":
			Env.DB_DATABASE = value
		case "DB_USERNAME":
			Env.DB_USERNAME = value
		case "DB_PASSWORD":
			Env.DB_PASSWORD = value
		case "DB_URI":
			Env.DB_URI = value
			if Env.DB_URI == "" {
				switch Env.DB_CONNECTION {
				case "mongodb":
					if Env.DB_USERNAME != "" && Env.DB_PASSWORD != "" {
						Env.DB_URI = fmt.Sprintf(
							"mongodb://%s:%s@%s:%s/%s?authSource=admin",
							Env.DB_USERNAME, Env.DB_PASSWORD,
							Env.DB_HOST, Env.DB_PORT, Env.DB_DATABASE,
						)
					} else {
						Env.DB_URI = fmt.Sprintf(
							"mongodb://%s:%s/%s",
							Env.DB_HOST, Env.DB_PORT, Env.DB_DATABASE,
						)
					}

				case "mysql":
					if Env.DB_USERNAME != "" && Env.DB_PASSWORD != "" {
						Env.DB_URI = fmt.Sprintf(
							"%s:%s@tcp(%s:%s)/%s?parseTime=true",
							Env.DB_USERNAME, Env.DB_PASSWORD,
							Env.DB_HOST, Env.DB_PORT, Env.DB_DATABASE,
						)
					} else {
						Env.DB_URI = fmt.Sprintf(
							"tcp(%s:%s)/%s?parseTime=true",
							Env.DB_HOST, Env.DB_PORT, Env.DB_DATABASE,
						)
					}

				case "postgresql":
					if Env.DB_USERNAME != "" && Env.DB_PASSWORD != "" {
						Env.DB_URI = fmt.Sprintf(
							"postgres://%s:%s@%s:%s/%s?sslmode=disable",
							Env.DB_USERNAME, Env.DB_PASSWORD,
							Env.DB_HOST, Env.DB_PORT, Env.DB_DATABASE,
						)
					} else {
						Env.DB_URI = fmt.Sprintf(
							"postgres://%s:%s/%s?sslmode=disable",
							Env.DB_HOST, Env.DB_PORT, Env.DB_DATABASE,
						)
					}
				default:
					Env.DB_URI = value
				}
			}

		case "LOG_LEVEL":
			switch strings.ToUpper(strings.TrimSpace(value)) {
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
				case "RELATIVEFILE":
					flags |= LOG_FLAG_RELATIVEFILE
				case "FUNCTION":
					flags |= LOG_FLAG_FUNCTION
				case "LINE":
					flags |= LOG_FLAG_LINE
				case "PREFIX":
					flags |= LOG_FLAG_PREFIX
				case "CONSOLE_AS_JSON":
					flags |= LOG_FLAG_CONSOLE_AS_JSON
				case "FILE_AS_JSON":
					flags |= LOG_FLAG_FILE_AS_JSON
				case "CONTEXT":
					flags |= LOG_FLAG_CONTEXT
				case "DUMP":
					flags |= LOG_FLAG_DUMP
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
		case "LOG_PATH":
			Env.LOG_PATH = value
		case "LOG_DATE_FORMAT":
			Env.LOG_DATE_FORMAT = value

		case "MAIL_MAILER":
			Env.MAIL_MAILER = value
		case "MAIL_SCHEME":
			Env.MAIL_SCHEME = value
		case "MAIL_HOST":
			Env.MAIL_HOST = value
		case "MAIL_PORT":
			port, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("MAIL_PORT inválido en la línea %d: %v", i, value)
			}
			Env.MAIL_PORT = port
		case "MAIL_USERNAME":
			Env.MAIL_USERNAME = value
		case "MAIL_PASSWORD":
			Env.MAIL_PASSWORD = value
		case "MAIL_ENCRYPTION":
			Env.MAIL_ENCRYPTION = value
		case "MAIL_FROM_ADDRESS":
			Env.MAIL_FROM_ADDRESS = value
		case "MAIL_FROM_NAME":
			Env.MAIL_FROM_NAME = value
		default:
			log.Printf("%v no es una variable de entrono", key)
		}
	}

	if scanner.Err() != nil {
		Log.Emergency("Fallo crítico al cargar las variables de entorno desde el archivo {file}", map[string]any{
			"env":  Env,
			"file": f,
		})
		return scanner.Err()
	}

	Log.Info("Variables de entorno cargadas exitosamente desde el archivo {file}", map[string]any{
		"env":  Env,
		"file": f,
	})
	return nil
}
