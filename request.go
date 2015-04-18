package sleepy

// sleep Request wraps http.Requst and holds an Items map that Filters may
// use to store data for the handler to use later
type Request struct {
	Req   *http.Request
	Items map[string]interface{}
}
