package handlers

import (
	"HR-system/employee_api/models"
	"HR-system/employee_api/responses"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var employees []models.Employee

	idQuery := vars["id"]
	// convert string to int
	id, err := strconv.Atoi(idQuery)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, []string{"Invalid id"})
	}

	results := s.DB.First(&employees, id)
	if results.Error != nil {
		responses.Error(w, http.StatusBadRequest, []string{results.Error.Error()})
		return
	}
	responses.Success(w, http.StatusOK, employees)
}
