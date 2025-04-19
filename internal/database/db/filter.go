package db

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const WITHOUT_TRASH = 0
const WITH_TRASH = 1
const ONLY_TRASH = 2

var MongoFilterMap = map[string]string{
	"eq":  "$eq",  // Igual a
	"=":   "$eq",  // Igual a
	"ne":  "$ne",  // Distinto de
	"!=":  "$ne",  // Distinto de
	"gt":  "$gt",  // Mayor que
	">":   "$gt",  // Mayor que
	"gte": "$gte", // Mayor o igual que
	">=":  "$gte", // Mayor o igual que
	"lt":  "$lt",  // Menor que
	"<":   "$lt",  // Menor que
	"lte": "$lte", // Menor o igual que
	"<=":  "$lte", // Menor o igual que

	"lk":    "$regex", // LIKE (puede usarse con opciones de case-insensitive)
	"like":  "$regex", // LIKE
	"ilk":   "$regex", // ILIKE (insensitive)
	"ilike": "$regex", // ILIKE

	"in":  "$in",  // Dentro de una lista
	"nin": "$nin", // No dentro de una lista

	"null":    "$eq", // Es nulo (con valor `nil`)
	"nnull":   "$ne", // No es nulo (con valor `nil`)
	"notnull": "$ne", // No es nulo (alias de `nnull`)

	"between": "$gte_lte", // Entre dos valores (requiere lógica especial)
	"bw":      "$gte_lte", // Alias de between
}

type Filter struct {
	Key    string
	Filter string
	Value  string
}

type Sort struct {
	Field     string
	Direction int // 1 para asc, -1 para desc (Mongo style)
}

type QueryFilter struct {
	Filters         []Filter
	Sort            []Sort
	GroupBy         []string
	Page            int // si es cero la paginacion es por cursor
	PerPage         int
	Cursor          string
	CursorDirection int    // 1 para asc, -1 para desc (Mongo style)
	Trash           int    // 0 para without(default) 1 para with y 2 only
	Path            string // ruta del path para la paginacion
}

func NewQueryFilter() *QueryFilter {
	return &QueryFilter{
		Filters:         []Filter{},
		Sort:            []Sort{},
		GroupBy:         []string{},
		Page:            0,
		PerPage:         15,
		Cursor:          "",
		CursorDirection: 1,
		Trash:           0,
	}
}

func (qf *QueryFilter) AppendFilter(key, filter, value string) {
	qf.Filters = append(qf.Filters, Filter{key, filter, value})
}

func (qf *QueryFilter) AppendSort(field string, direction int) {
	qf.Sort = append(qf.Sort, Sort{field, direction})
}

func (qf *QueryFilter) AppendGrouBy(field string) {
	qf.GroupBy = append(qf.GroupBy, field)
}

func (qf *QueryFilter) WithoutTrash() {
	qf.Trash = 0
}

func (qf *QueryFilter) WithTrash() {
	qf.Trash = 1
}

func (qf *QueryFilter) OnlyTrash() {
	qf.Trash = 2
}

func (qf *QueryFilter) Equals(key, value string) {
	qf.Filters = append(qf.Filters, Filter{Key: key, Filter: "eq", Value: value})
}

func (qf *QueryFilter) NotEquals(key, value string) {
	qf.Filters = append(qf.Filters, Filter{Key: key, Filter: "ne", Value: value})
}

func (qf *QueryFilter) GreaterThan(key, value string) {
	qf.Filters = append(qf.Filters, Filter{Key: key, Filter: "gt", Value: value})
}

func (qf *QueryFilter) GreaterOrEqual(key, value string) {
	qf.Filters = append(qf.Filters, Filter{Key: key, Filter: "gte", Value: value})
}

func (qf *QueryFilter) LessThan(key, value string) {
	qf.Filters = append(qf.Filters, Filter{Key: key, Filter: "lt", Value: value})
}

func (qf *QueryFilter) LessOrEqual(key, value string) {
	qf.Filters = append(qf.Filters, Filter{Key: key, Filter: "lte", Value: value})
}

func (qf *QueryFilter) In(key string, values ...string) {
	qf.Filters = append(qf.Filters, Filter{
		Key:    key,
		Filter: "in",
		Value:  strings.Join(values, ","),
	})
}

func (qf *QueryFilter) NotIn(key string, values ...string) {
	qf.Filters = append(qf.Filters, Filter{
		Key:    key,
		Filter: "nin",
		Value:  strings.Join(values, ","),
	})
}

func (qf *QueryFilter) Between(key, from, to string) {
	qf.Filters = append(qf.Filters, Filter{
		Key:    key,
		Filter: "between",
		Value:  fmt.Sprintf("%s,%s", from, to),
	})
}

