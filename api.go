package sleepy

import (
	"fmt"
	"net/http"

	"github.com/tortis/sleepy/mux"
)

////////////////////////////////////////////////////////////////////////////////
// A place to store arbitary data while the request is bounding between       //
// different calls and handlers.                                              //
////////////////////////////////////////////////////////////////////////////////
type CallData map[string]interface{}

////////////////////////////////////////////////////////////////////////////////
// Then end handler function that will service an API call.                   //
////////////////////////////////////////////////////////////////////////////////
type Handler func(http.ResponseWriter, *http.Request, CallData) (interface{}, Error)

////////////////////////////////////////////////////////////////////////////////
// A middleware handler that can do work before the API call handler. Filters //
// can be applied to the whole API, a single resource, or a single call.      //
// Filter's should not write any data to the ResponseWriter, instead they     //
// should write data to CallData and return an Error if appropriate.          //
////////////////////////////////////////////////////////////////////////////////
type Filter func(http.ResponseWriter, *http.Request, CallData) Error

type API struct {
	basePath       string
	resources      []*Resource
	router         *mux.Router
	resourceRouter *mux.Router
	filters        []Filter
}

////////////////////////////////////////////////////////////////////////////////
// Create a new API that will handle HTTP requests on the base path. Use the  //
// Register method to add resources to the API.                               //
////////////////////////////////////////////////////////////////////////////////
func New(basePath string) *API {
	api := &API{
		resources: make([]*Resource, 0),
		router:    mux.NewRouter(),
		basePath:  basePath,
	}
	api.resourceRouter = api.router.PathPrefix(basePath).Subrouter()
	return api
}

////////////////////////////////////////////////////////////////////////////////
// Add a filter to the whole API. These filters will run before any           //
// resource or call filters, and before the call handler. The filters will    //
// run in the order that they are added.                                      //
////////////////////////////////////////////////////////////////////////////////
func (api *API) Filter(f Filter) {
	api.filters = append(api.filters, f)
}

////////////////////////////////////////////////////////////////////////////////
// Adds a resource to the API. This will result in the API handling all of    //
// the resource's calls at '/base/path/resourcepath/call/path'.               //
////////////////////////////////////////////////////////////////////////////////
func (api *API) Register(r *Resource) {
	// Add the resource to our list
	api.resources = append(api.resources, r)
	r.construct(api.basePath)
	api.resourceRouter.PathPrefix(r.path).Handler(r)
}

// Implement http Handler interface
func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Create data space
	data := make(map[string]interface{})

	// Run API level filters
	for _, filter := range api.filters {
		err := filter(w, r, data)
		if err != nil {
			http.Error(w, err.Message(), err.StatusCode())
			logResult(r, err)
			return
		}
	}

	api.router.ServeHTTP(w, r, data)
}

// Should always be called before response handling ends
func logResult(r *http.Request, e Error) {
	if e == nil {
		fmt.Printf("[200]  [client %s]->[%s %s]  OK\n", r.RemoteAddr, r.URL)
	} else {
		fmt.Printf("[%d]  [client %s]->[%s %s] %s: %s\n", e.StatusCode(), r.RemoteAddr, r.URL, e.Message(), e.Error())
	}
}
