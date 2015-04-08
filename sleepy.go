package sleepy

// Data model for a type that will be produced or consumed by an api call.
type DataModel interface{}

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
	Params       DataModel
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
