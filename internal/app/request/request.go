package request

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/pkg/errors"
	"github.com/donbarrigon/nuevo-proyecto/pkg/lang"
	"golang.org/x/exp/constraints"
)

func Validate(l string, req any, rules map[string][]string) errors.Error {
	err := errors.NewError()

	val := reflect.ValueOf(req)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return errors.SUnknown("La request no es un struct válido")
	}

	typ := val.Type()

	for key, r := range rules {
		found := false
		isRequired := false
		isNullable := false
		var value reflect.Value

		numFields := typ.NumField()
		for i := 0; i < numFields; i++ {
			field := typ.Field(i)
			tag := field.Tag.Get("json")
			tagName := strings.Split(tag, ",")[0]
			if tagName == key {
				value = val.Field(i)
				found = true
				break
			}
		}

		for _, rule := range r {
			if rule == "required" {
				isRequired = true
				break
			}
			if rule == "nullable" {
				isNullable = true
				break
			}
		}

		if !found {
			if isRequired {
				err.Append(key, lang.TT(l, "Este campo es requerido"))
				continue
			}
			if isNullable {
				continue
			}
		}

		if isRequired {
			msg := Required(l, value.Interface())
			if msg != "" {
				err.Append(key, msg)
				continue
			}
		}

		if isNullable {
			if isEmpty(value.Interface()) {
				continue
			}
		}

		// Aplicar reglas
		for _, rule := range r {
			param := ""
			rp := strings.Split(rule, ":")
			rule = rp[0]
			if len(rp) > 1 {
				param = rp[1]
			}

			switch rule {
			// case "required":
			// 	err.Append(key, Required(l, value.Interface()))
			case "min":
				limit, _ := strconv.Atoi(param)
				switch value.Kind() {
				case reflect.String:
					err.Append(key, MinString(l, value.String(), limit))
				case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
					err.Append(key, MinNumber(l, value.Int(), int64(limit)))
				case reflect.Float32, reflect.Float64:
					err.Append(key, MinNumber(l, value.Float(), float64(limit)))
				case reflect.Slice, reflect.Array:
					err.Append(key, MinSlice(l, value.Interface().([]any), limit))
				}
			case "max":
				limit, _ := strconv.Atoi(param)
				switch value.Kind() {
				case reflect.String:
					err.Append(key, MaxString(l, value.String(), limit))
				case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
					err.Append(key, MaxNumber(l, value.Int(), int64(limit)))
				case reflect.Float32, reflect.Float64:
					err.Append(key, MaxNumber(l, value.Float(), float64(limit)))
				case reflect.Slice, reflect.Array:
					err.Append(key, MaxSlice(l, value.Interface().([]any), limit))
				}
			case "email":
				err.Append(key, Email(l, value.String()))
			case "url":
				err.Append(key, URL(l, value.String()))
			case "uuid":
				err.Append(key, UUID(l, value.String()))
			case "json":
				err.Append(key, JSON(l, value.String()))
			case "slug":
				err.Append(key, Slug(l, value.String()))
			case "alpha":
				err.Append(key, Alpha(l, value.String()))
			case "alphaDash", "alpha_dash":
				err.Append(key, AlphaDash(l, value.String()))
			case "alphaSpaces", "alpha_espaces":
				err.Append(key, AlphaSpaces(l, value.String()))
			case "alphaDashSpaces", "alpha_dash_espaces":
				err.Append(key, AlphaDashSpaces(l, value.String()))
			case "alphaNum", "alpha_num":
				err.Append(key, AlphaNum(l, value.String()))
			case "alphaNumDash", "alpha_num_dash":
				err.Append(key, AlphaNumDash(l, value.String()))
			case "alphaNumSpaces", "alpha_num_espaces":
				err.Append(key, AlphaNumSpaces(l, value.String()))
			case "alphaNumDashSpaces", "alpha_num_dash_spaces":
				err.Append(key, AlphaNumDashSpaces(l, value.String()))
			case "alphaAccents", "alpha_accents":
				err.Append(key, AlphaAccents(l, value.String()))
			case "alphaDashAccents", "alpha_dash_accents":
				err.Append(key, AlphaDashAccents(l, value.String()))
			case "alphaSpacesAccents", "alpha_spaces_accents":
				err.Append(key, AlphaSpacesAccents(l, value.String()))
			case "alphaDashSpacesAccents", "alpha_dash_spaces_accents":
				err.Append(key, AlphaDashSpacesAccents(l, value.String()))
			case "alphaNumAccents", "alpha_num_accents":
				err.Append(key, AlphaNumAccents(l, value.String()))
			case "alphaNumDashAccents", "alpha_num_dash_accents":
				err.Append(key, AlphaNumDashAccents(l, value.String()))
			case "alphaNumSpacesAccents", "alpha_num_spaces_accents":
				err.Append(key, AlphaNumSpacesAccents(l, value.String()))
			case "alphaNumDashSpacesAccents", "alpha_num_dash_spaces_accents":
				err.Append(key, AlphaNumDashSpacesAccents(l, value.String()))
			case "username", "user_name":
				err.Append(key, Username(l, value.String()))
			case "startsWith", "starts_with":
				err.Append(key, StartsWith(l, value.String(), param))
			case "endsWith", "ends_with":
				err.Append(key, EndsWith(l, value.String(), param))
			case "contains":
				err.Append(key, Contains(l, value.String(), param))
			case "notContains", "not_contains":
				err.Append(key, NotContains(l, value.String(), param))
			case "in":
				values := strings.Split(param, ",")
				switch value.Kind() {
				case reflect.String:
					err.Append(key, In(l, value.String(), values...))
				case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
					ints := make([]int64, 0, len(values))
					for _, v := range values {
						if n, errConv := strconv.ParseInt(v, 10, 64); errConv == nil {
							ints = append(ints, n)
						}
					}
					err.Append(key, In(l, value.Int(), ints...))
				case reflect.Float32, reflect.Float64:
					floats := make([]float64, 0, len(values))
					for _, v := range values {
						if f, errConv := strconv.ParseFloat(v, 64); errConv == nil {
							floats = append(floats, f)
						}
					}
					err.Append(key, In(l, value.Float(), floats...))
				}
			case "nin", "not_in":
				values := strings.Split(param, ",")
				switch value.Kind() {
				case reflect.String:
					err.Append(key, Nin(l, value.String(), values...))
				case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
					ints := make([]int64, 0, len(values))
					for _, v := range values {
						if n, errConv := strconv.ParseInt(v, 10, 64); errConv == nil {
							ints = append(ints, n)
						}
					}
					err.Append(key, Nin(l, value.Int(), ints...))
				case reflect.Float32, reflect.Float64:
					floats := make([]float64, 0, len(values))
					for _, v := range values {
						if f, errConv := strconv.ParseFloat(v, 64); errConv == nil {
							floats = append(floats, f)
						}
					}
					err.Append(key, Nin(l, value.Float(), floats...))
				}
			case "unique":
				switch value.Kind() {
				case reflect.Slice:
					// convertir a []string si es posible esto hay que revisarlo luego
					// pendiente por revisar
					slice := value.Interface()
					if strSlice, ok := slice.([]string); ok {
						err.Append(key, Unique(l, strSlice))
					}
				}
			case "positive":
				switch value.Kind() {
				case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
					err.Append(key, Positive(l, value.Int()))
				case reflect.Float32, reflect.Float64:
					err.Append(key, Positive(l, value.Float()))
				}
			case "negative":
				switch value.Kind() {
				case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
					err.Append(key, Negative(l, value.Int()))
				case reflect.Float32, reflect.Float64:
					err.Append(key, Negative(l, value.Float()))
				}
			case "between":
				rangeVals := strings.Split(param, ",")
				if len(rangeVals) == 2 {
					switch value.Kind() {
					case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
						min, _ := strconv.ParseInt(rangeVals[0], 10, 64)
						max, _ := strconv.ParseInt(rangeVals[1], 10, 64)
						err.Append(key, Between(l, value.Int(), min, max))
					case reflect.Float32, reflect.Float64:
						min, _ := strconv.ParseFloat(rangeVals[0], 64)
						max, _ := strconv.ParseFloat(rangeVals[1], 64)
						err.Append(key, Between(l, value.Float(), min, max))
					}
				}
			case "before":
				t, e := time.Parse(time.RFC3339, param)
				if e != nil {
					err.Append(key, lang.TT(l, "Error inesperado: formato de fecha inválido %v", e.Error()))
					continue
				}
				err.Append(key, Before(l, value.Interface().(time.Time), t))
			case "after":
				t, e := time.Parse(time.RFC3339, param)
				if e != nil {
					err.Append(key, lang.TT(l, "Error inesperado: formato de fecha inválido %v", e.Error()))
					continue
				}
				err.Append(key, After(l, value.Interface().(time.Time), t))
			case "beforenow":
				err.Append(key, BeforeNow(l, value.Interface().(time.Time)))
			case "afternow":
				err.Append(key, AfterNow(l, value.Interface().(time.Time)))
			case "datebetween":
				rangeVals := strings.Split(param, ",")
				if len(rangeVals) == 2 {
					start, errStart := time.Parse(time.RFC3339, rangeVals[0])
					end, errEnd := time.Parse(time.RFC3339, rangeVals[1])
					if errStart == nil && errEnd == nil {
						err.Append(key, DateBetween(l, value.Interface().(time.Time), start, end))
						continue
					}
					if errStart != nil {
						err.Append(key, lang.TT(l, "Error inesperado: formato de fecha inválido %v", errStart.Error()))
					}
					if errEnd != nil {
						err.Append(key, lang.TT(l, "Error inesperado: formato de fecha inválido %v", errEnd.Error()))
					}
				}
			case "withoutall", "without", "withall", "with":
				otherKeys := strings.Split(param, ",")
				otherFields := make([]any, 0, len(otherKeys))
				for _, k := range otherKeys {
					foundOther := false
					for i := 0; i < typ.NumField(); i++ {
						otherField := typ.Field(i)
						tag := otherField.Tag.Get("json")
						tagName := strings.Split(tag, ",")[0]
						if tagName == k {
							otherFields = append(otherFields, val.Field(i).Interface())
							foundOther = true
							break
						}
					}
					if !foundOther {
						// Si no se encontró el campo, se agrega nil para mantener el orden
						otherFields = append(otherFields, nil)
					}
				}

				switch rule {
				case "withoutall":
					err.Append(key, WithoutAll(l, value.Interface(), otherKeys, otherFields...))
				case "without":
					err.Append(key, Without(l, value.Interface(), otherKeys, otherFields...))
				case "withall":
					err.Append(key, WithAll(l, value.Interface(), otherKeys, otherFields...))
				case "with":
					err.Append(key, With(l, value.Interface(), otherKeys, otherFields...))
				}
			}
		}
	}
	return err.Errors()
}

