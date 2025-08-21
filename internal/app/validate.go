package app

import (
	"context"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/exp/constraints"
)

func Validate(ctx *HttpContext, req any) Error {
	rulesMap := make(map[string][]string)

	t := reflect.TypeOf(req)

	// Si es un puntero, desreferencia
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		ruleTag := field.Tag.Get("rules")

		if ruleTag != "" {
			rules := strings.Split(ruleTag, "|")
			jsonTag := field.Tag.Get("json")
			jsonName := strings.Split(jsonTag, ",")[0] // Por si el tag es "name,omEntrypty"
			if jsonName == "" {
				jsonName = field.Name // fallback al nombre del campo si no hay tag
			}
			rulesMap[jsonName] = rules
		}
	}

	return ValidateRules(ctx, req, rulesMap)
}

func ValidateRules(ctx *HttpContext, req any, rules map[string][]string) Error {
	err := Errors.NewEmpty()

	val := reflect.ValueOf(req)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return Errors.Unknownf("The request is not a valid struct.")
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
				err.Append(&FieldError{
					FieldName:    key,
					Message:      "Este campo es requerido",
					Placeholders: List{{"attribute", key}},
				})
				continue
			}
			if isNullable {
				continue
			}
		}

		if isRequired {
			if e := ValidateRequired(key, value.Interface()); e != nil {
				err.Append(e)
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
			// 	err.Append(Required(key, value.Interface()))
			case "min":
				limit, _ := strconv.Atoi(param)
				switch value.Kind() {
				case reflect.String:
					err.Append(ValidateMinString(key, value.String(), limit))
				case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
					err.Append(ValidateMinNumber(key, value.Int(), int64(limit)))
				case reflect.Float32, reflect.Float64:
					err.Append(ValidateMinNumber(key, value.Float(), float64(limit)))
				case reflect.Slice, reflect.Array:
					err.Append(ValidateMinSlice(key, value.Interface().([]any), limit))
				}
			case "max":
				limit, _ := strconv.Atoi(param)
				switch value.Kind() {
				case reflect.String:
					err.Append(ValidateMaxString(key, value.String(), limit))
				case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
					err.Append(ValidateMaxNumber(key, value.Int(), int64(limit)))
				case reflect.Float32, reflect.Float64:
					err.Append(ValidateMaxNumber(key, value.Float(), float64(limit)))
				case reflect.Slice, reflect.Array:
					err.Append(ValidateMaxSlice(key, value.Interface().([]any), limit))
				}
			case "required_if":
				otherValue := getOtherFieldValueFromParam(val, param)
				err.Append(ValidateRequiredIf(key, value.Interface(), otherValue, param))
			case "required_unless":
				otherValue := getOtherFieldValueFromParam(val, param)
				err.Append(ValidateRequiredUnless(key, value.Interface(), otherValue, param))
			case "required_without", "required_with", "required_without_all", "required_with_all":
				otherKeys := strings.Split(param, ",")
				otherFields := make([]any, 0, len(otherKeys))
				for _, k := range otherKeys {
					otherFields = append(otherFields, getOtherFieldValueFromParam(val, k))
				}

				switch rule {
				case "required_without_all":
					err.Append(ValidateWithoutAll(key, value.Interface(), otherKeys, otherFields...))
				case "required_without":
					err.Append(ValidateWithout(key, value.Interface(), otherKeys, otherFields...))
				case "required_with_all":
					err.Append(ValidateWithAll(key, value.Interface(), otherKeys, otherFields...))
				case "required_with":
					err.Append(ValidateWith(key, value.Interface(), otherKeys, otherFields...))
				}
			case "same":
				otherValue := getOtherFieldValueFromParam(val, param)
				err.Append(ValidateSame(key, value.Interface(), param, otherValue))

			case "different":
				otherValue := getOtherFieldValueFromParam(val, param)
				err.Append(ValidateDifferent(key, value.Interface(), param, otherValue))

			case "confirmed":
				confirmationattribute := key + "_confirmation"
				confirmationValue := getOtherFieldValueFromParam(val, confirmationattribute)
				err.Append(ValidateConfirmed(key, value.Interface(), confirmationValue))
			case "accepted":
				err.Append(ValidateAccepted(key, value.Interface()))
			case "declined":
				err.Append(ValidateDeclined(key, value.Interface()))
			case "digits":
				limit, _ := strconv.Atoi(param)
				err.Append(ValidateDigits(key, value.Interface(), limit))
			case "digits_between":
				rangeParts := strings.Split(param, ",")
				if len(rangeParts) == 2 {
					min, _ := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
					max, _ := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
					err.Append(ValidateDigitsBetween(key, value.Interface(), min, max))
				} else {
					err.Append(&FieldError{
						FieldName: key,
						Message:   "The {attribute} has invalid parameters for digits_between.",
						Placeholders: List{
							{"attribute", key},
						},
					})
				}
			case "email":
				err.Append(ValidateEmail(key, value.String()))
			case "url":
				err.Append(ValidateURL(key, value.String()))
			case "uuid":
				err.Append(ValidateUUID(key, value.String()))
			case "ulid":
				err.Append(ValidateULID(key, value.String()))
			case "ip":
				err.Append(ValidateIP(key, value.String()))
			case "ipv4":
				err.Append(ValidateIPv4(key, value.String()))
			case "ipv6":
				err.Append(ValidateIPv6(key, value.String()))
			case "mac", "mac_address":
				err.Append(ValidateMACAddress(key, value.String()))
			case "ascii":
				err.Append(ValidateASCII(key, value.String()))
			case "lowercase":
				err.Append(ValidateLowercase(key, value.String()))
			case "uppercase":
				err.Append(ValidateUppercase(key, value.String()))
			case "hex":
				err.Append(ValidateHex(key, value.String()))
			case "hex_color":
				err.Append(ValidateHexColor(key, value.String()))
			case "json":
				err.Append(ValidateJSON(key, value.String()))
			case "slug":
				err.Append(ValidateSlug(key, value.String()))
			case "regex":
				err.Append(ValidateRegex(key, value.String(), param))
			case "not_regex":
				err.Append(ValidateNotRegex(key, value.String(), param))
			case "alpha":
				err.Append(ValidateAlpha(key, value.String()))
			case "alpha_dash":
				err.Append(ValidateAlphaDash(key, value.String()))
			case "alpha_espaces":
				err.Append(ValidateAlphaSpaces(key, value.String()))
			case "alpha_dash_espaces":
				err.Append(ValidateAlphaDashSpaces(key, value.String()))
			case "alpha_num":
				err.Append(ValidateAlphaNum(key, value.String()))
			case "alpha_num_dash":
				err.Append(ValidateAlphaNumDash(key, value.String()))
			case "alpha_num_espaces":
				err.Append(ValidateAlphaNumSpaces(key, value.String()))
			case "alpha_num_dash_spaces":
				err.Append(ValidateAlphaNumDashSpaces(key, value.String()))
			case "alpha_accents":
				err.Append(ValidateAlphaAccents(key, value.String()))
			case "alpha_dash_accents":
				err.Append(ValidateAlphaDashAccents(key, value.String()))
			case "alpha_spaces_accents":
				err.Append(ValidateAlphaSpacesAccents(key, value.String()))
			case "alpha_dash_spaces_accents":
				err.Append(ValidateAlphaDashSpacesAccents(key, value.String()))
			case "alpha_num_accents":
				err.Append(ValidateAlphaNumAccents(key, value.String()))
			case "alpha_num_dash_accents":
				err.Append(ValidateAlphaNumDashAccents(key, value.String()))
			case "alpha_num_spaces_accents":
				err.Append(ValidateAlphaNumSpacesAccents(key, value.String()))
			case "alpha_num_dash_spaces_accents":
				err.Append(ValidateAlphaNumDashSpacesAccents(key, value.String()))
			case "username", "user_name":
				err.Append(ValidateUsername(key, value.String()))
			case "starts_with":
				err.Append(ValidateStartsWith(key, value.String(), param))
			case "ends_with":
				err.Append(ValidateEndsWith(key, value.String(), param))
			case "contains":
				err.Append(ValidateContains(key, value.String(), param))
			case "not_contains":
				err.Append(ValidateNotContains(key, value.String(), param))
			case "in":
				values := strings.Split(param, ",")
				switch value.Kind() {
				case reflect.String:
					err.Append(ValidateIn(key, value.String(), values...))
				case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
					ints := make([]int64, 0, len(values))
					for _, v := range values {
						if n, errConv := strconv.ParseInt(v, 10, 64); errConv == nil {
							ints = append(ints, n)
						}
					}
					err.Append(ValidateIn(key, value.Int(), ints...))
				case reflect.Float32, reflect.Float64:
					floats := make([]float64, 0, len(values))
					for _, v := range values {
						if f, errConv := strconv.ParseFloat(v, 64); errConv == nil {
							floats = append(floats, f)
						}
					}
					err.Append(ValidateIn(key, value.Float(), floats...))
				}
			case "nin", "not_in":
				values := strings.Split(param, ",")
				switch value.Kind() {
				case reflect.String:
					err.Append(ValidateNin(key, value.String(), values...))
				case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
					ints := make([]int64, 0, len(values))
					for _, v := range values {
						if n, errConv := strconv.ParseInt(v, 10, 64); errConv == nil {
							ints = append(ints, n)
						}
					}
					err.Append(ValidateNin(key, value.Int(), ints...))
				case reflect.Float32, reflect.Float64:
					floats := make([]float64, 0, len(values))
					for _, v := range values {
						if f, errConv := strconv.ParseFloat(v, 64); errConv == nil {
							floats = append(floats, f)
						}
					}
					err.Append(ValidateNin(key, value.Float(), floats...))
				}
			case "exists":
				params := strings.Split(param, ",")
				if len(params) != 2 {
					err.Appendf(key, "the exists rule must have two parameters")
				} else {
					var v any
					if strings.HasSuffix(params[1], "_id") {
						var er error
						v, er = bson.ObjectIDFromHex(value.String())
						if er != nil {
							err.Appendf(key, "the value is not a valid id: :error", Entry{"error", er})
						}
					} else {
						v = value
					}
					err.Append(ValidateExists(key, params[0], params[1], v))
				}
			case "unique":
				params := strings.Split(param, ",")
				if len(params) != 2 {
					err.Appendf(key, "the unique rule must have two parameters")
				} else {
					err.Append(ValidateUnique(key, params[0], params[1], value, ctx.Params["id"]))
				}
			case "unique_in":
				switch value.Kind() {
				case reflect.Slice:
					// convertir a []string si es posible esto hay que revisarlo luego
					// pendiente por revisar
					slice := value.Interface()
					if strSlice, ok := slice.([]string); ok {
						err.Append(ValidateUniqueIn(key, strSlice))
					}
				}
			case "positive":
				switch value.Kind() {
				case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
					err.Append(ValidatePositive(key, value.Int()))
				case reflect.Float32, reflect.Float64:
					err.Append(ValidatePositive(key, value.Float()))
				}
			case "negative":
				switch value.Kind() {
				case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
					err.Append(ValidateNegative(key, value.Int()))
				case reflect.Float32, reflect.Float64:
					err.Append(ValidateNegative(key, value.Float()))
				}
			case "between":
				rangeVals := strings.Split(param, ",")
				if len(rangeVals) == 2 {
					switch value.Kind() {
					case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
						min, _ := strconv.ParseInt(rangeVals[0], 10, 64)
						max, _ := strconv.ParseInt(rangeVals[1], 10, 64)
						err.Append(ValidateBetween(key, value.Int(), min, max))
					case reflect.Float32, reflect.Float64:
						min, _ := strconv.ParseFloat(rangeVals[0], 64)
						max, _ := strconv.ParseFloat(rangeVals[1], 64)
						err.Append(ValidateBetween(key, value.Float(), min, max))
					}
				}
			case "before":
				t, e := time.Parse(time.RFC3339, param)
				if e != nil {
					err.Append(&FieldError{
						FieldName: key,
						Message:   "The {attribute} has an invalid date format: {error}.",
						Placeholders: List{
							{"attribute", key},
							{"error", e.Error()},
						},
					})
					continue
				}
				err.Append(ValidateBefore(key, value.Interface().(time.Time), t))
			case "after":
				t, e := time.Parse(time.RFC3339, param)
				if e != nil {
					err.Append(&FieldError{
						FieldName: key,
						Message:   "The {attribute} has an invalid date format: {error}.",
						Placeholders: List{
							{"attribute", key},
							{"error", e.Error()},
						},
					})
					continue
				}
				err.Append(ValidateAfter(key, value.Interface().(time.Time), t))
			case "before_now":
				err.Append(ValidateBeforeNow(key, value.Interface().(time.Time)))
			case "after_now":
				err.Append(ValidateAfterNow(key, value.Interface().(time.Time)))
			case "date_between":
				rangeVals := strings.Split(param, ",")
				if len(rangeVals) == 2 {
					start, errStart := time.Parse(time.RFC3339, rangeVals[0])
					end, errEnd := time.Parse(time.RFC3339, rangeVals[1])
					if errStart == nil && errEnd == nil {
						err.Append(ValidateDateBetween(key, value.Interface().(time.Time), start, end))
						continue
					}
					if errStart != nil {
						err.Append(&FieldError{
							FieldName: "start",
							Message:   "The {attribute} has an invalid date format: {error}.",
							Placeholders: List{
								{"attribute", "start"},
								{"error", errStart.Error()},
							},
						})
					}
					if errEnd != nil {
						err.Append(&FieldError{
							FieldName: "end",
							Message:   "The {attribute} has an invalid date format: {error}.",
							Placeholders: List{
								{"attribute", "end"},
								{"error", errEnd.Error()},
							},
						})
					}
				}
			}
		}
	}
	return err.Errors()
}

