package sleepy

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"reflect"
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
			endCall(w, r, apiErr)
			return
		}
		apiErr := c.model.validateTagsIn(payload)
		if apiErr != nil {
			endCall(w, r, apiErr)
			return
		}
		d["body"] = payload
	}

	// and validate that required fields are present

	// Call filters
	for _, filter := range c.filters {
		err := filter(w, r, d)
		if err != nil {
			endCall(w, r, err)
			return
		}
	}

	// Call handler
	result, apiErr := c.handler(w, r, d)
	if apiErr != nil {
		endCall(w, r, apiErr)
		return
	}

	if result == nil {
		apiErr = newInternalError(errors.New("Call handler for " + c.operationName + " did not return a response or an error."))
		endCall(w, r, apiErr)
		return
	}

	jb, err := json.Marshal(result)
	if err != nil {
		apiErr = newInternalError(errors.New("Response from call handler for " + c.operationName + " could not be parsed to JSON."))
		endCall(w, r, apiErr)
		return
	}
	w.Header().Set("Content-Type", "Application/JSON")
	w.Write(jb)
	endCall(w, r, nil)
}

////////////////////////////////////////////////////////////////////////////////
//
////////////////////////////////////////////////////////////////////////////////
func (c *Call) Method(method string) *Call {
	if c.model.bodyIn != nil {
		log.Fatal("Cannot set method of call to GET since Reads() was set. GET calls do not have a body.")
	}
	c.method = method
	return c
}

////////////////////////////////////////////////////////////////////////////////
//
////////////////////////////////////////////////////////////////////////////////
func (c *Call) To(fn Handler) *Call {
	c.handler = fn
	return c
}

////////////////////////////////////////////////////////////////////////////////
//
////////////////////////////////////////////////////////////////////////////////
func (c *Call) Filter(f Filter) *Call {
	c.filters = append(c.filters, f)
	return c
}

////////////////////////////////////////////////////////////////////////////////
//
////////////////////////////////////////////////////////////////////////////////
func (c *Call) PathParam(name, desc string) *Call {
	c.model.pathVarsDoc = append(c.model.pathVarsDoc, inputVar{
		typ:      "string",
		name:     name,
		desc:     desc,
		required: true})
	return c
}

////////////////////////////////////////////////////////////////////////////////
//
////////////////////////////////////////////////////////////////////////////////
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
	c.model.identifyFieldTags(nil)
	return c
}

////////////////////////////////////////////////////////////////////////////////
//
////////////////////////////////////////////////////////////////////////////////
func (c *Call) Returns(m interface{}) *Call {
	c.model.bodyOut = m
	return c
}

////////////////////////////////////////////////////////////////////////////////
//
////////////////////////////////////////////////////////////////////////////////
func (c *Call) OperationName(name string) *Call {
	c.operationName = name
	return c
}
