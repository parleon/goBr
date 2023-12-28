package gobr

import (
	"net/http"
	"net/http/httptest"

)

type IRuntime interface {
	NewService() *Service
	NewClient() *Client
}

type Runtime struct {
	localMux *http.ServeMux

	broadcastPorts map[string]*http.ServeMux
}

func NewRuntime() *Runtime {
	return &Runtime{
		localMux: http.NewServeMux(),
		broadcastPorts: make(map[string]*http.ServeMux),
	}
}

func (r *Runtime) NewService() *Service {
	return &Service{
		runtime: r,
		Handlers: make(map[string]http.HandlerFunc),
	}
}

func (r *Runtime) NewClient() *Client {
	return &Client{
		runtime: r,
	}
}

func (r *Runtime) addHandlersToPort(port string, handlers map[string]http.HandlerFunc) {
	portMux, exists := r.broadcastPorts[port]

	if !exists {
		portMux = http.NewServeMux()
		r.broadcastPorts[port] = portMux
	}

	for pattern, handlerFunc := range handlers {
		portMux.HandleFunc(pattern, handlerFunc)
	}
}

func (r *Runtime) ServeAllPorts() {
	for port, mux := range r.broadcastPorts {
		go http.ListenAndServe(port, mux)
	}
}

// do() will first attempt to complete the request locally before sending a remote request.
// It does this by checking if the request matches a pattern in the localMux.
//
// If there is a local match, then the request is simulated by executeLocal()
// Otherwise, the request is handled by the default behavior of http.Client.Do().
func (r *Runtime) do(req *http.Request) (*http.Response, error) {
	function, found := r.findLocalHandlerFunc(req)

	if found {
		return executeLocal(function , req)
	}

	client := http.Client{}

	return client.Do(req)
}

func (r *Runtime) findLocalHandlerFunc(req *http.Request) (http.HandlerFunc, bool) {
	handler, pattern := r.localMux.Handler(req)

	if pattern == "" {
		return nil, false
	}

	return handler.ServeHTTP, true
}

// TODO: replace with a more robust solution that doesn't use httptest.ResponseRecorder
// executeLocal() simulates a request by passing in a httptest.ResponseRecorder as the http.ResponseWriter
func executeLocal(handlerFunc http.HandlerFunc, req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()

	handlerFunc(w, req)

	return w.Result(), nil
}



