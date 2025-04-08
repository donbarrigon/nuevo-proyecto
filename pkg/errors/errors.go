package errors

import (
	"context"
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver"
)

func Mongo(err error) Error {

	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		return NoDocuments(err)
	case errors.Is(err, mongo.ErrClientDisconnected):
		return ClientDisconnected(err)
	case errors.Is(err, mongo.ErrNilDocument):
		return BadRequest(err)
	case errors.Is(err, context.DeadlineExceeded):
		return Timeout(err)
	case errors.As(err, &mongo.WriteException{}):
		return HandleWriteException(err)
	case errors.As(err, &mongo.CommandError{}):
		return Command(err)
	case errors.As(err, &mongo.BulkWriteException{}):
		return BulkWrite(err)
	case errors.As(err, &driver.Error{}):
		return Driver(err)
	default:
		return Unknown(err)
	}
}

func NotFound(err error) Error {
	return &Err{
		Status:  http.StatusNotFound,
		Message: "El recurso no existe",
		Err:     err.Error(),
	}
}

func NoDocuments(err error) Error {
	return &Err{
		Status:  http.StatusNotFound,
		Message: "No se encontraron registros",
		Err:     err.Error(),
	}
}

func ClientDisconnected(err error) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Cliente desconectado",
		Err:     err.Error(),
	}
}

func BadRequest(err error) Error {
	return &Err{
		Status:  http.StatusBadRequest,
		Message: "Solicitud incorrecta",
		Err:     err.Error(),
	}
}

func Timeout(err error) Error {
	return &Err{
		Status:  http.StatusRequestTimeout,
		Message: "Tiempo de espera agotado",
		Err:     err.Error(),
	}
}

func HandleWriteException(err error) Error {
	var writeEx mongo.WriteException
	if errors.As(err, &writeEx) {
		for _, we := range writeEx.WriteErrors {
			if we.Code == 11000 {
				return &Err{
					Status:  http.StatusConflict,
					Message: "Registro duplicado",
					Err:     err.Error(),
				}
			}
		}
	}
	return Write(err)
}

func Write(err error) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Error al escribir el registro",
		Err:     err.Error(),
	}
}

func Update(err error) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Error al modificar el registro",
		Err:     err.Error(),
	}
}

func Delete(err error) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Error al eliminar el registro",
		Err:     err.Error(),
	}
}

func Restore(err error) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Error al restaurar el registro",
		Err:     err.Error(),
	}
}

func ForceDelete(err error) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Error al eliminar el registro permanentemente",
		Err:     err.Error(),
	}
}

func Command(err error) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Error de comando",
		Err:     err.Error(),
	}
}

func BulkWrite(err error) Error {
	return &Err{
		Status:  http.StatusBadRequest,
		Message: "Error en escritura masiva",
		Err:     err.Error(),
	}
}

func Driver(err error) Error {
	return &Err{
		Status:  http.StatusBadGateway,
		Message: "Error del driver",
		Err:     err.Error(),
	}
}

func Unknown(err error) Error {
	return &Err{
		Status:  http.StatusInternalServerError,
		Message: "Error inesperado",
		Err:     err.Error(),
	}
}

func HexID(err error) Error {
	return &Err{
		Status:  http.StatusBadRequest,
		Message: "El id no es un hexadecimal válido",
		Err:     err.Error(),
	}
}
