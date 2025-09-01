package app

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func InitMongoDB() error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(Env.DB_CONNECTION_STRING).SetServerAPIOptions(serverAPI)
	clientOptions.SetMaxPoolSize(100)
	clientOptions.SetMinPoolSize(5)
	clientOptions.SetRetryWrites(true)
	clientOptions.SetTimeout(30 * time.Second)

	var err error

	DBClient, err = mongo.Connect(clientOptions)
	if err != nil {
		log.Fatalf("Error al conectar con mongodb: %v", err)
		return err
	}
	DB = DBClient.Database(Env.DB_DATABASE)

	PrintInfo("Conectado exitosamente a MongoDB: :con - Base de datos: :db",
		Entry{"con", Env.DB_CONNECTION_STRING},
		Entry{"con", Env.DB_DATABASE})
	return nil
}

func CloseMongoDB() error {

	if DBClient == nil {
		return nil
	}

	err := DBClient.Disconnect(context.TODO())
	if err != nil {
		e := fmt.Errorf("error al cerrar la conexión con MongoDB: %w", err)
		Print(e.Error())
		return e
	}

	return nil
}

// FillChanges compara model con request y retorna los valores antiguos y los nuevos.
// Usa el tag bson como clave.
// Además actualiza el model con los valores nuevos.
// @return original, dirty, error
func Fill(model any, request any) (map[string]any, map[string]any, Error) {

	original := map[string]any{}
	dirty := map[string]any{}

	modelValue := reflect.ValueOf(model)
	requestValue := reflect.ValueOf(request)

	if modelValue.Kind() != reflect.Ptr || requestValue.Kind() != reflect.Ptr {
		return original, dirty, Errors.Unknownf("The parameters model and request must be pointers")
	}

	modelValue = modelValue.Elem()
	requestValue = requestValue.Elem()

	if modelValue.Kind() != reflect.Struct || requestValue.Kind() != reflect.Struct {
		return original, dirty, Errors.Unknownf("The parameters model and request must be structs")
	}

	modelType := modelValue.Type()
	requestType := requestValue.Type()

	for i := 0; i < requestType.NumField(); i++ {
		requestField := requestType.Field(i)
		fieldName := requestField.Name
		requestFieldValue := requestValue.Field(i)

		if !requestField.IsExported() {
			continue
		}

		modelField := modelValue.FieldByName(fieldName)
		if !modelField.IsValid() || !modelField.CanSet() {
			continue
		}

		if requestFieldValue.IsZero() {
			continue
		}

		var newValue reflect.Value
		if modelField.Type() == reflect.TypeOf(bson.ObjectID{}) &&
			requestFieldValue.Kind() == reflect.String {

			oid, err := bson.ObjectIDFromHex(requestFieldValue.String())
			if err != nil {
				continue
			}
			newValue = reflect.ValueOf(oid)

		} else if requestFieldValue.Type().AssignableTo(modelField.Type()) {
			newValue = requestFieldValue
		} else if requestFieldValue.Type().ConvertibleTo(modelField.Type()) {
			newValue = requestFieldValue.Convert(modelField.Type())
		} else {
			continue
		}

		currentValue := modelField
		if !reflect.DeepEqual(currentValue.Interface(), newValue.Interface()) {
			// sacar el tag bson
			modelFieldStruct, found := modelType.FieldByName(fieldName)
			if !found {
				continue
			}

			bsonTag := modelFieldStruct.Tag.Get("bson")
			if bsonTag == "" {
				bsonTag = strings.ToLower(fieldName)
			} else if commaIndex := strings.Index(bsonTag, ","); commaIndex != -1 {
				bsonTag = bsonTag[:commaIndex]
			}

			// Guardamos viejo y nuevo
			original[bsonTag] = currentValue.Interface()
			dirty[bsonTag] = newValue.Interface()

			// Actualizamos el modelo
			modelField.Set(newValue)
		}
	}

	return original, dirty, nil
}

// func Fill(model any, request any) Error {
// 	modelValue := reflect.ValueOf(model)
// 	requestValue := reflect.ValueOf(request)

// 	if modelValue.Kind() != reflect.Pointer || requestValue.Kind() != reflect.Pointer {
// 		return Errors.Unknownf("The parameters model and request must be pointers")
// 	}

