package sleepy

type Call struct {
	path          string
	method        string
	operationName string
	handler       Handler
	filters       []Filter
	model         callDataModel
}

// Implement the Handler interface.
func (c *Call) ServeHTTP(http.ResponseWriter, *http.Request) {
	// Create a sleepy request + response

	// Check headers (make sure method mathes the call method)

	// Parse url/path variables and store them in the *sleepy.Request

	// Parse the request body into the reads model if applicable
	// and validate that required fields are present

	// Call filters

	// Call handler
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

// -- Example data model