func (qf *QueryFilter) IsNull(key string) {
	qf.Filters = append(qf.Filters, Filter{Key: key, Filter: "null", Value: ""})
}

func (qf *QueryFilter) IsNotNull(key string) {
	qf.Filters = append(qf.Filters, Filter{Key: key, Filter: "notnull", Value: ""})
}

func (qf *QueryFilter) Like(key, pattern string) {
	qf.Filters = append(qf.Filters, Filter{Key: key, Filter: "like", Value: pattern})
}

func (qf *QueryFilter) ILike(key, pattern string) {
	qf.Filters = append(qf.Filters, Filter{Key: key, Filter: "ilike", Value: pattern})
}

func (qf *QueryFilter) CursorPaginate() {
	qf.Page = 0
}

func (qf *QueryFilter) All() {
	qf.Page = 0
	qf.PerPage = 0
	qf.Cursor = ""
	qf.CursorDirection = 0
}

func (qf *QueryFilter) Pipeline() mongo.Pipeline {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: qf.BsonD()}},
	}

	if len(qf.GroupBy) > 0 {
		groupID := bson.D{}
		for _, field := range qf.GroupBy {
			groupID = append(groupID, bson.E{Key: field, Value: "$" + field})
		}
		pipeline = append(pipeline, bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: groupID},
				{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
			}},
		})
	}

	if len(qf.Sort) > 0 {
		sortDoc := bson.D{}
		for _, s := range qf.Sort {
			sortDoc = append(sortDoc, bson.E{Key: s.Field, Value: s.Direction})
		}
		pipeline = append(pipeline, bson.D{{Key: "$sort", Value: sortDoc}})
	}

	if qf.Page > 0 {
		skip := (qf.Page - 1) * qf.PerPage
		pipeline = append(pipeline, bson.D{{Key: "$skip", Value: skip}})
	}

	if qf.PerPage > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$limit", Value: qf.PerPage}})
	}

	return pipeline
}

func (qf *QueryFilter) BsonD() *bson.D {
	filters := bson.D{}

	switch qf.Trash {
	case 0: // sin basura
		filters = append(filters, bson.E{Key: "deleted_at", Value: nil})
	case 2: // solo basura
		filters = append(filters, bson.E{Key: "deleted_at", Value: bson.M{"$ne": nil}})
	}

	if qf.Page == 0 && qf.Cursor != "" {
		if len(qf.Sort) > 0 {
			sortField := qf.Sort[0].Field
			cursorOp := "$gt"
			if qf.CursorDirection == -1 {
				cursorOp = "$lt"
			}
			filters = append(filters, bson.E{Key: sortField, Value: bson.M{cursorOp: qf.Cursor}})
		} else {
			cursorOp := "$gt"
			if qf.CursorDirection == -1 {
				cursorOp = "$lt"
			}
			filters = append(filters, bson.E{Key: "_id", Value: bson.M{cursorOp: qf.Cursor}})
		}
	}

	for _, f := range qf.Filters {
		mongoOp, ok := MongoFilterMap[f.Filter]
		if !ok {
			continue // Operador no soportado
		}

		switch f.Filter {
		case "in", "nin":
			// Split por comas para múltiples valores
			values := strings.Split(f.Value, ",")
			arr := make(bson.A, 0, len(values))
			for _, val := range values {
				arr = append(arr, strings.TrimSpace(val))
			}
			filters = append(filters, bson.E{Key: f.Key, Value: bson.M{mongoOp: arr}})

		case "between", "bw":
			values := strings.Split(f.Value, ",")
			if len(values) != 2 {
				continue // Ignorar si no hay 2 valores exactos
			}
			gte := strings.TrimSpace(values[0])
			lte := strings.TrimSpace(values[1])
			filters = append(filters, bson.E{Key: f.Key, Value: bson.M{"$gte": gte, "$lte": lte}})

		case "null":
			filters = append(filters, bson.E{Key: f.Key, Value: nil})

		case "nnull", "notnull":
			filters = append(filters, bson.E{Key: f.Key, Value: bson.M{"$ne": nil}})

		case "lk", "like", "ilk", "ilike":
			regex := f.Value
			options := ""
			if f.Filter == "ilk" || f.Filter == "ilike" {
				options = "i"
			}
			filters = append(filters, bson.E{Key: f.Key, Value: bson.M{"$regex": regex, "$options": options}})

		default:
			filters = append(filters, bson.E{Key: f.Key, Value: bson.M{mongoOp: f.Value}})
		}
	}

	return &filters
}
