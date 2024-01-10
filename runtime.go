package gobr

import (
	"net/http"
	"net/http/httptest"
)

type IRuntime interface {
	NewService(port string) *Service
	NewClient() *Client

	Run()
	Refresh()
}

type Runtime struct {
	Active bool

	Services  map[string]*Service
	localMux  *http.ServeMux
	portMux   map[string]*http.ServeMux
	openPorts map[string]*http.Server
}

func NewRuntime() *Runtime {
	return &Runtime{
		Active:    false,
		Services:  make(map[string]*Service),
		localMux:  http.NewServeMux(),
		portMux:   make(map[string]*http.ServeMux),
		openPorts: make(map[string]*http.Server),
	}
}

func (r *Runtime) NewService(cname string, port string) *Service {
	ns := &Service{
		Port:            port,
		Enabled:         true,
		runtime:         r,
		middlewareStack: make([]func(next http.HandlerFunc) http.HandlerFunc, 0),
		handlers:        make(map[string]http.HandlerFunc),
	}

	r.Services[cname] = ns
	return ns
}

func (r *Runtime) NewClient() *Client {
	return &Client{
		runtime: r,
	}
}

func (r *Runtime) servePort(port string) {
	server := &http.Server{Addr: ":" + port, Handler: r.portMux[port]}
	go server.ListenAndServe()
	r.openPorts[port] = server
}

func (r *Runtime) closePort(port string) {
	r.openPorts[port].Close()
	delete(r.openPorts, port)
}

func (r *Runtime) ServeAllPorts() {
	for port := range r.portMux {
		r.servePort(port)
	}
}

func (r *Runtime) CloseAllPorts() {
	for port := range r.openPorts {
		r.closePort(port)
	}
}

func (r *Runtime) addHandlersToPort(port string, handlers map[string]http.HandlerFunc) {
	portMux, exists := r.portMux[port]

	if !exists {
		portMux = http.NewServeMux()
		r.portMux[port] = portMux
	}

	for pattern, handlerFunc := range handlers {
		portMux.HandleFunc(pattern, handlerFunc)
	}
}

func (r *Runtime) refreshServices() {
	r.portMux = make(map[string]*http.ServeMux)
	localMux := http.NewServeMux()
	for _, service := range r.Services {
		if service.Enabled {

			for pattern, handlerFunc := range service.handlers {
				localMux.HandleFunc(pattern, handlerFunc)
			}

			r.addHandlersToPort(service.Port, service.handlers)
		}
	}
	r.localMux = localMux
}

func (r *Runtime) Open() {
	r.refreshServices()
	r.ServeAllPorts()
	r.Active = true
}

func (r *Runtime) Close() {
	r.CloseAllPorts()
	r.localMux = http.NewServeMux()
	r.Active = false
}

func (r *Runtime) Reboot() {
	r.Close()
	r.Open()
}

// do() will first attempt to complete the request locally before sending a remote request.
// It does this by checking if the request matches a pattern in the localMux.
//
// If there is a local match, then the request is simulated by executeLocal()
// Otherwise, the request is handled by the default behavior of http.Client.Do().
func (r *Runtime) do(req *http.Request) (*http.Response, error) {
	function, found := r.findLocalHandlerFunc(req)

	if found {
		resp, err := executeLocal(function, req)
		return resp, err
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
