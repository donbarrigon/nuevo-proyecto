package qb

import (
	"strconv"
	"strings"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// FindOptions toma los parámetros de búsqueda del contexto (ctx) y los convierte en un conjunto de opciones de búsqueda (FindOptions) de MongoDB.
// @param page: número de página para la paginación (predeterminado: 0)
// @param per_page: número de elementos por página (predeterminado: 15, máximo: 1000)
// @param sort: orden de los resultados.
//
//	Formato: sort=campo1,-campo2,+campo3 o sort=campo1&sort=-campo2&sort=+campo3
//	En este caso, "+" indica orden ascendente y "-" orden descendente si no se especifica por defecto es ascendente.
//	Ejemplo: sort=campo1,-campo2 ordena primero por 'campo1' ascendente y luego por 'campo2' descendente.
//
// @param projection: campos a incluir o excluir en los resultados.
//
//	Predeterminado: incluye solo el campo '_id'. Ejemplo de uso:
//	projection=campo1,campo2,campo3 (solo incluye estos campos en el resultado)
//	projection=campo1&projection=campo2&projection=campo3 (incluye estos tres campos en el resultado)
//	projection=id excluye '_id' (el campo por defecto siempre se incluye a menos que se indique lo contrario).
func FindOptions(ctx *app.HttpContext) *options.FindOptionsBuilder {
	findOptions := options.Find()
	urlValues := ctx.Request.URL.Query()

	limit, er := strconv.ParseInt(urlValues.Get("per_page"), 10, 64)
	if er != nil {
		limit = 15
	}
	if limit > 1000 {
		limit = 1000
	}
	page, er := strconv.ParseInt(urlValues.Get("page"), 10, 64)
	if er != nil {
		page = 1
	}
	findOptions.SetSkip((page - 1) * limit)
	findOptions.SetLimit(limit)

	sort := bson.D{}
	for _, field := range urlValues["sort"] {
		fields := strings.Split(field, ",")
		for _, f := range fields {
			if strings.HasPrefix(field, "-") {
				sort = append(sort, bson.E{Key: strings.TrimPrefix(f, "-"), Value: -1})
			} else {
				sort = append(sort, bson.E{Key: strings.TrimPrefix(f, "+"), Value: 1})
			}
		}
	}
	findOptions.SetSort(sort)

	projection := bson.D{}
	for _, field := range urlValues["projection"] {
		fields := strings.Split(field, ",")
		for _, f := range fields {
			if f == "id" {
				projection = append(projection, bson.E{Key: "_id", Value: 0})
			} else {
				projection = append(projection, bson.E{Key: f, Value: 1})
			}
		}

	}
	findOptions.SetProjection(projection)

	return findOptions
}
