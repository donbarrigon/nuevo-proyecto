package db

import (
	"context"
	"fmt"
	"reflect"

	"github.com/donbarrigon/nuevo-proyecto/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type PaginateResource struct {
	Data  any            `json:"data"`
	Links any            `json:"links"`
	Meta  map[string]any `json:"meta"`
}

func Paginate(model MongoModel, result any, qf *QueryFilter) (*PaginateResource, errors.Error) {
	cursor, err := Mongo.Database.Collection(model.CollectionName()).Aggregate(context.TODO(), qf.Pipeline())
	if err != nil {
		return nil, errors.Mongo(err)
	}
	if err = cursor.All(context.TODO(), result); err != nil {
		return nil, errors.Mongo(err)
	}

	paginated := &PaginateResource{
		Data: result,
	}

	meta := make(map[string]any)
	meta["path"] = qf.Path
	meta["per_page"] = qf.PerPage

	val := reflect.ValueOf(result)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	result_length := 0
	if val.Kind() == reflect.Slice {
		result_length = val.Len()
	}

	if qf.Page > 0 {
		// Paginación por offset
		meta["current_page"] = qf.Page
		paginated.Links = make([]map[string]string, 0)

		// Clonamos el filtro sin paginación para contar
		countQF := *qf
		countQF.Page = 0
		countQF.PerPage = 0
		countQF.Cursor = ""
		countQF.CursorDirection = 0

		total, err := Mongo.Database.Collection(model.CollectionName()).CountDocuments(context.TODO(), countQF.Pipeline())
		if err != nil {
			return nil, errors.Mongo(err)
		}

		lastPage := int((int64(total) + int64(qf.PerPage) - 1) / int64(qf.PerPage))
		from := (qf.Page-1)*qf.PerPage + 1
		to := from + result_length - 1

		meta["total"] = total
		meta["last_page"] = lastPage
		meta["from"] = from
		meta["to"] = to

		links := make([]map[string]string, 0)
		for p := 1; p <= lastPage; p++ {
			link := map[string]string{
				"url":    fmt.Sprintf("%s?page=%d&per_page=%d", qf.Path, p, qf.PerPage),
				"label":  fmt.Sprintf("%d", p),
				"active": fmt.Sprintf("%v", p == qf.Page),
			}
			links = append(links, link)
		}

		if qf.Page > 1 {
			prev := map[string]string{
				"url":    fmt.Sprintf("%s?page=%d&per_page=%d", qf.Path, qf.Page-1, qf.PerPage),
				"label":  "« Anterior",
				"active": "false",
			}
			links = append([]map[string]string{prev}, links...)
		}

		if qf.Page < lastPage {
			next := map[string]string{
				"url":    fmt.Sprintf("%s?page=%d&per_page=%d", qf.Path, qf.Page+1, qf.PerPage),
				"label":  "Siguiente »",
				"active": "false",
			}
			links = append(links, next)
		}

		paginated.Links = links
		paginated.Meta = meta

	} else {
		// Paginación por cursor

		sortField := "id"
		if len(qf.Sort) > 0 {
			sortField = qf.Sort[0].Field
		}

		if result_length > 0 {
			prev := val.Index(0).Interface()
			next := val.Index(result_length - 1).Interface()

			prevVal := getFieldValueByJSONTag(prev, sortField)
			if objID, ok := prevVal.(bson.ObjectID); ok {
				prevVal = objID.Hex()
			}
			nextVal := getFieldValueByJSONTag(next, sortField)
			if objID, ok := nextVal.(bson.ObjectID); ok {
				nextVal = objID.Hex()
			}

			if prevVal != nil {
				meta["prev_cursor"] = fmt.Sprintf("%v", prevVal)
			}
			if nextVal != nil {
				meta["next_cursor"] = fmt.Sprintf("%v", nextVal)
			}
		}

		links := make(map[string]string)
		if meta["next_cursor"] != nil {
			links["next"] = fmt.Sprintf("%s?cursor[asc]=%s", qf.Path, meta["next_cursor"])
		}
		if meta["prev_cursor"] != nil {
			links["prev"] = fmt.Sprintf("%s?cursor[desc]=%s", qf.Path, meta["prev_cursor"])
		}

		paginated.Links = links
		paginated.Meta = meta
	}

	return paginated, nil
}
func getFieldValueByJSONTag(obj any, tag string) any {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == tag {
			return v.Field(i).Interface()
		}
	}
	return nil
}
