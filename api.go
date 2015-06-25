package sleepy

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/op/go-logging"
	"github.com/tortis/sleepy/mux"
)

var log = logging.MustGetLogger("example")
var format = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
)

////////////////////////////////////////////////////////////////////////////////
// A handler function that can service an API call.                           //
////////////////////////////////////////////////////////////////////////////////
type Handler func(http.ResponseWriter, *http.Request, CallData) (interface{}, *Error)

////////////////////////////////////////////////////////////////////////////////
// A middleware handler that can do work before the API call handler. Filters //
// can be applied to the whole API, a single resource, or a single call.      //
// Filter's should not write any data to the ResponseWriter, instead they     //
// should write data to CallData and return an Error if appropriate.          //
////////////////////////////////////////////////////////////////////////////////
type Filter func(*http.Request, CallData) *Error

////////////////////////////////////////////////////////////////////////////////
// A place to store arbitary data while the request is bounding between       //
// different calls and handlers.                                              //
////////////////////////////////////////////////////////////////////////////////
type CallData map[string]interface{}

type API struct {
	basePath       string
	resources      []*Resource
	router         *mux.Router
	resourceRouter *mux.Router
	filters        []Filter
	enableCORS     bool
}

////////////////////////////////////////////////////////////////////////////////
// Create a new API that will handle HTTP requests on the base path. Use the  //
// Register method to add resources to the API.                               //
////////////////////////////////////////////////////////////////////////////////
func New(basePath string, enableCORS bool) *API {
	api := &API{
		resources:  make([]*Resource, 0),
		router:     mux.NewRouter(),
		basePath:   basePath,
		enableCORS: enableCORS,
	}
	api.resourceRouter = api.router.PathPrefix(basePath).Subrouter()

	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")
	logging.SetBackend(backend1Leveled, backend2Formatter)
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

	// Time the application level call handling
	data["_start"] = time.Now()

	// Check for OPTIONS methods to handle CORS
	if api.enableCORS {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			endCall(w, r, nil, data)
			return
		}
	}

	// Run API level filters
	for _, filter := range api.filters {
		err := filter(r, data)
		if err != nil {
			endCall(w, r, err, data)
			return
		}
	}

	api.router.ServeHTTP(w, r, data)
}

////////////////////////////////////////////////////////////////////////////////
// Function that should be called at the very end of every single request.    //
// It is responsible for logging the request and result. If the error is nil, //
// the request will be logged as having been handled successfully.            //
// If the error is not nil, then it will be logged AND this function will     //
// write the appropriate error details to the client.                         //
////////////////////////////////////////////////////////////////////////////////
func endCall(w http.ResponseWriter, r *http.Request, e *Error, d CallData) {
	startTime := d["_start"].(time.Time)
	duration := time.Since(startTime) / 1000
	if e == nil {
		log.Notice("[200] [client %s]->[%s %s] [%d us] OK\n", r.RemoteAddr, r.Method, r.URL, duration)
	} else {
		w.Header().Set("Content-Type", "Application/JSON")
		w.WriteHeader(e.HttpCode)
		// Marshal the error into JSON
		jb, err := json.Marshal(e)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			log.Critical("%s", err.Error())
		}
		w.Write(jb)
		if e.HttpCode >= 500 {
			log.Error("[%d] [client %s]->[%s %s] [%d us] %s: %s\n", e.HttpCode, r.RemoteAddr, r.Method, r.URL, duration, e.Msg, e.Err)
		} else {
			log.Warning("[%d] [client %s]->[%s %s] [%d us] %s: %s\n", e.HttpCode, r.RemoteAddr, r.Method, r.URL, duration, e.Msg, e.Err)
		}
	}
}
