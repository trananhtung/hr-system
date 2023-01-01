package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Employee struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Birthday  string `json:"birthday"`
	Position  string `json:"position"`
}

var posts []Employee

func Get(w http.ResponseWriter, r *http.Request) {
	posts = []Employee{
		{ID: 1, FirstName: "Tung", LastName: "Tran", Email: "tungtran@gmail.com", Phone: "0123456789", Birthday: "1999-01-01", Position: "Developer"},
		{ID: 2, FirstName: "Kiet", LastName: "Tran", Email: "kiettran@gmail.com", Phone: "0123456782", Birthday: "1999-02-01", Position: "Tester"},
	}

	// set content type to json
	w.Header().Set("Content-Type", "application/json")
	// GET /api/v1/employee/{id}
	id := r.URL.Path[len("/api/v1/employee/"):]
	ID, error := strconv.Atoi(id)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID"))
		return
	}
	// loop through posts and find one with the id from the request
	for _, post := range posts {
		if post.ID == ID {
			// convert the post to json
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(post)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 - Employee not found"))

}
