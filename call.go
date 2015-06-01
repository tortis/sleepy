package sleepy

import (
	"encoding/json"
	"errors"
	"net/http"
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
	// Create a sleepy request + response

	// Check headers (make sure method mathes the call method)

	// Parse url/path variables and store them in the *sleepy.Request

	// Parse the request body into the reads model if applicable
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

type inputVar struct {
	typ      string
	name     string
	desc     string
	required bool
}

// Call builder method
func (c *Call) Method(method string) *Call {
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