func MinNumber[T constraints.Integer | constraints.Float](l string, value T, limit T) string {
	if value < limit {
		return lang.TT(l, "Mínimo %v", limit)
	}
	return ""
}

func MaxNumber[T constraints.Integer | constraints.Float](l string, value T, limit T) string {
	if value > limit {
		return lang.TT(l, "Máximo %v", limit)
	}
	return ""
}

func MinString(l string, value string, limit int) string {
	if len(value) > limit {
		return lang.TT(l, "Minimo %v caracteres", limit)
	}
	return ""
}

func MaxString(l string, value string, limit int) string {
	if len(value) > limit {
		return lang.TT(l, "Máximo %v caracteres", limit)
	}
	return ""
}

func MinSlice(l string, value []any, limit int) string {
	if len(value) < limit {
		return lang.TT(l, "Minimo %v elementos", limit)
	}
	return ""
}

func MaxSlice(l string, value []any, limit int) string {
	if len(value) > limit {
		return lang.TT(l, "Máximo %v elementos", limit)
	}
	return ""
}

// isEmpty es una funcion auxiliar no hace parte de las validaciones
func isEmpty(value any) bool {
	if value == nil {
		return true
	}

	val := reflect.ValueOf(value)
	if !val.IsValid() {
		return true
	}

	switch val.Kind() {
	case reflect.String:
		return val.Len() == 0
	case reflect.Slice, reflect.Map, reflect.Array, reflect.Chan:
		return val.Len() == 0
	case reflect.Ptr, reflect.Interface:
		if val.IsNil() {
			return true
		}
		return isEmpty(val.Elem().Interface())
	case reflect.Func:
		return val.IsNil()
	default:
		return val.IsZero()
	}
}

