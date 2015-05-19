package sleepy

import "net/http"

type API struct {
	basePath  string
	resources []*Resource
	router    *http.ServeMux
}

func New(basePath string) *API {
	return &API{
		resources: make([]*Resource, 0),
		router:    http.NewServeMux(),
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
	api.resource = append(api.resources, r)
	// Create a new route with the resource path
	resourceRoute := api.router.PathPrefix(basePath + r.path)
	// Create a sub router on this route
	resourceRouter := resourceRoute.Subrouter()
	// Let the resource attach its handlers to the subrouter
	r.construct(resourceRouter)
	// Let the resource handle the route
	resourceRoute.Handler(r)
}
