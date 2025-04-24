package request

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/pkg/lang"
	"golang.org/x/exp/constraints"
)

func MinNumber[T constraints.Integer | constraints.Float](l string, value T, min T) string {
	if value < min {
		return lang.TT(l, "Mínimo %v", min)
	}
	return ""
}

func MaxNumber[T constraints.Integer | constraints.Float](l string, value T, max T) string {
	if value > max {
		return lang.TT(l, "Máximo %v", max)
	}
	return ""
}

func MinString(l string, value string, min int) string {
	if len(value) > min {
		return lang.TT(l, "Minimo %v caracteres", min)
	}
	return ""
}

func MaxString(l string, value string, max int) string {
	if len(value) > max {
		return lang.TT(l, "Máximo %v caracteres", max)
	}
	return ""
}

func MinSlice(l string, value []any, min int) string {
	if len(value) < min {
		return lang.TT(l, "Minimo %v elementos", min)
	}
	return ""
}

func MaxSlice(l string, value []any, max int) string {
	if len(value) > max {
		return lang.TT(l, "Máximo %v elementos", max)
	}
	return ""
}

func Required(l string, value any) string {
	if str, ok := value.(string); ok && str == "" {
		return lang.TT(l, "Este campo es obligatorio")
	}

	val := reflect.ValueOf(value)
	if (val.Kind() == reflect.Slice || val.Kind() == reflect.Map) && val.Len() == 0 {
		return lang.TT(l, "Este campo es obligatorio")
	}

	if val.IsZero() {
		return lang.TT(l, "Este campo es obligatorio")
	}

	return ""
}

// WithoutAll verifica si el campo debe estar presente cuando todos los otros campos están vacíos
func WithoutAll(l string, field any, otherFields ...any) string {
	allEmpty := true
	for _, otherField := range otherFields {
		if !reflect.ValueOf(otherField).IsZero() {
			allEmpty = false
			break
		}
	}

	if allEmpty {
		if str, ok := field.(string); ok && str == "" {
			return lang.TT(l, "Este campo es obligatorio cuando todos los demás campos están vacíos")
		}

		val := reflect.ValueOf(field)
		if (val.Kind() == reflect.Slice || val.Kind() == reflect.Map) && val.Len() == 0 {
			return lang.TT(l, "Este campo es obligatorio cuando todos los demás campos están vacíos")
		}

		if val.IsZero() {
			return lang.TT(l, "Este campo es obligatorio cuando todos los demás campos están vacíos")
		}
	}

	return ""
}

// Without verifica si el campo debe estar presente cuando cualquiera de los otros campos está vacío
func Without(l string, field any, otherFields ...any) string {
	// Verificamos si alguno de los otros campos está vacío
	anyEmpty := false
	for _, otherField := range otherFields {
		if reflect.ValueOf(otherField).IsZero() {
			anyEmpty = true
			break
		}
	}

	// Si alguno de los otros campos está vacío, el campo principal debe tener valor
	if anyEmpty && reflect.ValueOf(field).IsZero() {
		return lang.TT(l, "Este campo es obligatorio cuando algún otro campo está vacío")
	}

	return ""
}

// WithAll verifica si el campo debe estar presente cuando todos los otros campos tienen valor
func WithAll(l string, field any, otherFields ...any) string {
	// Verificamos si todos los otros campos tienen valor
	allFilled := true
	for _, otherField := range otherFields {
		if reflect.ValueOf(otherField).IsZero() {
			allFilled = false
			break
		}
	}

	// Si todos los otros campos tienen valor, el campo principal debe tener valor
	if allFilled && reflect.ValueOf(field).IsZero() {
		return lang.TT(l, "Este campo es obligatorio cuando todos los demás campos tienen valor")
	}

	return ""
}

// With verifica si el campo debe estar presente cuando cualquiera de los otros campos tiene valor
func With(l string, field any, otherFields ...any) string {
	// Verificamos si alguno de los otros campos tiene valor
	anyFilled := false
	for _, otherField := range otherFields {
		if !reflect.ValueOf(otherField).IsZero() {
			anyFilled = true
			break
		}
	}

	// Si alguno de los otros campos tiene valor, el campo principal debe tener valor
	if anyFilled && reflect.ValueOf(field).IsZero() {
		return lang.TT(l, "Este campo es obligatorio cuando algún otro campo tiene valor")
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

func Alpha(l string, value string) string {
	alphaRegex := regexp.MustCompile(`^[a-zA-Z]+$`)
	if !alphaRegex.MatchString(value) {
		return lang.TT(l, "Solo se permiten letras")
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

func In[T comparable](l string, value T, allowed ...T) string {
	for _, v := range allowed {
		if value == v {
			return ""
		}
	}
	return lang.TT(l, "Valor no permitido, debe ser uno de: %v", allowed)
}

func Nin[T comparable](l string, value T, denied ...T) string {
	for _, v := range denied {
		if value == v {
			return lang.TT(l, "Valor no permitido, no puede ser uno de: %v", denied)
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
