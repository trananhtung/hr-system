package handlers

import (
	"net/http"
)

func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header()
	id := r.URL.Path[len("/api/v1/employee/"):]
	w.Write([]byte("Delete account for id: " + id))
	w.WriteHeader(http.StatusOK)
}
