package handlers

import (
	"HR-system/employee_api/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var posts []models.Employee

func (s *Server) Get(w http.ResponseWriter, r *http.Request) {
	posts = []models.Employee{
		{ID: 1, FirstName: "Tung", LastName: "Tran", Email: "tungtran@gmail.com", Phone: "0123456789", Birthday: "1999-01-01", Position: "Developer"},
		{ID: 2, FirstName: "Kiet", LastName: "Tran", Email: "kiettran@gmail.com", Phone: "0123456782", Birthday: "1999-02-01", Position: "Tester"},
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	id := vars["id"]
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
	message := fmt.Sprintf("404 - Employee with id %d not found", ID)
	w.Write([]byte(message))

}
