package gobr

import "net/http"

type IService interface {

}

type Service struct {
	Handlers map[string]http.HandlerFunc

	runtime *Runtime
}

func (s *Service) HandleFunc(pattern string, handlerFunc http.HandlerFunc) {
	s.Handlers[pattern] = handlerFunc
	s.runtime.localMux.HandleFunc(pattern, handlerFunc)
}

func (s *Service) BroadcastService(port string) {
	s.runtime.addHandlersToPort(port, s.Handlers)
}


