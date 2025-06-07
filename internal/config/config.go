package config

import (
	"bufio"
	"fmt"
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

	SERVER_PORT         string
	SERVER_DEFAULT_LANG string
	SERVER_USE_HTTPS    bool

	DB_CONNECTION string
	DB_HOST       string
	DB_PORT       string
	DB_DATABASE   string
	DB_USERNAME   string
	DB_PASSWORD   string
	DB_URI        string

	MAIL_MAILER       string
	MAIL_HOST         string
	MAIL_PORT         int
	MAIL_USERNAME     string
	MAIL_PASSWORD     string
	MAIL_ENCRYPTION   string
	MAIL_FROM_ADDRESS string
	MAIL_FROM_NAME    string
}

// si se proporciona una ruta, se usa la primera; de lo contrario, se carga el archivo .env por defecto
var Env Environment

func Load(path ...string) error {
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
		value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)

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
		case "SERVER_DEFAULT_LANG":
			Env.SERVER_DEFAULT_LANG = value
		case "SERVER_USE_HTTPS":
			Env.SERVER_USE_HTTPS = value == "true"

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
		}
	}

	return scanner.Err()
}

// deprecado
// func Load(path ...string) error {
// 	p := ".env"
// 	if len(path) > 0 {
// 		p = path[0]
// 	}

// 	file, err := os.Open(p)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)
// 	i := 0

// 	for scanner.Scan() {
// 		i++
// 		line := strings.TrimSpace(scanner.Text())

// 		if line == "" || strings.HasPrefix(line, "#") {
// 			continue
// 		}

// 		parts := strings.SplitN(line, "=", 2)
// 		if len(parts) != 2 {
// 			return fmt.Errorf("error de sintaxis en variables de entorno en la línea %d: %v", i, line)
// 		}

// 		key := strings.TrimSpace(parts[0])
// 		value := strings.TrimSpace(parts[1])
// 		value = strings.Trim(value, `"'`)

// 		if key == "" {
// 			return fmt.Errorf("clave vacía al cargar variables de entorno en la línea %d: %v", i, line)
// 		}

// 		os.Setenv(key, value)
// 	}

// 	return scanner.Err()
// }
