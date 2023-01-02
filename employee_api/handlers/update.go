package handlers

import (
	"HR-system/employee_api/models"
	"HR-system/employee_api/responses"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) Update(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	responses.CheckError(w, err != nil, nil)

	vars := mux.Vars(r)
	idQuery := vars["id"]
	id, err := strconv.Atoi(idQuery)
	responses.CheckError(w, err != nil, nil)

	var updateEmployee models.Employee
	err = json.Unmarshal(body, &updateEmployee)
	responses.CheckError(w, err != nil, nil)

	// update employee
	tx := s.DB.Model(&models.Employee{}).Where("id = ?", id).Updates(updateEmployee)
	if tx.Error != nil {
		messages := []string{tx.Error.Error()}
		responses.Error(w, http.StatusBadRequest, messages)
		return
	}

	w.WriteHeader(http.StatusOK)
	responses.Success(w, http.StatusOK, struct {
		ID int `json:"id"`
	}{
		ID: id,
	})
}
