package sleepy

import (
	"encoding/json"
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

	// Validate that required queryVars are present
	apiErr := c.model.validateQueryVars(r, d)
	if apiErr != nil {
		endCall(w, r, apiErr, d)
		return
	}

	//Parse the request body into the reads model if applicable
	if c.model.bodyIn.model != nil && r.Method != "GET" {
		dec := json.NewDecoder(r.Body)
		payload := reflect.New(reflect.TypeOf(c.model.bodyIn.model)).Interface()
		err := dec.Decode(payload)
		if err != nil {
			apiErr := ErrBadRequest(err.Error(), "Could not parse the request.", ERR_PARSE_REQUEST)
			endCall(w, r, apiErr, d)
			return
		}
		apiErr := c.model.validateTagsIn(payload)
		if apiErr != nil {
			endCall(w, r, apiErr, d)
			return
		}
		d["body"] = payload
	}

	// Call filters
	for _, filter := range c.filters {
		err := filter(r, d)
		if err != nil {
			endCall(w, r, err, d)
			return
		}
	}

	// Call handler
	result, apiErr := c.handler(w, r, d)
	if apiErr != nil {
		endCall(w, r, apiErr, d)
		return
	}

	if result == nil {
		apiErr = ErrInternal("Call handler for " + c.operationName + " did not return a response or an error.")
		endCall(w, r, apiErr, d)
		return
	}

	// Remove any fields that are write only
	// Get the value of the result
	if c.model.bodyOut.model != nil {
		rawVal := reflect.ValueOf(result)
		var val reflect.Value
		if rawVal.Kind() == reflect.Ptr {
			val = rawVal.Elem()
		} else {
			val = rawVal
		}
		// Then zero each field in the result that is marked as writeonly
		for _, fieldIndex := range c.model.bodyOut.woFields {
			field := val.FieldByIndex(fieldIndex)
			field.Set(reflect.Zero(field.Type()))
		}
	}

	// Marshal the result into json and write the response
	jb, err := json.Marshal(result)
	if err != nil {
		apiErr = ErrInternal("Response from call handler for " + c.operationName + " could not be parsed to JSON.")
		endCall(w, r, apiErr, d)
		return
	}
	w.Header().Set("Content-Type", "Application/JSON")
	w.Write(jb)
	endCall(w, r, nil, d)
}

////////////////////////////////////////////////////////////////////////////////
//
////////////////////////////////////////////////////////////////////////////////
func (c *Call) Method(method string) *Call {
	if c.model.bodyIn.model != nil {
		log.Critical("Cannot set method of call to GET since Reads() was set. GET calls do not have a body.")
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
	c.model.pathVars = append(c.model.pathVars, inputVar{
		typ:      "string",
		name:     name,
		desc:     desc,
		required: true,
	})
	return c
}

////////////////////////////////////////////////////////////////////////////////
//
////////////////////////////////////////////////////////////////////////////////
func (c *Call) QueryVar(name, description string, required bool) *Call {
	c.model.queryVars = append(c.model.queryVars, inputVar{
		typ:      "string",
		name:     name,
		desc:     description,
		required: required,
	})
	return c
}

////////////////////////////////////////////////////////////////////////////////
//
////////////////////////////////////////////////////////////////////////////////
func (c *Call) Reads(m interface{}) *Call {
	// No body allowed in a GET call
	if c.method == "GET" {
		log.Critical("A GET call can not use Reads() because the GET method does not have a request body.")
	}

	// Confirm the model is a struct
	if reflect.TypeOf(m).Kind() != reflect.Struct {
		log.Critical("The model given to Reads() must be of kind 'Struct'")
	}

	c.model.bodyIn.model = m
	c.model.identifyFieldTagsIn(nil)
	return c
}

////////////////////////////////////////////////////////////////////////////////
//
////////////////////////////////////////////////////////////////////////////////
func (c *Call) Returns(m interface{}) *Call {
	c.model.bodyOut.model = m
	c.model.identifyFieldTagsOut(nil)
	return c
}

////////////////////////////////////////////////////////////////////////////////
//
////////////////////////////////////////////////////////////////////////////////
func (c *Call) OperationName(name string) *Call {
	c.operationName = name
	return c
}
