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

// A middleware function that may do some work before a handler.
// Filters can store arbitrary data inside the Request.Items map.
type Filter func(Response, *Request) (Response, *Request)

// End function that will serve the API call
type APIHandler func(http.ResponseWriter, *http.Request)

// sleepy response wraps http.ResponseWriter and provides helper functions
type Response struct {
	Res http.ResponseWriter
}

// sleep Request wraps http.Requst and holds an Items map that Filters may
// use to store data for the handler to use later
type Request struct {
	Req   *http.Request
	Items map[string]interface{}
}

// Documentation that is generated for an API call
type CallDoc struct {
	Path         string
	Summary      string
	Description  string
	OperatonName string
	Params       DataIn
	Returns      DataModel
}

// A modular API call. Contains a handler function and filter functions
// that will do the heavy lifting
type APICall struct {
	path    string
	handler APIHandler
	filters []Filter
	doc     CallDoc
}
