package sleepy

import (
	"net/http"

	"github.com/tortis/sleepy/mux"
)

type Resource struct {
	path    string
	name    string
	calls   []*Call
	filters []Filter
	router  *mux.Router
}

func NewResource(path string) *Resource {
	return &Resource{path: path, router: mux.NewRouter()}
}

// Start call builder
func (r *Resource) Route(path string) *Call {
	c := &Call{path: path, operationName: path, filters: make([]Filter, 0)}
	r.calls = append(r.calls, c)
	return c
}

func (r *Resource) Filter(f Filter) {
	r.filters = append(r.filters, f)
}

// ServeHTTP of a resource is used by any path with the Resource.path prefix.
// This handler will apply any set filters and then use a subrouter to determine
// which call handler should be used.
// The construct() method should be called prior to serving any requests.
func (res *Resource) ServeHTTP(w http.ResponseWriter, r *http.Request, d map[string]interface{}) {
	// Call all filters
	for _, filter := range res.filters {
		err := filter(w, r, d)
		if err != nil {
			http.Error(w, err.Message(), err.StatusCode())
			logResult(r, err)
			return
		}
	}

	// Route to the appropriate call handler
	res.router.ServeHTTP(w, r, d)
}

// Give the resource a subrouter for its base path so that it can attach its
// call handlers to their respective paths.
func (r *Resource) construct(pathPrefix string) {
	for _, call := range r.calls {
		r.router.Handle(pathPrefix+r.path+call.path, call).Methods(call.method)
	}
}
