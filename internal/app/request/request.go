package request

import (
	"encoding/json"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/pkg/errors"
	"github.com/donbarrigon/nuevo-proyecto/pkg/lang"
	"golang.org/x/exp/constraints"
)

type Validator struct {
	Lang string
	Err  errors.Err
}

func NewValidator(l string) *Validator {
	return &Validator{
		Lang: l,
		Err: errors.Err{
			ErrMap: make(map[string][]string),
		},
	}
}

func (v *Validator) Rules(field string, value any, rules ...string) {
	f := 10.5
	m := MinNumber("es", f, 10)
	log.Println(m)
}

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
		return lang.TT(l, "Campo obligatorio")
	}
	return ""
}

// WithoutAll verifica si el campo debe estar presente cuando todos los otros campos están vacíos
func WithoutAll(l string, field any, otherFields ...any) string {
	allEmpty := true
	for _, otherField := range otherFields {
		if !isEmpty(otherField) {
			allEmpty = false
			break
		}
	}

	if allEmpty && isEmpty(field) {
		return lang.TT(l, "Obligatorio cuando los demás están vacíos")
	}

	return ""
}

// Without verifica si el campo debe estar presente cuando cualquiera de los otros campos está vacío
func Without(l string, field any, otherFields ...any) string {
	anyEmpty := false
	for _, otherField := range otherFields {
		if isEmpty(otherField) {
			anyEmpty = true
			break
		}
	}

	if anyEmpty && isEmpty(field) {
		return lang.TT(l, "Obligatorio si algún otro está vacío")
	}

	return ""
}

// WithAll verifica si el campo debe estar presente cuando todos los otros campos tienen valor
func WithAll(l string, field any, otherFields ...any) string {
	allFilled := true
	for _, otherField := range otherFields {
		if isEmpty(otherField) {
			allFilled = false
			break
		}
	}

	if allFilled && isEmpty(field) {
		return lang.TT(l, "Obligatorio si todos tienen valor")
	}

	return ""
}