// getOtherFieldValueFromParam es una funcion auxiliar no hace parte de las validaciones
func getOtherFieldValueFromParam(val reflect.Value, param string) any {
	fieldKey := param
	for _, op := range []string{">=", "<=", "!=", "==", ">", "<"} {
		if strings.Contains(param, op) {
			parts := strings.SplitN(param, op, 2)
			fieldKey = strings.TrimSpace(parts[0])
			break
		}
	}
	if parts := strings.Split(param, ","); len(parts) > 1 {
		fieldKey = strings.TrimSpace(parts[0])
	}
	fieldKey = strings.TrimSpace(fieldKey)

	t := val.Type()
	for i := 0; i < t.NumField(); i++ {
		jsonTag := t.Field(i).Tag.Get("json")
		tagName := strings.Split(jsonTag, ",")[0]
		if tagName == fieldKey {
			return val.Field(i).Interface()
		}
	}
	return nil
}

func ValidateMinNumber[T constraints.Integer | constraints.Float](attribute string, value T, limit T) *FieldError {
	if value < limit {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} must be at least {limit}.",
			Placeholders: List{
				{"attribute", attribute},
				{"limit", limit},
			},
		}
	}
	return nil
}

func ValidateMaxNumber[T constraints.Integer | constraints.Float](attribute string, value T, limit T) *FieldError {
	if value > limit {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} may not be greater than {limit}.",
			Placeholders: List{
				{"attribute", attribute},
				{"limit", limit},
			},
		}
	}
	return nil
}