func Required(l string, value any) string {
	if isEmpty(value) {
		return lang.TT(l, "Este campo es requerido")
	}
	return ""
}

// WithoutAll verifica si el campo debe estar presente cuando todos los otros campos están vacíos
func WithoutAll(l string, field any, otherFieldsNames []string, otherFields ...any) string {
	allEmpty := true
	for _, otherField := range otherFields {
		if !isEmpty(otherField) {
			allEmpty = false
			break
		}
	}

	if allEmpty && isEmpty(field) {
		if len(otherFieldsNames) > 1 {
			return lang.TT(l, "Este campo es requerido cuando [%v] están vacíos", otherFieldsNames)
		}
		return lang.TT(l, "Este campo es requerido cuando %v está vacío", otherFieldsNames)
	}

	return ""
}

// Without verifica si el campo debe estar presente cuando cualquiera de los otros campos está vacío
func Without(l string, field any, otherFieldsNames []string, otherFields ...any) string {
	anyEmpty := false
	for _, otherField := range otherFields {
		if isEmpty(otherField) {
			anyEmpty = true
			break
		}
	}

	if anyEmpty && isEmpty(field) {
		return lang.TT(l, "Este campo es requerido si algúno de estos [%v] está vacío", otherFieldsNames)
	}

	return ""
}

