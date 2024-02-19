package web

import "net/http"

type (
	Server struct {
		Path
	}
)

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	
}
