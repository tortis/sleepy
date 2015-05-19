package sleepy

import (
	"net/http"
	"reflect"
	"strings"
)

// Data model for a type that will be produced or consumed by an api call.
type DataModel interface{}

type Param struct {
	Type     string
	Required bool
}

type DataIn struct {
	Body         DataModel
	HeaderParams map[string]Param
	PathParams   map[string]Param
	QueryParams  map[string]Param
}

type DataOut struct {
	Body DataModel
}

func (dm *DataIn) isRequired(field string) bool {
	// Only consider struct fields
	t := reflect.TypeOf(dm.body)
	if t.Kind() != reflect.Struct {
		return false
	}
	if f, e := t.FieldByName(field); e {
		tag := f.StructTag.Get("sleepy")
		props := strings.Split(tag, ",")
		for i := 0; i < len(props); i++ {
			if props[i] == "required" {
				return true
			}
		}
		return false
	} else {
		return false
	}
}