// WithAll verifica si el campo debe estar presente cuando todos los otros campos tienen valor
func WithAll(l string, field any, otherFieldsNames []string, otherFields ...any) string {
	allFilled := true
	for _, otherField := range otherFields {
		if isEmpty(otherField) {
			allFilled = false
			break
		}
	}

	if allFilled && isEmpty(field) {
		if len(otherFieldsNames) > 1 {
			return lang.TT(l, "Este campo es requerido si [%v] no estan vacios", otherFieldsNames)
		}
		return lang.TT(l, "Este campo es requerido si %v no esta vacio", otherFieldsNames)
	}

	return ""
}

// With verifica si el campo debe estar presente cuando cualquiera de los otros campos tiene valor
func With(l string, field any, otherFieldsNames []string, otherFields ...any) string {
	anyFilled := false
	for _, otherField := range otherFields {
		if !isEmpty(otherField) {
			anyFilled = true
			break
		}
	}

	if anyFilled && isEmpty(field) {
		if len(otherFieldsNames) > 1 {
			return lang.TT(l, "Este campo es requerido si alguno de estos [%v] no esta vacio", otherFieldsNames)
		}
		return lang.TT(l, "Este campo es requerido si %v no esta vacio", otherFieldsNames)
	}

	return ""
}

func Email(l string, value string) string {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		return lang.TT(l, "Correo electrónico inválido")
	}
	return ""
}

func URL(l string, value string) string {
	urlRegex := regexp.MustCompile(`^(https?://)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,6}(/[\w\-\./?%&=]*)?$`)
	if !urlRegex.MatchString(value) {
		return lang.TT(l, "URL inválida")
	}
	return ""
}

func UUID(l string, value string) string {
	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89ab][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
	if !uuidRegex.MatchString(value) {
		return lang.TT(l, "UUID inválido (debe ser v4)")
	}
	return ""
}

func JSON(l string, value string) string {
	var js map[string]interface{}
	if err := json.Unmarshal([]byte(value), &js); err != nil {
		return lang.TT(l, "El formato JSON es inválido")
	}
	return ""
}

func Slug(l string, value string) string {
	slugRegex := regexp.MustCompile(`^[a-z0-9]+(?:[-_][a-z0-9]+)*$`)
	if !slugRegex.MatchString(value) {
		return lang.TT(l, "Solo se permiten letras minúsculas, números, guiones y guiones bajos (sin empezar o terminar con ellos)")
	}
	return ""
}

func Alpha(l string, value string) string {
	alphaRegex := regexp.MustCompile(`^[a-zA-Z]+$`)
	if !alphaRegex.MatchString(value) {
		return lang.TT(l, "Solo se permiten letras")
	}
	return ""
}

func AlphaDash(l string, value string) string {
	alphaDashRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !alphaDashRegex.MatchString(value) {
		return lang.TT(l, "Solo se permiten letras, números, guiones y guiones bajos")
	}
	return ""
}

func AlphaSpaces(l string, value string) string {
	alphaSpacesRegex := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	if !alphaSpacesRegex.MatchString(value) {
		return lang.TT(l, "Solo se permiten letras y espacios")
	}
	return ""
}

func AlphaDashSpaces(l string, value string) string {
	alphaDashSpacesRegex := regexp.MustCompile(`^[a-zA-Z0-9_\s-]+$`)
	if !alphaDashSpacesRegex.MatchString(value) {
		return lang.TT(l, "Solo se permiten letras, números, guiones, guiones bajos y espacios")
	}
	return ""
}

func AlphaNum(l string, value string) string {
	alphaNumRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !alphaNumRegex.MatchString(value) {
		return lang.TT(l, "Solo se permiten letras y números")
	}
	return ""
}

