package db

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// HasMany crea un bson.D para realizar una busqueda con relacion hasMany
func HasMany(collection string, foreignField string) bson.D {
	return bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: collection},           // colección relacionada
			{Key: "localField", Value: "_id"},          // campo en colecion local
			{Key: "foreignField", Value: foreignField}, // campo en coleccion relacionada
			{Key: "as", Value: collection},             // campo para el resultado
		}},
	}
}

// HasManyf crea un bson.D para realizar una busqueda con relacion hasMany para aquellos que no siguien las conveciones
func HasManyf(collection string, localField string, foreignField string, as string) bson.D {
	return bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: collection},           // colección relacionada
			{Key: "localField", Value: localField},     // campo en colecion local
			{Key: "foreignField", Value: foreignField}, // campo en coleccion relacionada
			{Key: "as", Value: as},                     // campo para el resultado
		}},
	}
}

// ManyToMany crea un bson.D para realizar una busqueda con relacion manyToMany
func ManyToMany(collection string, localField string) bson.D {
	return bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: collection},       // colección relacionada
			{Key: "localField", Value: localField}, // campo en colecion local
			{Key: "foreignField", Value: "_id"},    // campo en coleccion relacionada
			{Key: "as", Value: collection},         // campo para el resultado
		}},
	}
}

// ManyToManyf crea un bson.D para realizar una busqueda con relacion manyToMany para aquellos que no siguien las conveciones
func ManyToManyf(collection string, localField string, foreignField string, as string) bson.D {
	return bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: collection},           // colección relacionada
			{Key: "localField", Value: localField},     // campo en colecion local
			{Key: "foreignField", Value: foreignField}, // campo en coleccion relacionada
			{Key: "as", Value: as},                     // campo para el resultado
		}},
	}
}

// HasOne crea un bson.D para realizar una busqueda con relacion hasOne
func HasOne(collection string, foreignField string, as string) []bson.D {
	return []bson.D{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: collection},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: foreignField},
			{Key: "as", Value: as},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$" + as},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
	}
}

// HasOnef crea un bson.D para realizar una busqueda con relacion hasOne para aquellos que no siguien las conveciones
func HasOnef(collection string, localField string, foreignField string, as string) []bson.D {
	return []bson.D{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: collection},
			{Key: "localField", Value: localField},
			{Key: "foreignField", Value: foreignField},
			{Key: "as", Value: as},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$" + as},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
	}
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

// BelongsTof crea un bson.D para realizar una busqueda con relacion belongsTo para aquellos que no siguien las conveciones
func BelongsTof(collection string, localField string, foreignField string, as string) []bson.D {
	return []bson.D{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: collection},           // colección relacionada
			{Key: "localField", Value: localField},     // campo en el documento local (foreign key)
			{Key: "foreignField", Value: foreignField}, // campo en la colección relacionada
			{Key: "as", Value: as},                     // nombre del campo resultado
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
			{Key: "from", Value: collection}, // colección relacionada
			{Key: "let", Value: bson.D{
				{Key: "local_id", Value: "$_id"}, // pasamos el _id del documento local
			}},
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

// BelongsToManyf crea un bson.D para realizar una busqueda con relacion manyToMany inverso para aquellos que no siguien las conveciones
func BelongsToManyf(collection string, localID string, foreignArrayField string, as string) bson.D {
	return bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: collection}, // colección relacionada
			{Key: "let", Value: bson.D{
				{Key: "local_id", Value: localID}, // pasamos el _id del documento local
			}},
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

func HasMany2(collection string, foreignField string) bson.D {
	return bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: collection},           // colección relacionada
			{Key: "localField", Value: "_id"},          // campo en colecion local
			{Key: "foreignField", Value: foreignField}, // campo en coleccion relacionada
			{Key: "as", Value: collection},             // campo para el resultado
		}},
	}
}

func HasManyNested(
	parentCollection string, // roles
	parentForeignField string, // _id del user está en role_ids
	childCollection string, // permissions
	childLocalField string, // permission_ids
	childForeignField string, // _id en permissions
	as string, // nombre para roles
	childAs string, // nombre para permissions
) bson.D {
	return bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: parentCollection},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: parentForeignField},
			{Key: "pipeline", Value: mongo.Pipeline{
				{{Key: "$lookup", Value: bson.D{
					{Key: "from", Value: childCollection},
					{Key: "localField", Value: childLocalField},
					{Key: "foreignField", Value: childForeignField},
					{Key: "as", Value: childAs},
				}}},
			}},
			{Key: "as", Value: as},
		}},
	}
}