// With verifica si el campo debe estar presente cuando cualquiera de los otros campos tiene valor
func With(l string, field any, otherFields ...any) string {
	anyFilled := false
	for _, otherField := range otherFields {
		if !isEmpty(otherField) {
			anyFilled = true
			break
		}
	}

	if anyFilled && isEmpty(field) {
		return lang.TT(l, "Obligatorio si alguno tiene valor")
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

func AlphaNumSpaces(l string, value string) string {
	alphaNumSpacesRegex := regexp.MustCompile(`^[a-zA-Z0-9\s]+$`)
	if !alphaNumSpacesRegex.MatchString(value) {
		return lang.TT(l, "Solo se permiten letras, números y espacios")
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

func AlphaNumAccents(l string, value string) string {
	alphaNumAccentsRegex := regexp.MustCompile(`^[a-zA-Z0-9áéíóúÁÉÍÓÚñÑüÜ]+$`)
	if !alphaNumAccentsRegex.MatchString(value) {
		return lang.TT(l, "Solo se permiten letras (con tildes), eñes y números")
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

func AlphaNumSpacesAccents(l string, value string) string {
	alphaNumSpacesAccentsRegex := regexp.MustCompile(`^[a-zA-Z0-9áéíóúÁÉÍÓÚñÑüÜ\s]+$`)
	if !alphaNumSpacesAccentsRegex.MatchString(value) {
		return lang.TT(l, "Solo se permiten letras (con tildes), eñes, números y espacios")
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

func Slug(l string, value string) string {
	slugRegex := regexp.MustCompile(`^[a-z0-9]+(?:[-_][a-z0-9]+)*$`)
	if !slugRegex.MatchString(value) {
		return lang.TT(l, "Solo se permiten letras minúsculas, números, guiones y guiones bajos (sin empezar o terminar con ellos)")
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

////////////////////////

// func (v *Validator) MinNumber [Number constraints.Integer | constraints.Float](field string, value Number, min Number) {
// 	if value < min {
// 		v.Err.Append(field, lang.TT(v.Lang, "Mínimo %v", min))
// 	}
// }

// func (v *Validator[Number, Comparable]) MaxNumber(field string, value Number, max Number) {
// 	if value > max {
// 		v.Err.Append(field, lang.TT(v.Lang, "Máximo %v", max))
// 	}
// }

// func (v *Validator[Number, Comparable]) MinString(field string, value string, min int) {
// 	if len(value) < min {
// 		v.Err.Append(field, lang.TT(v.Lang, "Mínimo %v caracteres", min))
// 	}
// }

// func (v *Validator[Number, Comparable]) MaxString(field string, value string, max int) {
// 	if len(value) > max {
// 		v.Err.Append(field, lang.TT(v.Lang, "Máximo %v caracteres", max))
// 	}
// }

// func (v *Validator[Number, Comparable]) MinSlice(field string, value []any, min int) {
// 	if len(value) < min {
// 		v.Err.Append(field, lang.TT(v.Lang, "Mínimo %v elementos", min))
// 	}
// }

// func (v *Validator[Number, Comparable]) MaxSlice(field string, value []any, max int) {
// 	if len(value) > max {
// 		v.Err.Append(field, lang.TT(v.Lang, "Máximo %v elementos", max))
// 	}
// }

// func (v *Validator[Number, Comparable]) Required(field string, value any) {
// 	if isEmpty(value) {
// 		v.Err.Append(field, lang.TT(v.Lang, "Campo obligatorio"))
// 	}
// }

// func (v *Validator[Number, Comparable]) WithoutAll(field string, value any, otherFields ...any) {
// 	allEmpty := true
// 	for _, otherField := range otherFields {
// 		if !isEmpty(otherField) {
// 			allEmpty = false
// 			break
// 		}
// 	}

// 	if allEmpty && isEmpty(value) {
// 		v.Err.Append(field, lang.TT(v.Lang, "Obligatorio cuando los demás están vacíos"))
// 	}
// }

// func (v *Validator[Number, Comparable]) Without(field string, value any, otherFields ...any) {
// 	anyEmpty := false
// 	for _, otherField := range otherFields {
// 		if isEmpty(otherField) {
// 			anyEmpty = true
// 			break
// 		}
// 	}

// 	if anyEmpty && isEmpty(value) {
// 		v.Err.Append(field, lang.TT(v.Lang, "Obligatorio si algún otro está vacío"))
// 	}
// }

// func (v *Validator[Number, Comparable]) WithAll(field string, value any, otherFields ...any) {
// 	allFilled := true
// 	for _, otherField := range otherFields {
// 		if isEmpty(otherField) {
// 			allFilled = false
// 			break
// 		}
// 	}

// 	if allFilled && isEmpty(value) {
// 		v.Err.Append(field, lang.TT(v.Lang, "Obligatorio si todos tienen valor"))
// 	}
// }

// func (v *Validator[Number, Comparable]) With(field string, value any, otherFields ...any) {
// 	anyFilled := false
// 	for _, otherField := range otherFields {
// 		if !isEmpty(otherField) {
// 			anyFilled = true
// 			break
// 		}
// 	}

// 	if anyFilled && isEmpty(value) {
// 		v.Err.Append(field, lang.TT(v.Lang, "Obligatorio si algun otro tiene valor"))
// 	}
// }

// func (v *Validator[Number, Comparable]) Email(field string, value string) {
// 	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
// 	if !emailRegex.MatchString(value) {
// 		v.Err.Append(field, lang.TT(v.Lang, "Correo electrónico inválido"))
// 	}
// }

// func (v *Validator[Number, Comparable]) URL(field string, value string) {
// 	urlRegex := regexp.MustCompile(`^(https?://)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,6}(/[\w\-\./?%&=]*)?$`)
// 	if !urlRegex.MatchString(value) {
// 		v.Err.Append(field, lang.TT(v.Lang, "URL inválida"))
// 	}
// }

// func (v *Validator[Number, Comparable]) UUID(field string, value string) {
// 	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89ab][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
// 	if !uuidRegex.MatchString(value) {
// 		v.Err.Append(field, lang.TT(v.Lang, "UUID inválido (debe ser v4)"))
// 	}
// }

// func (v *Validator[Number, Comparable]) JSON(field string, value string) {
// 	var js map[string]interface{}
// 	if err := json.Unmarshal([]byte(value), &js); err != nil {
// 		v.Err.Append(field, lang.TT(v.Lang, "El formato JSON es inválido"))
// 	}
// }

// func (v *Validator[Number, Comparable]) Alpha(field string, value string) {
// 	if !regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(value) {
// 		v.Err.Append(field, lang.TT(v.Lang, "Solo se permiten letras"))
// 	}
// }

// func (v *Validator[Number, Comparable]) AlphaNum(field string, value string) {
// 	if !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(value) {
// 		v.Err.Append(field, lang.TT(v.Lang, "Solo se permiten letras y números"))
// 	}
// }

// func (v *Validator[Number, Comparable]) AlphaDash(field string, value string) {
// 	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(value) {
// 		v.Err.Append(field, lang.TT(v.Lang, "Solo se permiten letras, números, guiones y guiones bajos"))
// 	}
// }

// func (v *Validator[Number, Comparable]) AlphaSpaces(field string, value string) {
// 	if !regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(value) {
// 		v.Err.Append(field, lang.TT(v.Lang, "Solo se permiten letras y espacios"))
// 	}
// }

// func (v *Validator[Number, Comparable]) AlphaNumSpaces(field string, value string) {
// 	if !regexp.MustCompile(`^[a-zA-Z0-9\s]+$`).MatchString(value) {
// 		v.Err.Append(field, lang.TT(v.Lang, "Solo se permiten letras, números y espacios"))
// 	}
// }

// func (v *Validator[Number, Comparable]) AlphaAccents(field string, value string) {
// 	if !regexp.MustCompile(`^[a-zA-ZáéíóúÁÉÍÓÚñÑüÜ]+$`).MatchString(value) {
// 		v.Err.Append(field, lang.TT(v.Lang, "Solo se permiten letras, incluyendo tildes y eñes"))
// 	}
// }

// func (v *Validator[Number, Comparable]) AlphaNumAccents(field string, value string) {
// 	if !regexp.MustCompile(`^[a-zA-Z0-9áéíóúÁÉÍÓÚñÑüÜ]+$`).MatchString(value) {
// 		v.Err.Append(field, lang.TT(v.Lang, "Solo se permiten letras (con tildes), eñes y números"))
// 	}
// }

// func (v *Validator[Number, Comparable]) AlphaSpacesAccents(field string, value string) {
// 	if !regexp.MustCompile(`^[a-zA-ZáéíóúÁÉÍÓÚñÑüÜ\s]+$`).MatchString(value) {
// 		v.Err.Append(field, lang.TT(v.Lang, "Solo se permiten letras, tildes, eñes y espacios"))
// 	}
// }

// func (v *Validator[Number, Comparable]) AlphaNumSpacesAccents(field string, value string) {
// 	if !regexp.MustCompile(`^[a-zA-Z0-9áéíóúÁÉÍÓÚñÑüÜ\s]+$`).MatchString(value) {
// 		v.Err.Append(field, lang.TT(v.Lang, "Solo se permiten letras (con tildes), eñes, números y espacios"))
// 	}
// }

// func (v *Validator[Number, Comparable]) AlphaDashAccents(field string, value string) {
// 	if !regexp.MustCompile(`^[a-zA-ZáéíóúÁÉÍÓÚñÑüÜ_-]+$`).MatchString(value) {
// 		v.Err.Append(field, lang.TT(v.Lang, "Solo se permiten letras (con tildes), eñes, guiones y guiones bajos"))
// 	}
// }

// func (v *Validator[Number, Comparable]) Slug(field string, value string) {
// 	if !regexp.MustCompile(`^[a-z0-9]+(?:[-_][a-z0-9]+)*$`).MatchString(value) {
// 		v.Err.Append(field, lang.TT(v.Lang, "Solo se permiten letras minúsculas, números, guiones y guiones bajos (sin empezar o terminar con ellos)"))
// 	}
// }

// func isAlphaNumeric(r rune) bool {
// 	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
// }

// func (v *Validator[Number, Comparable]) Username(field string, value string) {
// 	if len(value) == 0 || !isAlphaNumeric(rune(value[0])) || !isAlphaNumeric(rune(value[len(value)-1])) {
// 		v.Err.Append(field, lang.TT(v.Lang, "Usuario inválido: solo letras, números, '.' o '_', sin empezar o terminar con ellos ni usarlos consecutivos"))
// 		return
// 	}

// 	for i := 0; i < len(value)-1; i++ {
// 		if (value[i] == '.' || value[i] == '_') && (value[i+1] == '.' || value[i+1] == '_') {
// 			v.Err.Append(field, lang.TT(v.Lang, "Usuario inválido: solo letras, números, '.' o '_', sin empezar o terminar con ellos ni usarlos consecutivos"))
// 			return
// 		}
// 	}

// 	for _, char := range value {
// 		if !isAlphaNumeric(char) && char != '.' && char != '_' {
// 			v.Err.Append(field, lang.TT(v.Lang, "Usuario inválido: solo letras, números, '.' o '_', sin empezar o terminar con ellos ni usarlos consecutivos"))
// 			return
// 		}
// 	}
// }