func AlphaNumDash(l string, value string) string {
	alphaNumDashRegex := regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
	if !alphaNumDashRegex.MatchString(value) {
		return lang.TT(l, "Solo se permiten letras, números y guiones")
	}
	return ""
}

func AlphaNumSpaces(l string, value string) string {
	alphaNumSpacesRegex := regexp.MustCompile(`^[a-zA-Z0-9\s]+$`)
	if !alphaNumSpacesRegex.MatchString(value) {
		return lang.TT(l, "Solo se permiten letras, números y espacios")
	}
	return ""
}

func AlphaNumDashSpaces(l string, value string) string {
	alphaNumDashSpacesRegex := regexp.MustCompile(`^[a-zA-Z0-9_\s-]+$`)
	if !alphaNumDashSpacesRegex.MatchString(value) {
		return lang.TT(l, "Solo se permiten letras, números, guiones, guiones bajos y espacios")
	}
	return ""
}

func AlphaAccents(l string, value string) string {
	alphaAccentsRegex := regexp.MustCompile(`^[a-zA-ZáéíóúÁÉÍÓÚñÑüÜ]+$`)
	if !alphaAccentsRegex.MatchString(value) {
		return lang.TT(l, "Solo se permiten letras, incluyendo tildes y eñes")
	}
	return ""
}

func AlphaDashAccents(l string, value string) string {
	alphaDashAccentsRegex := regexp.MustCompile(`^[a-zA-ZáéíóúÁÉÍÓÚñÑüÜ_-]+$`)
	if !alphaDashAccentsRegex.MatchString(value) {
		return lang.TT(l, "Solo se permiten letras (con tildes), eñes, guiones y guiones bajos")
	}
	return ""
}

func AlphaSpacesAccents(l string, value string) string {
	alphaSpacesAccentsRegex := regexp.MustCompile(`^[a-zA-ZáéíóúÁÉÍÓÚñÑüÜ\s]+$`)
	if !alphaSpacesAccentsRegex.MatchString(value) {
		return lang.TT(l, "Solo se permiten letras, tildes, eñes y espacios")
	}
	return ""
}

func AlphaDashSpacesAccents(l string, value string) string {
	alphaDashSpacesAccentsRegex := regexp.MustCompile(`^[a-zA-ZáéíóúÁÉÍÓÚñÑüÜ_\s-]+$`)
	if !alphaDashSpacesAccentsRegex.MatchString(value) {
		return lang.TT(l, "Solo se permiten letras (con tildes), eñes, guiones, guiones bajos y espacios")
	}
	return ""
}

func AlphaNumAccents(l string, value string) string {
	alphaNumAccentsRegex := regexp.MustCompile(`^[a-zA-Z0-9áéíóúÁÉÍÓÚñÑüÜ]+$`)
	if !alphaNumAccentsRegex.MatchString(value) {
		return lang.TT(l, "Solo se permiten letras (con tildes), eñes y números")
	}
	return ""
}

func AlphaNumDashAccents(l string, value string) string {
	alphaNumDashAccentsRegex := regexp.MustCompile(`^[a-zA-Z0-9áéíóúÁÉÍÓÚñÑüÜ_-]+$`)
	if !alphaNumDashAccentsRegex.MatchString(value) {
		return lang.TT(l, "Solo se permiten letras (con tildes), eñes, números, guiones y guiones bajos")
	}
	return ""
}

func AlphaNumSpacesAccents(l string, value string) string {
	alphaNumSpacesAccentsRegex := regexp.MustCompile(`^[a-zA-Z0-9áéíóúÁÉÍÓÚñÑüÜ\s]+$`)
	if !alphaNumSpacesAccentsRegex.MatchString(value) {
		return lang.TT(l, "Solo se permiten letras (con tildes), eñes, números y espacios")
	}
	return ""
}

func AlphaNumDashSpacesAccents(l string, value string) string {
	alphaNumDashSpacesAccentsRegex := regexp.MustCompile(`^[a-zA-Z0-9áéíóúÁÉÍÓÚñÑüÜ_\s-]+$`)
	if !alphaNumDashSpacesAccentsRegex.MatchString(value) {
		return lang.TT(l, "Solo se permiten letras (con tildes), eñes, números, guiones, guiones bajos y espacios")
	}
	return ""
}

