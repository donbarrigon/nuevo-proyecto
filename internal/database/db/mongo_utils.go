package db

import (
	"reflect"
	"strings"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
)

// Fill copia los valores de request a model si tienen el mismo nombre y tipo compatible
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

	requestType := requestValue.Type()

	for i := 0; i < requestType.NumField(); i++ {
		requestField := requestType.Field(i)
		fieldName := requestField.Name
		requestFieldValue := requestValue.Field(i)

		// Solo procesamos campos exportados (que empiecen con mayúscula)
		if !requestField.IsExported() {
			continue
		}

		// Buscamos el campo en el modelo
		modelField := modelValue.FieldByName(fieldName)

		if !modelField.IsValid() || !modelField.CanSet() {
			continue
		}

		// Verificamos que no sea un valor zero (opcional)
		if requestFieldValue.IsZero() {
			continue // Puedes cambiar esto si quieres copiar valores zero
		}

		// Verificamos compatibilidad de tipos
		if requestFieldValue.Type().AssignableTo(modelField.Type()) {
			modelField.Set(requestFieldValue)
		} else if requestFieldValue.Type().ConvertibleTo(modelField.Type()) {
			// Intentamos conversión si es posible (ej: int32 -> int64)
			modelField.Set(requestFieldValue.Convert(modelField.Type()))
		}
	}
	return nil
}

// Fill copia los valores de request a model si tienen el mismo nombre y tipo compatible
// Retorna un mapa con los campos que cambiaron, usando el tag bson como clave
func FillDirty(model any, request any) (map[string]any, app.Error) {
	dirty := make(map[string]any)

	modelValue := reflect.ValueOf(model)
	requestValue := reflect.ValueOf(request)

	if modelValue.Kind() != reflect.Ptr || requestValue.Kind() != reflect.Ptr {
		return nil, app.Errors.Unknownf("The parameters model and request must be pointers")
	}

	modelValue = modelValue.Elem()
	requestValue = requestValue.Elem()

	if modelValue.Kind() != reflect.Struct || requestValue.Kind() != reflect.Struct {
		return nil, app.Errors.Unknownf("The parameters model and request must be structs")
	}

	modelType := modelValue.Type()
	requestType := requestValue.Type()

	for i := 0; i < requestType.NumField(); i++ {
		requestField := requestType.Field(i)
		fieldName := requestField.Name
		requestFieldValue := requestValue.Field(i)

		// Solo procesamos campos exportados (que empiecen con mayúscula)
		if !requestField.IsExported() {
			continue
		}

		// Buscamos el campo en el modelo
		modelField := modelValue.FieldByName(fieldName)

		if !modelField.IsValid() || !modelField.CanSet() {
			continue
		}

		// Verificamos que no sea un valor zero (opcional)
		if requestFieldValue.IsZero() {
			continue
		}

		// Verificamos compatibilidad de tipos
		var newValue reflect.Value
		if requestFieldValue.Type().AssignableTo(modelField.Type()) {
			newValue = requestFieldValue
		} else if requestFieldValue.Type().ConvertibleTo(modelField.Type()) {
			newValue = requestFieldValue.Convert(modelField.Type())
		} else {
			continue // Tipos incompatibles, saltamos este campo
		}

		// Comparamos el valor actual con el nuevo valor
		currentValue := modelField
		if !reflect.DeepEqual(currentValue.Interface(), newValue.Interface()) {
			// Los valores son diferentes, actualizamos el modelo y guardamos en dirty

			// Obtenemos el tag bson del campo del modelo
			modelFieldStruct, found := modelType.FieldByName(fieldName)
			if !found {
				continue
			}

			bsonTag := modelFieldStruct.Tag.Get("bson")
			if bsonTag == "" {
				// Si no tiene tag bson, usamos el nombre del campo en minúsculas
				bsonTag = strings.ToLower(fieldName)
			} else {
				// Si tiene tag bson, extraemos solo el nombre (sin opciones como ",omitempty")
				if commaIndex := strings.Index(bsonTag, ","); commaIndex != -1 {
					bsonTag = bsonTag[:commaIndex]
				}
			}

			// Guardamos el nuevo valor en dirty
			dirty[bsonTag] = newValue.Interface()

			// Actualizamos el campo en el modelo
			modelField.Set(newValue)
		}
	}

	return dirty, nil
}

// // Versión alternativa que también incluye una función helper para extraer el tag bson
// func FillWithHelper(model any, request any) (map[string]any, app.Error) {
// 	dirty := make(map[string]any)

// 	modelValue := reflect.ValueOf(model)
// 	requestValue := reflect.ValueOf(request)

// 	if modelValue.Kind() != reflect.Ptr || requestValue.Kind() != reflect.Ptr {
// 		return nil, app.Errors.Unknownf("The parameters model and request must be pointers")
// 	}

// 	modelValue = modelValue.Elem()
// 	requestValue = requestValue.Elem()

// 	if modelValue.Kind() != reflect.Struct || requestValue.Kind() != reflect.Struct {
// 		return nil, app.Errors.Unknownf("The parameters model and request must be structs")
// 	}

// 	modelType := modelValue.Type()
// 	requestType := requestValue.Type()

// 	for i := 0; i < requestType.NumField(); i++ {
// 		requestField := requestType.Field(i)
// 		fieldName := requestField.Name
// 		requestFieldValue := requestValue.Field(i)

// 		if !requestField.IsExported() {
// 			continue
// 		}

// 		modelField := modelValue.FieldByName(fieldName)

// 		if !modelField.IsValid() || !modelField.CanSet() {
// 			continue
// 		}

// 		if requestFieldValue.IsZero() {
// 			continue
// 		}

// 		var newValue reflect.Value
// 		if requestFieldValue.Type().AssignableTo(modelField.Type()) {
// 			newValue = requestFieldValue
// 		} else if requestFieldValue.Type().ConvertibleTo(modelField.Type()) {
// 			newValue = requestFieldValue.Convert(modelField.Type())
// 		} else {
// 			continue
// 		}

// 		// Comparamos valores
// 		if !reflect.DeepEqual(modelField.Interface(), newValue.Interface()) {
// 			modelFieldStruct, found := modelType.FieldByName(fieldName)
// 			if !found {
// 				continue
// 			}

// 			bsonTag := getBsonFieldName(modelFieldStruct, fieldName)
// 			dirty[bsonTag] = newValue.Interface()
// 			modelField.Set(newValue)
// 		}
// 	}

// 	return dirty, nil
// }

// // Helper function para extraer el nombre del campo bson
// func getBsonFieldName(field reflect.StructField, defaultName string) string {
// 	bsonTag := field.Tag.Get("bson")
// 	if bsonTag == "" {
// 		return strings.ToLower(defaultName)
// 	}

// 	// Extraemos solo el nombre, ignorando opciones como ",omitempty"
// 	if commaIndex := strings.Index(bsonTag, ","); commaIndex != -1 {
// 		return bsonTag[:commaIndex]
// 	}

// 	return bsonTag
// }