// 	modelValue = modelValue.Elem()
// 	requestValue = requestValue.Elem()

// 	if modelValue.Kind() != reflect.Struct || requestValue.Kind() != reflect.Struct {
// 		return Errors.Unknownf("The parameters model and request must be structs")
// 	}

// 	requestType := requestValue.Type()

// 	for i := 0; i < requestType.NumField(); i++ {
// 		requestField := requestType.Field(i)
// 		fieldName := requestField.Name
// 		requestFieldValue := requestValue.Field(i)

// 		// Solo procesamos campos exportados (que empiecen con mayúscula)
// 		if !requestField.IsExported() {
// 			continue
// 		}

// 		// Buscamos el campo en el modelo
// 		modelField := modelValue.FieldByName(fieldName)

// 		if !modelField.IsValid() || !modelField.CanSet() {
// 			continue
// 		}

// 		// Verificamos que no sea un valor zero (opcional)
// 		if requestFieldValue.IsZero() {
// 			continue // Puedes cambiar esto si quieres copiar valores zero
// 		}

// 		// --- 🔥 NUEVO: si el destino es ObjectID y la fuente es string
// 		if modelField.Type() == reflect.TypeOf(bson.ObjectID{}) &&
// 			requestFieldValue.Kind() == reflect.String {

// 			oid, err := bson.ObjectIDFromHex(requestFieldValue.String())
// 			if err != nil {
// 				continue // si no es un ObjectID válido, lo saltamos
// 			}
// 			modelField.Set(reflect.ValueOf(oid))
// 			continue
// 		}

// 		// Verificamos compatibilidad de tipos
// 		if requestFieldValue.Type().AssignableTo(modelField.Type()) {
// 			modelField.Set(requestFieldValue)
// 		} else if requestFieldValue.Type().ConvertibleTo(modelField.Type()) {
// 			// Intentamos conversión si es posible (ej: int32 -> int64)
// 			modelField.Set(requestFieldValue.Convert(modelField.Type()))
// 		}
// 	}
// 	return nil
// }

// // Fill copia los valores de request a model si tienen el mismo nombre y tipo compatible
// // Retorna un mapa con los campos que cambiaron, usando el tag bson como clave
// func FillDirty(model any, request any) (map[string]any, Error) {
// 	dirty := make(map[string]any)

// 	modelValue := reflect.ValueOf(model)
// 	requestValue := reflect.ValueOf(request)

// 	if modelValue.Kind() != reflect.Ptr || requestValue.Kind() != reflect.Ptr {
// 		return nil, Errors.Unknownf("The parameters model and request must be pointers")
// 	}

// 	modelValue = modelValue.Elem()
// 	requestValue = requestValue.Elem()

// 	if modelValue.Kind() != reflect.Struct || requestValue.Kind() != reflect.Struct {
// 		return nil, Errors.Unknownf("The parameters model and request must be structs")
// 	}

// 	modelType := modelValue.Type()
// 	requestType := requestValue.Type()

// 	for i := 0; i < requestType.NumField(); i++ {
// 		requestField := requestType.Field(i)
// 		fieldName := requestField.Name
// 		requestFieldValue := requestValue.Field(i)

// 		// Solo procesamos campos exportados (que empiecen con mayúscula)
// 		if !requestField.IsExported() {
// 			continue
// 		}

// 		// Buscamos el campo en el modelo
// 		modelField := modelValue.FieldByName(fieldName)
// 		if !modelField.IsValid() || !modelField.CanSet() {
// 			continue
// 		}

// 		// Verificamos que no sea un valor zero (opcional)
// 		if requestFieldValue.IsZero() {
// 			continue
// 		}

// 		var newValue reflect.Value

// 		// --- 🔥 NUEVO: si el destino es ObjectID y la fuente es string
// 		if modelField.Type() == reflect.TypeOf(bson.ObjectID{}) &&
// 			requestFieldValue.Kind() == reflect.String {

// 			oid, err := bson.ObjectIDFromHex(requestFieldValue.String())
// 			if err != nil {
// 				continue // no es un ObjectID válido, ignoramos
// 			}
// 			newValue = reflect.ValueOf(oid)

