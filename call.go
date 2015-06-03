package sleepy

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type Call struct {
	path          string
	method        string
	operationName string
	handler       Handler
	filters       []Filter
	model         callDataModel
}

// Implement the Handler interface.
func (c *Call) ServeHTTP(w http.ResponseWriter, r *http.Request, d map[string]interface{}) {
	// Parse url/path variables and store them in the *sleepy.Request

	//Parse the request body into the reads model if applicable
	if c.model.bodyIn != nil && r.Method != "GET" {
		dec := json.NewDecoder(r.Body)
		payload := reflect.New(reflect.TypeOf(c.model.bodyIn)).Interface()
		err := dec.Decode(payload)
		if err != nil {
			apiErr := newRequestError("Could not parse the request", err)
			http.Error(w, apiErr.Message(), apiErr.StatusCode())
			logResult(r, apiErr)
			return
		}
		apiErr := c.model.validateTagsIn(payload)
		if apiErr != nil {
			http.Error(w, apiErr.Message(), apiErr.StatusCode())
			logResult(r, apiErr)
			return
		}
		d["body"] = payload
	}

	// and validate that required fields are present

	// Call filters
	for _, filter := range c.filters {
		err := filter(w, r, d)
		if err != nil {
			http.Error(w, err.Message(), err.StatusCode())
			logResult(r, err)
			return
		}
	}

	// Call handler
	result, apiErr := c.handler(w, r, d)
	if apiErr != nil {
		http.Error(w, apiErr.Message(), apiErr.StatusCode())
		logResult(r, apiErr)
		return
	}

	if result == nil {
		apiErr = newInternalError(errors.New("Call handler for " + c.operationName + " did not return a response or an error."))
		http.Error(w, apiErr.Message(), apiErr.StatusCode())
		logResult(r, apiErr)
		return
	}

	jb, err := json.Marshal(result)
	if err != nil {
		apiErr = newInternalError(errors.New("Response from call handler for " + c.operationName + " could not be parsed to JSON."))
		http.Error(w, apiErr.Message(), apiErr.StatusCode())
		logResult(r, apiErr)
		return
	}
	w.Header().Set("Content-Type", "Application/JSON")
	w.Write(jb)
	logResult(r, nil)
}

type callDataModel struct {
	bodyIn    interface{}
	bodyOut   interface{}
	pathVars  []inputVar
	queryVars []inputVar
}

func (model *callDataModel) validateTagsIn(payload interface{}) Error {
	pValue := reflect.ValueOf(payload).Elem()
	modType := reflect.TypeOf(model.bodyIn)
	for i := 0; i < modType.NumField(); i++ {
		isRequired := false
		isReadOnly := false
		tags := strings.Split(modType.Field(i).Tag.Get("sleepy"), ",")
		for _, tag := range tags {
			if tag == "required" {
				isRequired = true
				log.Println(modType.Field(i).Name + " is required")
			} else if tag == "readonly" {
				isReadOnly = true
			}
		}

		if isRequired {
			if IsZeroOfUnderlyingType(pValue.Field(i).Interface()) {
				return newRequestError("Required field: "+modType.Field(i).Name+" is missing.", errors.New("Failed while validating tags for the payload."))
			}
		}

		if isReadOnly {
			if !IsZeroOfUnderlyingType(pValue.Field(i).Interface()) {
				return newRequestError("Attempting to set read-only field: "+modType.Field(i).Name+".", errors.New("Failed while validating tags for the payload."))
			}
		}

	}
	return nil
}

func IsZeroOfUnderlyingType(x interface{}) bool {
	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}

type inputVar struct {
	typ      string
	name     string
	desc     string
	required bool
}

// Call builder method
func (c *Call) Method(method string) *Call {
	if c.model.bodyIn != nil {
		log.Fatal("Cannot set method of call to GET since Reads() was set. GET calls do not have a body.")
	}
	c.method = method
	return c
}

// Call builder method
func (c *Call) To(fn Handler) *Call {
	c.handler = fn
	return c
}

// Call builder method
func (c *Call) Filter(f Filter) *Call {
	c.filters = append(c.filters, f)
	return c
}

// Call builder method
func (c *Call) PathParam(name, desc string) *Call {
	c.model.pathVars = append(c.model.pathVars, inputVar{
		typ:      "string",
		name:     name,
		desc:     desc,
		required: true})
	return c
}

// Call builder method
func (c *Call) Reads(m interface{}) *Call {
	// No body allowed in a GET call
	if c.method == "GET" {
		log.Fatal("A GET call can not use Reads() because the GET method does not have a request body.")
	}

	// Confirm the model is a struct
	if reflect.TypeOf(m).Kind() != reflect.Struct {
		log.Fatal("The model given to Reads() must be of kind 'Struct'")
	}

	c.model.bodyIn = m
	return c
}

// Call builder method
func (c *Call) Returns(m interface{}) *Call {
	c.model.bodyOut = m
	return c
}

// Call builder method
func (c *Call) OperationName(name string) *Call {
	c.operationName = name
	return c
}
