package qb

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/exp/constraints"
)

func Match(value ...bson.E) bson.D {
	return bson.D{bson.E{Key: "$match", Value: bson.D(value)}}
}

func Where(key string, value ...bson.E) bson.E {
	return bson.E{Key: key, Value: bson.D(value)}
}

func OrWhere(key string, value ...bson.E) bson.D {
	return bson.D{bson.E{Key: key, Value: bson.D(value)}}
}

func Document(value ...bson.E) bson.D {
	return bson.D(value)
}

func Filter(value ...bson.E) bson.D {
	return bson.D(value)
}

func Element(key string, value any) bson.E {
	return bson.E{Key: key, Value: value}
}

func Array(value ...any) bson.A {
	return value
}

func Pipeline(value ...bson.D) mongo.Pipeline {
	return value
}

// operadores logicos ----------------------------------------------------------------

func And(value ...any) bson.E {
	return bson.E{Key: "$and", Value: value}
}

func Or(value ...any) bson.E {
	return bson.E{Key: "$or", Value: value}
}

func Not(key string, value any) bson.E {
	return bson.E{
		Key: key,
		Value: bson.D{
			{Key: "$not", Value: value},
		},
	}
}

func Nor(values ...any) bson.E {
	return bson.E{Key: "$nor", Value: values}
}

// operadores de comparacion ----------------------------------------------------------------

func Eq(value any) bson.E {
	return bson.E{Key: "$eq", Value: value}
}

func Ne(value any) bson.E {
	return bson.E{Key: "$ne", Value: value}
}

func Gt(value any) bson.E {
	return bson.E{Key: "$gt", Value: value}
}

func Gte(value any) bson.E {
	return bson.E{Key: "$gte", Value: value}
}

func Lt(value any) bson.E {
	return bson.E{Key: "$lt", Value: value}
}

func Lte(value any) bson.E {
	return bson.E{Key: "$lte", Value: value}
}

func In(value ...any) bson.E {
	return bson.E{Key: "$in", Value: value}
}

func Nin(value ...any) bson.E {
	return bson.E{Key: "$nin", Value: value}
}

// operadores de elemento ----------------------------------------------------------------

func Exists() bson.E {
	return bson.E{Key: "$exists", Value: true}
}

func NotExists() bson.E {
	return bson.E{Key: "$exists", Value: false}
}

func Type(value string) bson.E {
	return bson.E{Key: "$type", Value: value}
}

// funciones de softdelete ----------------------------------------------------------------

func OnlyTrashed() bson.E {
	return bson.E{Key: "deleted_at", Value: bson.D{bson.E{Key: "$exists", Value: true}}}
}

func WithOutTrashed() bson.E {
	return bson.E{Key: "deleted_at", Value: bson.D{bson.E{Key: "$exists", Value: false}}}
}

// operadores de evacuacion ----------------------------------------------------------------

func Expr(value ...bson.E) bson.E {
	return bson.E{Key: "$expr", Value: bson.D(value)}
}

func JsonSchema(value any) bson.E {
	return bson.E{Key: "$jsonSchema", Value: value}
}

func Mod[D constraints.Integer | constraints.Float, R constraints.Integer | constraints.Float](divisor D, remainder R) bson.E {
	return bson.E{Key: "$mod", Value: bson.A([]any{divisor, remainder})}
}

func Regex(value any) bson.E {
	return bson.E{Key: "$regex", Value: value} // esta de aca queda pendiente para cuando investigue mejor
}

// operadores geoespaciales ----------------------------------------------------------------

//quedan pendientes por que estan complejos y no la quiero cagar
// ademas no los uso aun.

// operadores de array ----------------------------------------------------------------

func All(value ...any) bson.E {
	return bson.E{Key: "$all", Value: bson.A(value)}
}

func ElemMatch(value ...bson.E) bson.E {
	return bson.E{Key: "$elemMatch", Value: bson.D(value)}
}

func Size(value int) bson.E {
	return bson.E{Key: "$size", Value: value}
}

// operadores bitwise ----------------------------------------------------------------

// quedan pendientes

// operadores de proyeccion ----------------------------------------------------------------

func Slice(value int) bson.E {
	return bson.E{Key: "$slice", Value: value}
}

// operadores de agregacion ----------------------------------------------------------------

func Unset(keys ...string) bson.D {
	elements := bson.D{}
	for _, key := range keys {
		elements = append(elements, bson.E{Key: key, Value: nil})
	}
	return bson.D{bson.E{Key: "$unset", Value: elements}}
}

func Set(value ...bson.E) bson.D {
	return bson.D{bson.E{Key: "$set", Value: value}}
}

// operadores miselanios ----------------------------------------------------------------

func Rand() bson.E {
	return bson.E{Key: "$rand", Value: bson.D{}}
}

// 1️⃣ Operadores de etapa (Aggregation Pipeline Stages)
// Son los que usas en el pipeline (mongo.Pipeline{} en tu caso). Ejemplos:

// $match → filtra documentos (similar a find pero en el pipeline)

// $project → selecciona / transforma campos

// $group → agrupa documentos y calcula valores agregados

// $sort → ordena documentos

// $limit → limita el número de resultados

// $skip → salta documentos

// $unwind → desestructura arrays en documentos separados

// $addFields → agrega nuevos campos calculados en pipeline

// $replaceRoot / $replaceWith → reemplaza el documento raíz

// $count → cuenta documentos que pasan el pipeline

// $facet → ejecuta múltiples sub-pipelines en paralelo

// $sortByCount → agrupa y cuenta por valores de un campo

// $merge → escribe resultados en otra colección

// $out → sobrescribe una colección con el resultado

// 2️⃣ Operadores de expresión (dentro de $match, $project, $group, etc.)
// Estos son como "funciones" que puedes usar para comparar, operar con strings, arrays, fechas, etc.

// 🔹 Operadores de array
// $size

// $all

// $elemMatch

// $slice

// 🔹 Strings
// $concat

// $substr (deprecado, usar $substrBytes o $substrCP)

// $toLower / $toUpper

// $split

// $regexMatch

// 🔹 Fechas
// $year

// $month

// $dayOfMonth

// $dayOfWeek

// $hour

// $minute

// $second

// $dateFromString

// $dateToString

// 🔹 Condicionales
// $cond

// $ifNull

// $switch

// 3️⃣ Operadores de acumulador (solo dentro de $group, $setWindowFields, etc.)
// Sirven para calcular valores agregados:

// $sum

// $avg

// $min

// $max

// $first

// $last

// $push (mete valores en un array)

// $addToSet (mete valores únicos en un array)

// $count (en stages específicos)

// $stdDevPop / $stdDevSamp