// 		} else if requestFieldValue.Type().AssignableTo(modelField.Type()) {
// 			newValue = requestFieldValue
// 		} else if requestFieldValue.Type().ConvertibleTo(modelField.Type()) {
// 			newValue = requestFieldValue.Convert(modelField.Type())
// 		} else {
// 			continue // Tipos incompatibles, saltamos este campo
// 		}

// 		// Comparamos el valor actual con el nuevo valor
// 		currentValue := modelField
// 		if !reflect.DeepEqual(currentValue.Interface(), newValue.Interface()) {
// 			// Los valores son diferentes, actualizamos el modelo y guardamos en dirty

// 			// Obtenemos el tag bson del campo del modelo
// 			modelFieldStruct, found := modelType.FieldByName(fieldName)
// 			if !found {
// 				continue
// 			}

// 			bsonTag := modelFieldStruct.Tag.Get("bson")
// 			if bsonTag == "" {
// 				// Si no tiene tag bson, usamos el nombre del campo en minúsculas
// 				bsonTag = strings.ToLower(fieldName)
// 			} else {
// 				// Si tiene tag bson, extraemos solo el nombre (sin opciones como ",omitempty")
// 				if commaIndex := strings.Index(bsonTag, ","); commaIndex != -1 {
// 					bsonTag = bsonTag[:commaIndex]
// 				}
// 			}

// 			// Guardamos el nuevo valor en dirty
// 			dirty[bsonTag] = newValue.Interface()

// 			// Actualizamos el campo en el modelo
// 			modelField.Set(newValue)
// 		}
// 	}

// 	return dirty, nil
// }

// func FillByMap(model any, request map[string]any) Error {
// 	// Verificar que model sea un puntero a struct
// 	modelValue := reflect.ValueOf(model)
// 	if modelValue.Kind() != reflect.Ptr || modelValue.IsNil() {
// 		return Errors.InternalServerErrorf("model debe ser un puntero no nulo")
// 	}

// 	modelElem := modelValue.Elem()
// 	if modelElem.Kind() != reflect.Struct {
// 		return Errors.InternalServerErrorf("model debe ser un puntero a struct")
// 	}

// 	modelType := modelElem.Type()

// 	// Iterar sobre los campos del struct
// 	for i := 0; i < modelElem.NumField(); i++ {
// 		field := modelElem.Field(i)
// 		fieldType := modelType.Field(i)

// 		// Verificar si el campo es exportado (público)
// 		if !field.CanSet() {
// 			continue
// 		}

// 		// Obtener el tag bson
// 		bsonTag := fieldType.Tag.Get("bson")
// 		if bsonTag == "" || bsonTag == "_id" {
// 			continue // Saltar campos sin tag bson o con _id
// 		}

// 		// Buscar el valor en el map
// 		mapValue, exists := request[bsonTag]
// 		if !exists {
// 			continue
// 		}

// 		// Asignar el valor al campo
// 		if err := setFieldValue(field, mapValue); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func setFieldValue(field reflect.Value, value any) Error {
// 	if value == nil {
// 		return nil
// 	}

// 	valueReflect := reflect.ValueOf(value)
// 	fieldType := field.Type()
// 	valueType := valueReflect.Type()

// 	// Si los tipos son iguales, asignar directamente
// 	if valueType.AssignableTo(fieldType) {
// 		field.Set(valueReflect)
// 		return nil
// 	}

// 	// Intentar conversiones de tipos comunes
// 	switch fieldType.Kind() {
// 	case reflect.String:
// 		switch v := value.(type) {
// 		case string:
// 			field.SetString(v)
// 		case []byte:
// 			field.SetString(string(v))
// 		default:
// 			field.SetString(fmt.Sprintf("%v", value))
// 		}

// 	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
// 		switch v := value.(type) {
// 		case int:
// 			field.SetInt(int64(v))
// 		case int32:
// 			field.SetInt(int64(v))
// 		case int64:
// 			field.SetInt(v)
// 		case float32:
// 			field.SetInt(int64(v))
// 		case float64:
// 			field.SetInt(int64(v))
// 		case string:
// 			if intVal, err := strconv.ParseInt(v, 10, 64); err == nil {
// 				field.SetInt(intVal)
// 			} else {
// 				return Errors.InternalServerErrorf("no se puede convertir string ':value' a int", Entry{"value", v})
// 			}
// 		default:
// 			return Errors.InternalServerErrorf("no se puede convertir :value a :type", Entry{"value", value}, Entry{"type", fieldType.Kind()})
// 		}

