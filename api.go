package sleepy

import "github.com/gorilla/mux"

type API struct {
	resources []*Resource
	router    *mux.Router
}

func New() *API {
	return &API{
		resources: make([]*Resource, 0),
		router:    mux.NewRouter(),
	}
}

func (api *API) GetRouter() *mux.Router {
	return api.router
}

// All of the resources calls need to be added to the router
func (api *API) Register(r *Resource) {
	// Add the resource to our local list
	api.resource = append(api.resources, r)

	// Attach handler for each call
	for _, call := range r.calls {

	}
}
