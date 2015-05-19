package sleepy

import (
	"net/http"

	"github.com/gorilla/mux"
)

type API struct {
	basePath  string
	resources []*Resource
	router    *mux.Router
}

func New(basePath string) *API {
	return &API{
		resources: make([]*Resource, 0),
		router:    mux.NewRouter(),
		basePath:  basePath,
	}
}

// Implement http Handler interface
func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.router.ServeHTTP(w, r)
}

// All of the resources calls need to be added to the router
func (api *API) Register(r *Resource) {
	// Add the resource to our list
	api.resources = append(api.resources, r)
	// Create a new route with the resource path
	resourceRoute := api.router.PathPrefix(api.basePath + r.path)
	// Create a sub router on this route
	resourceRouter := resourceRoute.Subrouter()
	// Let the resource attach its handlers to the subrouter
	r.construct(resourceRouter)
	// Let the resource handle the route
	resourceRoute.Handler(r)
}
