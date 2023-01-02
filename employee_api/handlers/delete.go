package handlers

import (
	"HR-system/employee_api/models"
	"HR-system/employee_api/responses"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idQuery := vars["id"]
	id, err := strconv.Atoi(idQuery)
	responses.CheckError(w, err != nil, nil)

	tx := s.DB.Delete(&models.Employee{}, id)
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
