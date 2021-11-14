package main

import (
	"net/http"
	"twitter/api"
	"twitter/routes"
	r "twitter/routes"

	"github.com/gorilla/mux"
)

// NewRouter creates a new *mux.Router
func NewRouter() *mux.Router {
	routes2 := r.Routes{}

	// Append routes2
	routes2.AppendRoutes(api.GetRoutes())

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes2 {
		h := addMiddleware(http.Handler(route.HandlerFunc), route)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(h)
	}
	return router
}

// addMiddleware adds all the standard middlewares to the handler.
func addMiddleware(h http.Handler, r routes.Route) http.Handler {

	h = addCustomContext(h)

	if r.VerifyJWT {
		// before THAT we will validate the token passed
		h = jwtMiddleware.Handler(h)
	}

	return h
}

// MyServer is a wrapper for mux.Router pointer
type MyServer struct {
	r *mux.Router
}

// ServeHTTP servers HTTP
func (s *MyServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
		rw.Header().Set("Access-Control-Max-Age", "86400")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-XSRF-Token, X-HTTP-Method-Override, X-Requested-With, Mobile-Cookie")
	}
	// Let Gorilla work
	s.r.ServeHTTP(rw, req)
}
