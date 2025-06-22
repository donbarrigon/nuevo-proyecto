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
	APP_ENV    string
	APP_DEBUG  bool
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
	APP_ENV:    "local",
	APP_DEBUG:  true,
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

func LoadEnv(path ...string) error {
	p := ".env"
	if len(path) > 0 {
		p = path[0]
	}

	file, err := os.Open(p)
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
		case "APP_ENV":
			Env.APP_ENV = value
		case "APP_DEBUG":
			Env.APP_DEBUG = value == "true"
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

	printEnv()

	return scanner.Err()
}

func printEnv() {
	fmt.Println(".env")
	fmt.Println("--------------------------------")
	fmt.Printf("APP_ENV: %v\n", Env.APP_ENV)
	fmt.Printf("APP_DEBUG: %v\n", Env.APP_DEBUG)
	fmt.Printf("APP_KEY: %v\n", Env.APP_KEY)
	fmt.Printf("APP_URL: %v\n", Env.APP_URL)
	fmt.Printf("APP_LOCALE: %v\n", Env.APP_LOCALE)

	fmt.Printf("SERVER_PORT: %v\n", Env.SERVER_PORT)
	fmt.Printf("SERVER_HTTPS_ENABLED: %v\n", Env.SERVER_HTTPS_ENABLED)
	fmt.Printf("SERVER_HTTPS_CERT_PATH: %v\n", Env.SERVER_HTTPS_CERT_PATH)
	fmt.Printf("SERVER_HTTPS_KEY_PATH: %v\n", Env.SERVER_HTTPS_KEY_PATH)

	fmt.Printf("DB_CONNECTION: %v\n", Env.DB_CONNECTION)
	fmt.Printf("DB_HOST: %v\n", Env.DB_HOST)
	fmt.Printf("DB_PORT: %v\n", Env.DB_PORT)
	fmt.Printf("DB_DATABASE: %v\n", Env.DB_DATABASE)
	fmt.Printf("DB_USERNAME: %v\n", Env.DB_USERNAME)
	fmt.Printf("DB_PASSWORD: %v\n", Env.DB_PASSWORD)
	fmt.Printf("DB_URI: %v\n", Env.DB_URI)

	fmt.Printf("MAIL_MAILER: %v\n", Env.MAIL_MAILER)
	fmt.Printf("MAIL_SCHEME: %v\n", Env.MAIL_SCHEME)
	fmt.Printf("MAIL_HOST: %v\n", Env.MAIL_HOST)
	fmt.Printf("MAIL_PORT: %v\n", Env.MAIL_PORT)
	fmt.Printf("MAIL_USERNAME: %v\n", Env.MAIL_USERNAME)
	fmt.Printf("MAIL_PASSWORD: %v\n", Env.MAIL_PASSWORD)
	fmt.Printf("MAIL_ENCRYPTION: %v\n", Env.MAIL_ENCRYPTION)
	fmt.Printf("MAIL_FROM_ADDRESS: %v\n", Env.MAIL_FROM_ADDRESS)
	fmt.Printf("MAIL_FROM_NAME: %v\n", Env.MAIL_FROM_NAME)
	fmt.Println("--------------------------------")
}
