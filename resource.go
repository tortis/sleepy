package sleepy

type Resource struct {
	path    string
	name    string
	calls   []*Call
	filters []*Filter
	router  *mux.Router
}

func (r *Resource) Route(path string) *Call {
	c := &Call{path: path, operationName: path}
	r.calls = append(r.calls, c)
	return c
}

func (r *Resource) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Call all filters
	for _, filter := range r.filters {
		err := filter(w, r)
		if err != nil {
			// Fail here and stop handling the request
		}
	}
	r.router.ServeHTTP(w, r)
}

func (r *Resource) construct(router *mux.Router) {
	r.router = router
	for _, call := range r.calls {
		r.router.Handle(call.path, call).Methods(call.Method)
	}
}