func ValidateMinString(attribute string, value string, limit int) *FieldError {
	if len(value) < limit {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} must be at least {limit} characters.",
			Placeholders: List{
				{"attribute", attribute},
				{"limit", limit},
			},
		}
	}
	return nil
}

func ValidateMaxString(attribute string, value string, limit int) *FieldError {
	if len(value) > limit {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} may not be greater than {limit} characters.",
			Placeholders: List{
				{"attribute", attribute},
				{"limit", limit},
			},
		}
	}
	return nil
}

func ValidateMinSlice(attribute string, value []any, limit int) *FieldError {
	if len(value) < limit {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} must have at least {limit} Entrys.",
			Placeholders: List{
				{"attribute", attribute},
				{"limit", limit},
			},
		}
	}
	return nil
}

func ValidateMaxSlice(attribute string, value []any, limit int) *FieldError {
	if len(value) > limit {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} may not have more than {limit} Entrys.",
			Placeholders: List{
				{"attribute", attribute},
				{"limit", limit},
			},
		}
	}
	return nil
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

func ValidateRequired(attribute string, value any) *FieldError {
	if isEmpty(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field is required.",
			Placeholders: List{
				{"attribute", attribute},
			},
		}
	}
	return nil
}

func ValidateRequiredIf[T comparable](attribute string, value any, other T, param string) *FieldError {
	comparisons := []string{">=", "<=", "!=", ">", "<", "=="}
	actual := fmt.Sprintf("%v", other)

	for _, op := range comparisons {
		if strings.Contains(param, op) {
			parts := strings.SplitN(param, op, 2)
			if len(parts) != 2 {
				return &FieldError{
					FieldName: attribute,
					Message:   "Invalid condition for required_if rule.",
					Placeholders: List{
						{"attribute", attribute},
					},
				}
			}

			key := strings.TrimSpace(parts[0])
			expected := strings.TrimSpace(parts[1])

			if isEmpty(value) {
				switch op {
				case "==":
					if actual == expected {
						return requiredIfError(attribute, key, expected)
					}
				case "!=":
					if actual != expected {
						return requiredIfError(attribute, key, expected)
					}
				case ">":
					if actual > expected {
						return requiredIfError(attribute, key, expected)
					}
				case "<":
					if actual < expected {
						return requiredIfError(attribute, key, expected)
					}
				case ">=":
					if actual >= expected {
						return requiredIfError(attribute, key, expected)
					}
				case "<=":
					if actual <= expected {
						return requiredIfError(attribute, key, expected)
					}
				}
			}
			return nil
		}
	}

	parts := strings.Split(param, ",")
	if len(parts) < 2 {
		return &FieldError{
			FieldName: attribute,
			Message:   "Invalid parameters for required_if rule.",
			Placeholders: List{
				{"attribute", attribute},
			},
		}
	}

	for _, expected := range parts[1:] {
		if actual == strings.TrimSpace(expected) && isEmpty(value) {
			return requiredIfError(attribute, parts[0], expected)
		}
	}
	return nil
}

