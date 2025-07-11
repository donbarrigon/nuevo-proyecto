package model

import (
	"reflect"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
)

// Fill llena los campos del modelo con los valores del request,
// pero solo si el campo del modelo tiene la etiqueta fillable
func Fill(model any, request any) app.Error {
	modelValue := reflect.ValueOf(model)
	requestValue := reflect.ValueOf(request)

	if modelValue.Kind() != reflect.Ptr || requestValue.Kind() != reflect.Ptr {
		return app.Errors.Unknownf("The parameters model and request must be pointers")
	}

	modelValue = modelValue.Elem()
	requestValue = requestValue.Elem()

	if modelValue.Kind() != reflect.Struct || requestValue.Kind() != reflect.Struct {
		return app.Errors.Unknownf("The parameters model and request must be structs")
	}

	modelType := modelValue.Type()

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		// if fillable, ok := field.Tag.Lookup("fillable"); ok && fillable == "true" {
		if _, ok := field.Tag.Lookup("fillable"); ok {
			fieldName := field.Name

			requestField := requestValue.FieldByName(fieldName)

			if requestField.IsValid() && requestField.Type().AssignableTo(field.Type) {
				modelField := modelValue.Field(i)
				if modelField.CanSet() {
					modelField.Set(requestField)
				}
			}
		}
	}
	return nil
}
