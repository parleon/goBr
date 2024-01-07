package gobr

import "net/http"

type IService interface {
}

type Service struct {
	middlewareStack []func(next http.HandlerFunc) http.HandlerFunc
	handlers        map[string]http.HandlerFunc
	runtime         *Runtime
}

// The returning middleware function must call next(w, r) inside to continue the chain.
func (s *Service) Use(middleware func(next http.HandlerFunc) http.HandlerFunc) {
	s.middlewareStack = append(s.middlewareStack, middleware)
}

func (s *Service) HandleFunc(pattern string, handlerFunc http.HandlerFunc) {

	if len(s.middlewareStack) > 0 {
		handlerFunc = s.chainMiddleware(handlerFunc)
	}

	s.handlers[pattern] = handlerFunc
	s.runtime.localMux.HandleFunc(pattern, handlerFunc)
}

func (s *Service) BroadcastService(port string) {
	s.runtime.addHandlersToPort(port, s.handlers)
}

// The middleware is executed recursively such that the nth function wraps the n+1th.
// e.g., in a middleware stack of [a, b, c], the order of operations will be:
//
//	a(w,r) {
//		// a logic before calling next(w,r)
//		b(w,r) {
//			// b logic before calling next(w,r)
//			c(w,r) {
//				// c logic before calling next(w,r)
//				handlerFunc(w,r)
//				// c logic after calling next(w,r)
//			}
//			// b logic after calling next(w,r)
//		}
//		// a logic after calling next(w,r)
//	}
func (s *Service) chainMiddleware(handlerFunc http.HandlerFunc) http.HandlerFunc {
	for i := len(s.middlewareStack) - 1; i >= 0; i-- {
		handlerFunc = s.middlewareStack[i](handlerFunc)
	}

	return handlerFunc
}