func requiredIfError(attribute string, otherField string, otherValue string) *FieldError {
	return &FieldError{
		FieldName: attribute,
		Message:   "The {attribute} field is required when {other} is {value}.",
		Placeholders: List{
			{"attribute", attribute},
			{"other", otherField},
			{"value", otherValue},
		},
	}
}

func ValidateRequiredUnless[T comparable](attribute string, value any, other T, param string) *FieldError {
	comparisons := []string{">=", "<=", "!=", ">", "<", "=="}
	actual := fmt.Sprintf("%v", other)

	for _, op := range comparisons {
		if strings.Contains(param, op) {
			parts := strings.SplitN(param, op, 2)
			if len(parts) != 2 {
				return &FieldError{
					FieldName: attribute,
					Message:   "Invalid parameters for required_unless rule.",
					Placeholders: List{
						{"attribute", attribute},
					},
				}
			}

			key := strings.TrimSpace(parts[0])
			expected := strings.TrimSpace(parts[1])

			if isEmpty(value) {
				switch op {
				case "==":
					if actual != expected {
						return requiredUnlessError(attribute, key, expected)
					}
				case "!=":
					if actual == expected {
						return requiredUnlessError(attribute, key, expected)
					}
				case ">":
					if actual <= expected {
						return requiredUnlessError(attribute, key, expected)
					}
				case "<":
					if actual >= expected {
						return requiredUnlessError(attribute, key, expected)
					}
				case ">=":
					if actual < expected {
						return requiredUnlessError(attribute, key, expected)
					}
				case "<=":
					if actual > expected {
						return requiredUnlessError(attribute, key, expected)
					}
				}
			}
			return nil
		}
	}

	parts := strings.Split(param, ",")
	if len(parts) < 2 {
		return &FieldError{
			FieldName: attribute,
			Message:   "Invalid parameters for required_unless rule.",
			Placeholders: List{
				{"attribute", attribute},
			},
		}
	}

	for _, expected := range parts[1:] {
		if actual != strings.TrimSpace(expected) && isEmpty(value) {
			return requiredUnlessError(attribute, parts[0], expected)
		}
	}
	return nil
}

