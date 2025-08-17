package qb

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Lookup(value bson.D) bson.D {
	return bson.D{{Key: "$lookup", Value: value}}
}

func With(collection string, localField string, foreignField string, as string) bson.D {
	return bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: collection},           // colección relacionada
			{Key: "localField", Value: localField},     // campo en colecion local
			{Key: "foreignField", Value: foreignField}, // campo en coleccion relacionada
			{Key: "as", Value: as},                     // campo para el resultado
		}},
	}
}

// HasMany crea un bson.D para realizar una busqueda con relacion hasMany
func HasMany(collection string, foreignField string, with ...bson.D) bson.D {
	if len(with) > 0 {
		return hasManyWith(collection, foreignField, with...)
	}
	return bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: collection},           // colección relacionada
			{Key: "localField", Value: "_id"},          // campo en colecion local
			{Key: "foreignField", Value: foreignField}, // campo en coleccion relacionada
			{Key: "as", Value: collection},             // campo para el resultado
		}},
	}
}

// HasManyWith crea un bson.D para realizar una busqueda con relacion hasMany habilitando relaciones en los documentos hijos
func hasManyWith(collection string, foreignField string, with ...bson.D) bson.D {
	parentID := collection + "_pid"

	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "$expr", Value: bson.D{{Key: "$eq", Value: bson.A{"$" + foreignField, "$$" + parentID}}}},
		}}},
	}

	// pasar el []bson.D → bson.A y no perder el type safety
	for _, w := range with {
		pipeline = append(pipeline, w)
	}

	return bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: collection},
			{Key: "let", Value: bson.D{{Key: parentID, Value: "$_id"}}},
			{Key: "pipeline", Value: pipeline},
			{Key: "as", Value: collection},
		}},
	}
}

// ManyToMany crea un bson.D para realizar una busqueda con relacion manyToMany
func ManyToMany(collection string, localField string, with ...bson.D) bson.D {
	if len(with) > 0 {
		return manyToManyWith(collection, localField, with...)
	}

	return bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: collection},       // colección relacionada
			{Key: "localField", Value: localField}, // campo en colecion local
			{Key: "foreignField", Value: "_id"},    // campo en coleccion relacionada
			{Key: "as", Value: collection},         // campo para el resultado
		}},
	}
}

// ManyToManyWith crea un bson.D para realizar una busqueda con relacion manyToMany habilitando relaciones en los documentos hijos
func manyToManyWith(collection string, localField string, with ...bson.D) bson.D {
	parentID := collection + "_pids"

	// Pipeline base: filtrar solo los documentos cuyo _id esté en el array localField
	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "$expr", Value: bson.D{
				{Key: "$in", Value: bson.A{"$_id", "$$" + parentID}},
			}},
		}}},
	}

	for _, w := range with {
		pipeline = append(pipeline, w)
	}

	return bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: collection},                                      // colección relacionada
			{Key: "let", Value: bson.D{{Key: parentID, Value: "$" + localField}}}, // pasar IDs desde el documento padre
			{Key: "pipeline", Value: pipeline},                                    // pipeline con hijas
			{Key: "as", Value: collection},                                        // alias del resultado
		}},
	}
}

// HasOne crea un bson.D para realizar una busqueda con relacion hasOne
func HasOne(collection string, foreignField string, as string, with ...bson.D) (bson.D, bson.D) {
	if len(with) > 0 {
		return hasOneWith(collection, foreignField, as, with...)
	}

	return bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: collection},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: foreignField},
			{Key: "as", Value: as},
		}}},
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$" + as},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}}
}

// HasOneWith crea un bson.D para realizar una busqueda con relacion hasOne habilitando relaciones en los documentos hijos
func hasOneWith(collection string, foreignField string, as string, with ...bson.D) (bson.D, bson.D) {
	parentID := collection + "_pid"

	// pipeline base para filtrar por relación padre-hijo
	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "$expr", Value: bson.D{{Key: "$eq", Value: bson.A{"$" + foreignField, "$$" + parentID}}}},
		}}},
	}

	for _, w := range with {
		pipeline = append(pipeline, w)
	}

	return bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: collection},
			{Key: "let", Value: bson.D{{Key: parentID, Value: "$_id"}}},
			{Key: "pipeline", Value: pipeline},
			{Key: "as", Value: as},
		}}}, bson.D{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$" + as},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			}}}
}

// BelongsTo crea un bson.D para realizar una busqueda con relacion belongsTo
func BelongsTo(collection string, localField string, as string) []bson.D {
	return []bson.D{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: collection},       // colección relacionada
			{Key: "localField", Value: localField}, // campo en el documento local (foreign key)
			{Key: "foreignField", Value: "_id"},    // campo en la colección relacionada
			{Key: "as", Value: as},                 // nombre del campo resultado
		}}},

		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$" + as},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
	}
}

// BelongsToWith crea un bson.D para realizar una busqueda con relacion belongsTo habilitando relaciones en el documentos padre
func BelongsToWith(collection string, localField string, as string, with ...bson.D) []bson.D {
	foreignID := collection + "_id"

	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "$expr", Value: bson.D{{Key: "$eq", Value: bson.A{"$_id", "$$" + foreignID}}}},
		}}},
	}

	// Agregar lookups anidados si existen
	for _, w := range with {
		pipeline = append(pipeline, w)
	}

	return []bson.D{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: collection},
			{Key: "let", Value: bson.D{{Key: foreignID, Value: "$" + localField}}},
			{Key: "pipeline", Value: pipeline},
			{Key: "as", Value: as},
		}}},

		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$" + as},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
	}
}

// BelongsToMany crea un bson.D para realizar una busqueda con relacion manyToMany inverso
func BelongsToMany(collection string, foreignArrayField string, as string) bson.D {
	return bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: collection},                              // colección relacionada
			{Key: "let", Value: bson.D{{Key: "local_id", Value: "$_id"}}}, // pasamos el _id del documento local
			{Key: "pipeline", Value: mongo.Pipeline{
				{{Key: "$match", Value: bson.D{
					{Key: "$expr", Value: bson.D{
						{Key: "$in", Value: bson.A{"$$local_id", "$" + foreignArrayField}},
					}},
				}}},
			}},
			{Key: "as", Value: as}, // nombre del campo resultado
		}},
	}
}

func BelongsToManyWith(collection string, foreignArrayField string, as string, with ...bson.D) bson.D {
	foreignID := collection + "_id"

	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "$expr", Value: bson.D{
				{Key: "$in", Value: bson.A{"$$" + foreignID, "$" + foreignArrayField}},
			}},
		}}},
	}

	// Añadir relaciones hijas si se pasan
	for _, w := range with {
		pipeline = append(pipeline, w)
	}

	return bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: collection},
			{Key: "let", Value: bson.D{{Key: foreignID, Value: "$_id"}}},
			{Key: "pipeline", Value: pipeline},
			{Key: "as", Value: as},
		}},
	}
}