// 	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
// 		switch v := value.(type) {
// 		case int:
// 			if v >= 0 {
// 				field.SetUint(uint64(v))
// 			} else {
// 				return Errors.InternalServerErrorf("valor negativo no puede ser asignado a campo unsigned")
// 			}
// 		case uint:
// 			field.SetUint(uint64(v))
// 		case uint32:
// 			field.SetUint(uint64(v))
// 		case uint64:
// 			field.SetUint(v)
// 		case float32:
// 			if v >= 0 {
// 				field.SetUint(uint64(v))
// 			} else {
// 				return Errors.InternalServerErrorf("valor negativo no puede ser asignado a campo unsigned")
// 			}
// 		case float64:
// 			if v >= 0 {
// 				field.SetUint(uint64(v))
// 			} else {
// 				return Errors.InternalServerErrorf("valor negativo no puede ser asignado a campo unsigned")
// 			}
// 		case string:
// 			if uintVal, err := strconv.ParseUint(v, 10, 64); err == nil {
// 				field.SetUint(uintVal)
// 			} else {
// 				return Errors.InternalServerErrorf("no se puede convertir string ':value' a uint", Entry{"value", v})
// 			}
// 		default:
// 			return Errors.InternalServerErrorf("no se puede convertir :value a :type", Entry{"value", value}, Entry{"type", fieldType.Kind()})
// 		}

// 	case reflect.Float32, reflect.Float64:
// 		switch v := value.(type) {
// 		case float32:
// 			field.SetFloat(float64(v))
// 		case float64:
// 			field.SetFloat(v)
// 		case int:
// 			field.SetFloat(float64(v))
// 		case int32:
// 			field.SetFloat(float64(v))
// 		case int64:
// 			field.SetFloat(float64(v))
// 		case string:
// 			if floatVal, err := strconv.ParseFloat(v, 64); err == nil {
// 				field.SetFloat(floatVal)
// 			} else {
// 				return Errors.InternalServerErrorf("no se puede convertir string ':value' a float", Entry{"value", v})
// 			}
// 		default:
// 			return Errors.InternalServerErrorf("no se puede convertir :value a :type", Entry{"value", value}, Entry{"type", fieldType.Kind()})
// 		}

// 	case reflect.Bool:
// 		switch v := value.(type) {
// 		case bool:
// 			field.SetBool(v)
// 		case string:
// 			if boolVal, err := strconv.ParseBool(v); err == nil {
// 				field.SetBool(boolVal)
// 			} else {
// 				return Errors.InternalServerErrorf("no se puede convertir string ':value' a bool", Entry{"value", v})
// 			}
// 		case int:
// 			field.SetBool(v != 0)
// 		default:
// 			return Errors.InternalServerErrorf("no se puede convertir :value a bool", Entry{"value", value})
// 		}

// 	case reflect.Ptr:
// 		if field.IsNil() {
// 			field.Set(reflect.New(fieldType.Elem()))
// 		}
// 		return setFieldValue(field.Elem(), value)

// 	case reflect.Slice:
// 		if valueReflect.Kind() == reflect.Slice {
// 			newSlice := reflect.MakeSlice(fieldType, valueReflect.Len(), valueReflect.Cap())
// 			for i := 0; i < valueReflect.Len(); i++ {
// 				elem := newSlice.Index(i)
// 				if err := setFieldValue(elem, valueReflect.Index(i).Interface()); err != nil {
// 					return err
// 				}
// 			}
// 			field.Set(newSlice)
// 		} else {
// 			return Errors.InternalServerErrorf("no se puede convertir :value a slice", Entry{"value", value})
// 		}

// 	default:
// 		// Para tipos que se pueden convertir directamente
// 		if valueReflect.Type().ConvertibleTo(fieldType) {
// 			field.Set(valueReflect.Convert(fieldType))
// 		} else {
// 			return Errors.InternalServerErrorf("no se puede convertir :value a :type", Entry{"value", value}, Entry{"type", fieldType})
// 		}
// 	}

// 	return nil
// }