func requiredUnlessError(attribute, other, expected string) *FieldError {
	return &FieldError{
		FieldName: attribute,
		Message:   "The {attribute} field is required unless {other} is in {expected}.",
		Placeholders: List{
			{"attribute", attribute},
			{"other", other},
			{"expected", expected},
		},
	}
}

// WithoutAll verifica si el campo debe estar presente cuando todos los otros campos están vacíos
func ValidateWithoutAll(attribute string, value any, otherFieldNames []string, otherValues ...any) *FieldError {
	allEmpty := true
	for _, val := range otherValues {
		if !isEmpty(val) {
			allEmpty = false
			break
		}
	}

	if allEmpty && isEmpty(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field is required when none of {others} are present.",
			Placeholders: List{
				{"attribute", attribute},
				{"others", strings.Join(otherFieldNames, ", ")},
			},
		}
	}
	return nil
}

func ValidateWithout(attribute string, value any, otherFieldNames []string, otherValues ...any) *FieldError {
	anyEmpty := false
	for _, val := range otherValues {
		if isEmpty(val) {
			anyEmpty = true
			break
		}
	}

	if anyEmpty && isEmpty(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field is required when {others} is empty.",
			Placeholders: List{
				{"attribute", attribute},
				{"others", strings.Join(otherFieldNames, ", ")},
			},
		}
	}
	return nil
}

func ValidateWithAll(attribute string, value any, otherFieldNames []string, otherValues ...any) *FieldError {
	allFilled := true
	for _, val := range otherValues {
		if isEmpty(val) {
			allFilled = false
			break
		}
	}

	if allFilled && isEmpty(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field is required when all {others} are present.",
			Placeholders: List{
				{"attribute", attribute},
				{"others", strings.Join(otherFieldNames, ", ")},
			},
		}
	}
	return nil
}

func ValidateWith(attribute string, value any, otherFieldNames []string, otherValues ...any) *FieldError {
	anyFilled := false
	for _, val := range otherValues {
		if !isEmpty(val) {
			anyFilled = true
			break
		}
	}

	if anyFilled && isEmpty(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field is required when {others} is present.",
			Placeholders: List{
				{"attribute", attribute},
				{"others", strings.Join(otherFieldNames, ", ")},
			},
		}
	}
	return nil
}

func ValidateSame[T comparable](attribute string, value T, otherattribute string, other T) *FieldError {
	if value != other {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} and {other} must match.",
			Placeholders: List{
				{"attribute", attribute},
				{"other", otherattribute},
			},
		}
	}
	return nil
}

func ValidateDifferent[T comparable](attribute string, value T, otherattribute string, other T) *FieldError {
	if value == other {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} and {other} must be different.",
			Placeholders: List{
				{"attribute", attribute},
				{"other", otherattribute},
			},
		}
	}
	return nil
}

func ValidateConfirmed[T comparable](attribute string, value T, confirmation T) *FieldError {
	if value != confirmation {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} confirmation does not match.",
			Placeholders: List{
				{"attribute", attribute},
			},
		}
	}
	return nil
}

func ValidateAccepted(attribute string, value any) *FieldError {
	v := fmt.Sprintf("%v", value)
	acceptedValues := []string{"yes", "on", "1", "true"}
	for _, a := range acceptedValues {
		if strings.EqualFold(v, a) {
			return nil
		}
	}
	return &FieldError{
		FieldName: attribute,
		Message:   "The {attribute} must be accepted.",
		Placeholders: List{
			{"attribute", attribute},
		},
	}
}

func ValidateDeclined(attribute string, value any) *FieldError {
	v := fmt.Sprintf("%v", value)
	declinedValues := []string{"no", "off", "0", "false"}
	for _, d := range declinedValues {
		if strings.EqualFold(v, d) {
			return nil
		}
	}
	return &FieldError{
		FieldName: attribute,
		Message:   "The {attribute} must be declined.",
		Placeholders: List{
			{Key: "attribute", Value: attribute},
		},
	}
}

func ValidateDigits(attribute string, value any, length int) *FieldError {
	v := fmt.Sprintf("%v", value)
	if len(v) != length || !regexp.MustCompile(`^\d+$`).MatchString(v) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} must be :digits digits.",
			Placeholders: List{
				{"attribute", attribute},
				{"digits", length},
			},
		}
	}
	return nil
}

func ValidateDigitsBetween(attribute string, value any, min, max int) *FieldError {
	v := fmt.Sprintf("%v", value)
	length := len(v)
	if length < min || length > max || !regexp.MustCompile(`^\d+$`).MatchString(v) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} must be between :min and :max digits.",
			Placeholders: List{
				{"attribute", attribute},
				{"min", min},
				{"max", max},
			},
		}
	}
	return nil
}