//	func Username(l string, value string) string {
//		usernameRegex := regexp.MustCompile(`^(?!.*[_.]{2})[a-zA-Z0-9](?:[a-zA-Z0-9._]*[a-zA-Z0-9])?$`)
//		if !usernameRegex.MatchString(value) {
//			return lang.TT(l, "Usuario inválido: solo letras, números, '.' o '_', sin empezar o terminar con ellos ni usarlos consecutivos")
//		}
//		return ""
//	}

// isAlphaNumeric es una funcion auxiliar no hace parte de las validaciones
func isAlphaNumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}
func Username(l string, value string) string {
	if len(value) == 0 || !isAlphaNumeric(rune(value[0])) || !isAlphaNumeric(rune(value[len(value)-1])) {
		return lang.TT(l, "Usuario inválido: solo letras, números, '.' o '_', sin empezar o terminar con ellos ni usarlos consecutivos")
	}

	for i := 0; i < len(value)-1; i++ {
		if (value[i] == '.' || value[i] == '_') && (value[i+1] == '.' || value[i+1] == '_') {
			return lang.TT(l, "Usuario inválido: solo letras, números, '.' o '_', sin empezar o terminar con ellos ni usarlos consecutivos")
		}
	}

	for _, char := range value {
		if !isAlphaNumeric(char) && char != '.' && char != '_' {
			return lang.TT(l, "Usuario inválido: solo letras, números, '.' o '_', sin empezar o terminar con ellos ni usarlos consecutivos")
		}
	}

	return ""
}

func StartsWith(l string, value string, prefix string) string {
	if !strings.HasPrefix(value, prefix) {
		return lang.TT(l, "Debe comenzar con: %v", prefix)
	}
	return ""
}

func EndsWith(l string, value string, suffix string) string {
	if !strings.HasSuffix(value, suffix) {
		return lang.TT(l, "Debe terminar con: %v", suffix)
	}
	return ""
}

func Contains(l string, value string, substr string) string {
	if !strings.Contains(value, substr) {
		return lang.TT(l, "Debe contener: %v", substr)
	}
	return ""
}

func NotContains(l string, value string, substr string) string {
	if strings.Contains(value, substr) {
		return lang.TT(l, "No debe contener: %v", substr)
	}
	return ""
}

func In[T comparable](l string, value T, allowed ...T) string {
	for _, v := range allowed {
		if value == v {
			return ""
		}
	}
	return lang.TT(l, "Valor no permitido, debe ser: %v", allowed)
}

func Nin[T comparable](l string, value T, denied ...T) string {
	for _, v := range denied {
		if value == v {
			return lang.TT(l, "Valor no permitido, no puede ser: %v", denied)
		}
	}
	return ""
}

func Unique[T comparable](l string, list []T) string {
	// Mapa para rastrear elementos vistos
	seen := make(map[T]bool)
	for _, item := range list {
		if seen[item] {
			return lang.TT(l, "el elemento [%v] esta duplicado", item)
		}
		seen[item] = true
	}
	return ""
}

func Positive[T constraints.Integer | constraints.Float](l string, value T) string {
	if value <= 0 {
		return lang.TT(l, "Debe ser mayor que 0")
	}
	return ""
}

func Negative[T constraints.Integer | constraints.Float](l string, value T) string {
	if value >= 0 {
		return lang.TT(l, "Debe ser menor que 0")
	}
	return ""
}

func Between[T constraints.Integer | constraints.Float](l string, value T, min T, max T) string {
	if value < min || value > max {
		return lang.TT(l, "Debe estar entre %v y %v", min, max)
	}
	return ""
}

func Before(l string, value time.Time, target time.Time) string {
	if value.After(target) || value.Equal(target) {
		return lang.TT(l, "Debe ser una fecha anterior a %v", target)
	}
	return ""
}

func After(l string, value time.Time, target time.Time) string {
	if value.Before(target) || value.Equal(target) {
		return lang.TT(l, "Debe ser una fecha posterior a %v", target)
	}
	return ""
}

func BeforeNow(l string, value time.Time) string {
	if value.After(time.Now()) {
		return lang.TT(l, "Debe ser una fecha anterior al momento actual")
	}
	return ""
}

func AfterNow(l string, value time.Time) string {
	if value.Before(time.Now()) {
		return lang.TT(l, "Debe ser una fecha posterior al momento actual")
	}
	return ""
}

func DateBetween(l string, value time.Time, start time.Time, end time.Time) string {
	if value.Before(start) || value.After(end) {
		return lang.TT(l, "Debe estar entre %v y %v", start, end)
	}
	return ""
}
