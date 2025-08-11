package querybuilder

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Match(value ...bson.E) bson.D {
	return bson.D{bson.E{Key: "$match", Value: bson.D(value)}}
}

func Where(key string, value ...bson.E) bson.E {
	return bson.E{Key: key, Value: bson.D(value)}
}

func Filter(value ...bson.E) bson.D {
	return bson.D(value)
}

func Field(key string, value any) bson.E {
	return bson.E{Key: key, Value: value}
}

// operadores logicos ----------------------------------------------------------------

func And(value ...any) bson.E {
	return bson.E{Key: "$and", Value: bson.A(value)}
}

func Or(value ...any) bson.E {
	return bson.E{Key: "$or", Value: bson.A(value)}
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
	return bson.E{Key: "$nor", Value: bson.A(values)}
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
	return bson.E{Key: "$in", Value: bson.A(value)}
}

func Nin(value ...any) bson.E {
	return bson.E{Key: "$nin", Value: bson.A(value)}
}

// ejemplos de prueba ----------------------------------------------------------------
func GetCoffe() mongo.Pipeline {
	return mongo.Pipeline{Match(Where("coffee", Eq("cafe")))}
}

func GetFruits1() bson.D {
	return Match(
		Or(
			Where("fruit", Eq("apple")),
			Where("fruit", Eq("banana")),
			Where("fruit", Eq("orange")),
		),
	)
}

func GetFruits3() bson.D {
	return Match(
		Where("size", Gt(10), Lt(20)),
	)
}

func GetFruits4() bson.D {
	return Match(
		And(
			Where("size", Gt(10), Lt(20)),
			Where("status", Eq("green")),
			Or(
				Where("fruit", Eq("apple")),
				Where("fruit", Eq("banana")),
				Where("fruit", Eq("orange")),
			),
		),
	)
}

func GetFruits5() bson.D {
	return Match(
		Where("size", Gt(10), Lt(20)),
		Where("status", Eq("green")),
		Or(
			Where("fruit", Eq("apple")),
			Where("fruit", Eq("banana")),
			Where("fruit", Eq("orange")),
		),
	)
}

func GetFruits6() bson.D {
	return Match(
		Or(
			Nor(
				Where("fruit", Eq("apple")),
				Where("fruit", Eq("orange")),
				Or(
					Where("size", Gt(30)),
					Where("status", Eq("green")),
				),
			),
			Where("fruit", Eq("watermelon")),
		),
		Where("country", Eq("mexico")),
	)
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

// $lookup → join con otra colección

// $addFields → agrega nuevos campos calculados

// $set → alias de $addFields

// $unset → elimina campos

// $replaceRoot / $replaceWith → reemplaza el documento raíz

// $count → cuenta documentos que pasan el pipeline

// $facet → ejecuta múltiples sub-pipelines en paralelo

// $sortByCount → agrupa y cuenta por valores de un campo

// $merge → escribe resultados en otra colección

// $out → sobrescribe una colección con el resultado

// 2️⃣ Operadores de expresión (dentro de $match, $project, $group, etc.)
// Estos son como "funciones" que puedes usar para comparar, operar con strings, arrays, fechas, etc.

// 🔹 Lógicos
// $and

// $or

// $not

// $nor

// 🔹 Comparación
// $eq

// $ne

// $gt

// $gte

// $lt

// $lte

// $in

// $nin

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