func ValidateEmail(attribute, value string) *FieldError {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} must be a valid email address.",
			Placeholders: List{
				{"attribute", attribute},
			},
		}
	}
	return nil
}

func ValidateURL(attribute, value string) *FieldError {
	urlRegex := regexp.MustCompile(`^(https?://)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,6}(/[\w\-\./?%&=]*)?$`)
	if !urlRegex.MatchString(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} must be a valid URL.",
			Placeholders: List{
				{"attribute", attribute},
			},
		}
	}
	return nil
}

func ValidateUUID(attribute, value string) *FieldError {
	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89ab][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
	if !uuidRegex.MatchString(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} must be a valid UUID.",
			Placeholders: List{
				{"attribute", attribute},
			},
		}
	}
	return nil
}

func ValidateULID(attribute, value string) *FieldError {
	ulidRegex := regexp.MustCompile(`^[0123456789ABCDEFGHJKMNPQRSTVWXYZ]{26}$`)
	if !ulidRegex.MatchString(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} must be a valid ULID.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}

	tsStr := value[:10]
	decoder := base32.NewEncoding("0123456789ABCDEFGHJKMNPQRSTVWXYZ").WithPadding(base32.NoPadding)
	decoded, err := decoder.DecodeString(strings.ToUpper(tsStr))
	if err != nil || len(decoded) < 6 {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} has an invalid ULID timestamp.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}

	timestampMs := uint64(decoded[0])<<40 |
		uint64(decoded[1])<<32 |
		uint64(decoded[2])<<24 |
		uint64(decoded[3])<<16 |
		uint64(decoded[4])<<8 |
		uint64(decoded[5])

	if timestampMs > uint64(time.Now().UnixMilli()) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} has a ULID timestamp in the future.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}

	return nil
}

