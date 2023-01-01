package routes

import "net/http"

func Post(w http.ResponseWriter, r *http.Request) {
	w.Header()
	id := r.URL.Path[len("/api/v1/employee/"):]
	w.Write([]byte("Create account for id: " + id))
	w.WriteHeader(http.StatusCreated)
}
