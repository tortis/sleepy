package sleepy

import (
	"net/http"

	"github.com/tortis/sleepy/mux"
)

type API struct {
	basePath       string
	resources      []*Resource
	router         *mux.Router
	resourceRouter *mux.Router
}

func New(basePath string) *API {
	api := &API{
		resources: make([]*Resource, 0),
		router:    mux.NewRouter(),
		basePath:  basePath,
	}
	api.resourceRouter = api.router.PathPrefix(basePath).Subrouter()
	return api
}

// Implement http Handler interface
func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.router.ServeHTTP(w, r)
}

// All of the resources calls need to be added to the router
func (api *API) Register(r *Resource) {
	// Add the resource to our list
	api.resources = append(api.resources, r)
	r.construct(api.basePath)
	api.resourceRouter.PathPrefix(r.path).Handler(r)
}