func ValidateIP(attribute, value string) *FieldError {
	if net.ParseIP(value) == nil {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} must be a valid IP address.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateIPv4(attribute, value string) *FieldError {
	ip := net.ParseIP(value)
	if ip == nil || ip.To4() == nil {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} must be a valid IPv4 address.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateIPv6(attribute, value string) *FieldError {
	ip := net.ParseIP(value)
	if ip == nil || ip.To4() != nil || strings.Contains(ip.String(), "::ffff:") {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} must be a valid IPv6 address.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateMACAddress(attribute, value string) *FieldError {
	if _, err := net.ParseMAC(value); err != nil {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} must be a valid MAC address.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateASCII(attribute, value string) *FieldError {
	if !regexp.MustCompile(`^[\x00-\x7F]+$`).MatchString(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field must only contain ASCII characters.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateLowercase(attribute, value string) *FieldError {
	if value != strings.ToLower(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field must be lowercase.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateUppercase(attribute, value string) *FieldError {
	if value != strings.ToUpper(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field must be uppercase.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateHex(attribute, value string) *FieldError {
	if !regexp.MustCompile(`^[0-9a-fA-F]+$`).MatchString(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field must be a hexadecimal string.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateHexColor(attribute, value string) *FieldError {
	if !regexp.MustCompile(`^#(?:[0-9a-fA-F]{3}){1,2}$`).MatchString(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field must be a valid hexadecimal color code.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateJSON(attribute string, value string) *FieldError {
	var js map[string]interface{}
	if err := json.Unmarshal([]byte(value), &js); err != nil {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field must be a valid JSON string.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateSlug(attribute string, value string) *FieldError {
	slugRegex := regexp.MustCompile(`^[a-z0-9]+(?:[-_][a-z0-9]+)*$`)
	if !slugRegex.MatchString(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field must be a valid slug (lowercase letters, numbers, hyphens and underscores only).",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateRegex(attribute string, value string, pattern string) *FieldError {

	re, err := regexp.Compile(pattern)
	if err != nil {
		PrintWarning("Invalid regular expression pattern: :pattern", Entry{"pattern", pattern})
		return &FieldError{
			FieldName: attribute,
			Message:   "Invalid regular expression pattern.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}

	}
	if !re.MatchString(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field format is invalid.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateNotRegex(attribute string, value string, pattern string) *FieldError {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return &FieldError{
			FieldName: attribute,
			Message:   "Invalid regular expression pattern.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	if re.MatchString(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field contains an invalid value.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateAlpha(attribute string, value string) *FieldError {
	alphaRegex := regexp.MustCompile(`^[a-zA-Z]+$`)
	if !alphaRegex.MatchString(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field may only contain letters.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateAlphaDash(attribute string, value string) *FieldError {
	regex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !regex.MatchString(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field may only contain letters, numbers, dashes and underscores.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateAlphaSpaces(attribute string, value string) *FieldError {
	regex := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	if !regex.MatchString(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field may only contain letters and spaces.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateAlphaDashSpaces(attribute string, value string) *FieldError {
	regex := regexp.MustCompile(`^[a-zA-Z0-9_\s-]+$`)
	if !regex.MatchString(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field may only contain letters, numbers, spaces, dashes and underscores.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateAlphaNum(attribute string, value string) *FieldError {
	regex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !regex.MatchString(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field may only contain letters and numbers.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateAlphaNumDash(attribute string, value string) *FieldError {
	regex := regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
	if !regex.MatchString(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field may only contain letters, numbers and dashes.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateAlphaNumSpaces(attribute string, value string) *FieldError {
	regex := regexp.MustCompile(`^[a-zA-Z0-9\s]+$`)
	if !regex.MatchString(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field may only contain letters, numbers, and spaces.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateAlphaNumDashSpaces(attribute string, value string) *FieldError {
	regex := regexp.MustCompile(`^[a-zA-Z0-9_\s-]+$`)
	if !regex.MatchString(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field may only contain letters, numbers, spaces, dashes, and underscores.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateAlphaAccents(attribute string, value string) *FieldError {
	regex := regexp.MustCompile(`^[a-zA-ZáéíóúÁÉÍÓÚñÑüÜ]+$`)
	if !regex.MatchString(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field may only contain letters, including accented characters and ñ.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateAlphaDashAccents(attribute string, value string) *FieldError {
	regex := regexp.MustCompile(`^[a-zA-ZáéíóúÁÉÍÓÚñÑüÜ_-]+$`)
	if !regex.MatchString(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field may only contain letters (including accented characters and ñ), dashes, and underscores.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateAlphaSpacesAccents(attribute string, value string) *FieldError {
	regex := regexp.MustCompile(`^[a-zA-ZáéíóúÁÉÍÓÚñÑüÜ\s]+$`)
	if !regex.MatchString(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field may only contain letters, accented characters, ñ, and spaces.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateAlphaDashSpacesAccents(attribute string, value string) *FieldError {
	regex := regexp.MustCompile(`^[a-zA-ZáéíóúÁÉÍÓÚñÑüÜ_\s-]+$`)
	if !regex.MatchString(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field may only contain letters (including accented characters and ñ), numbers, spaces, dashes, and underscores.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateAlphaNumAccents(attribute string, value string) *FieldError {
	regex := regexp.MustCompile(`^[a-zA-Z0-9áéíóúÁÉÍÓÚñÑüÜ]+$`)
	if !regex.MatchString(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field may only contain letters (including accented characters and ñ) and numbers.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateAlphaNumDashAccents(attribute string, value string) *FieldError {
	regex := regexp.MustCompile(`^[a-zA-Z0-9áéíóúÁÉÍÓÚñÑüÜ_-]+$`)
	if !regex.MatchString(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field may only contain letters (including accented characters and ñ), numbers, dashes, and underscores.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateAlphaNumSpacesAccents(attribute string, value string) *FieldError {
	regex := regexp.MustCompile(`^[a-zA-Z0-9áéíóúÁÉÍÓÚñÑüÜ\s]+$`)
	if !regex.MatchString(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field may only contain letters (including accented characters and ñ), numbers, and spaces.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateAlphaNumDashSpacesAccents(attribute string, value string) *FieldError {
	regex := regexp.MustCompile(`^[a-zA-Z0-9áéíóúÁÉÍÓÚñÑüÜ_\s-]+$`)
	if !regex.MatchString(value) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field may only contain letters (including accented characters and ñ), numbers, spaces, dashes, and underscores.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

// isAlphaNumeric es una funcion auxiliar no hace parte de las validaciones
func isAlphaNumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

func ValidateUsername(attribute string, value string) *FieldError {
	if len(value) == 0 || !isAlphaNumeric(rune(value[0])) || !isAlphaNumeric(rune(value[len(value)-1])) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} must start and end with a letter or number.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}

	for i := 0; i < len(value)-1; i++ {
		if (value[i] == '.' || value[i] == '_') && (value[i+1] == '.' || value[i+1] == '_') {
			return &FieldError{
				FieldName: attribute,
				Message:   "The {attribute} cannot contain consecutive dots or underscores.",
				Placeholders: List{
					{Key: "attribute", Value: attribute},
				},
			}
		}
	}

	for _, char := range value {
		if !isAlphaNumeric(char) && char != '.' && char != '_' {
			return &FieldError{
				FieldName: attribute,
				Message:   "The {attribute} may only contain letters, numbers, dots, and underscores.",
				Placeholders: List{
					{Key: "attribute", Value: attribute},
				},
			}
		}
	}

	return nil
}

func ValidateStartsWith(attribute string, value string, prefix string) *FieldError {
	if !strings.HasPrefix(value, prefix) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} must start with {prefix}.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
				{Key: "prefix", Value: prefix},
			},
		}
	}
	return nil
}

func ValidateEndsWith(attribute string, value string, suffix string) *FieldError {
	if !strings.HasSuffix(value, suffix) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} must end with {suffix}.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
				{Key: "suffix", Value: suffix},
			},
		}
	}
	return nil
}

func ValidateContains(attribute string, value string, substr string) *FieldError {
	if !strings.Contains(value, substr) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field must contain {substring}.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
				{Key: "substring", Value: substr},
			},
		}
	}
	return nil
}

func ValidateNotContains(attribute string, value string, substr string) *FieldError {
	if strings.Contains(value, substr) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field must not contain {substring}.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
				{Key: "substring", Value: substr},
			},
		}
	}
	return nil
}

func ValidateIn[T comparable](attribute string, value T, allowed ...T) *FieldError {
	for _, v := range allowed {
		if value == v {
			return nil
		}
	}
	return &FieldError{
		FieldName: attribute,
		Message:   "The selected {attribute} is invalid.",
		Placeholders: List{
			{Key: "attribute", Value: attribute},
			{Key: "options", Value: allowed},
		},
	}
}

func ValidateNin[T comparable](attribute string, value T, denied ...T) *FieldError {
	for _, v := range denied {
		if value == v {
			return &FieldError{
				FieldName: attribute,
				Message:   "The selected {attribute} is not allowed.",
				Placeholders: List{
					{Key: "attribute", Value: attribute},
					{Key: "restricted", Value: denied},
				},
			}
		}
	}
	return nil
}

func ValidateExists(attribute string, collection string, field string, value any) *FieldError {
	result := map[string]any{}

	err := DB.Collection(collection).FindOne(context.TODO(), bson.D{
		{Key: field, Value: value},
	}).Decode(&result)

	if err == mongo.ErrNoDocuments {
		// No existe → error
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} does not exist.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}

	if err != nil {
		// Error inesperado en la consulta
		PrintWarning("Failed to check existence in database: collection: :collection value: :value error: :error ",
			Entry{"error", err.Error()},
			Entry{"collection", collection},
			Entry{"value", value},
		)
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} failed to check existence in database.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}

	return nil
}

func ValidateUnique(attribute string, collection string, field string, value any, currentID string) *FieldError {
	result := map[string]any{}
	var id bson.ObjectID
	if currentID != "" {
		var er error
		id, er = bson.ObjectIDFromHex(currentID)
		if er != nil {
			PrintWarning("Failed to convert string [:input_id] to ObjectID :error ", Entry{"error", er.Error()}, Entry{"input_id", currentID})
		}
	}

	err := DB.Collection(collection).FindOne(context.TODO(), bson.D{
		{Key: field, Value: value},
		{Key: "_id", Value: bson.D{{Key: "$ne", Value: id}}},
	}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	if err != nil {
		PrintWarning("Failed to find document in database: collection: :collection value: :value error: :error ", Entry{"error", err.Error()}, Entry{"collection", collection}, Entry{"value", value})
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} failed to find document in database.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return &FieldError{
		FieldName: attribute,
		Message:   "The {attribute} has already been taken.",
		Placeholders: List{
			{Key: "attribute", Value: attribute},
		},
	}
}

func ValidateUniqueIn[T comparable](attribute string, list []T) *FieldError {
	seen := make(map[T]bool)
	for _, Entry := range list {
		if seen[Entry] {
			return &FieldError{
				FieldName: attribute,
				Message:   "The {attribute} field has a duplicate value.",
				Placeholders: List{
					{Key: "attribute", Value: attribute},
					{Key: "value", Value: Entry},
				},
			}
		}
		seen[Entry] = true
	}
	return nil
}

func ValidatePositive[T constraints.Integer | constraints.Float](attribute string, value T) *FieldError {
	if value <= 0 {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field must be greater than 0.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateNegative[T constraints.Integer | constraints.Float](attribute string, value T) *FieldError {
	if value >= 0 {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field must be less than 0.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateBetween[T constraints.Integer | constraints.Float](attribute string, value T, min T, max T) *FieldError {
	if value < min || value > max {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field must be between {min} and {max}.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
				{Key: "min", Value: min},
				{Key: "max", Value: max},
			},
		}
	}
	return nil
}

func ValidateBefore(attribute string, value time.Time, target time.Time) *FieldError {
	if !value.Before(target) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field must be a date before {date}.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
				{Key: "date", Value: target.Format(time.RFC3339)},
			},
		}
	}
	return nil
}

func ValidateAfter(attribute string, value time.Time, target time.Time) *FieldError {
	if !value.After(target) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field must be a date after {date}.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
				{Key: "date", Value: target.Format(time.RFC3339)},
			},
		}
	}
	return nil
}

func ValidateBeforeNow(attribute string, value time.Time) *FieldError {
	if value.After(time.Now()) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field must be a date in the past.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateAfterNow(attribute string, value time.Time) *FieldError {
	if value.Before(time.Now()) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field must be a date in the future.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
			},
		}
	}
	return nil
}

func ValidateDateBetween(attribute string, value time.Time, start time.Time, end time.Time) *FieldError {
	if value.Before(start) || value.After(end) {
		return &FieldError{
			FieldName: attribute,
			Message:   "The {attribute} field must be a date between {start} and {end}.",
			Placeholders: List{
				{Key: "attribute", Value: attribute},
				{Key: "start", Value: start.Format(time.RFC3339)},
				{Key: "end", Value: end.Format(time.RFC3339)},
			},
		}
	}
	return nil
}
