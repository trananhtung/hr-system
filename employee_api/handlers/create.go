package handlers

import (
	"HR-system/employee_api/models"
	"HR-system/employee_api/responses"
	"encoding/json"
	"io"
	"net/http"
)

func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	// read data from request body
	body, err := io.ReadAll(r.Body)
	responses.CheckError(w, err != nil, nil)

	// convert data to Employee struct
	employee := models.Employee{}
	err = json.Unmarshal(body, &employee)
	responses.CheckError(w, err != nil, nil)

	// validate data
	messages := employee.Validate()
	responses.CheckError(w, len(messages) > 0, messages)

	tx := s.DB.Create(&employee)
	if tx.Error != nil {
		messages = []string{tx.Error.Error()}
		responses.Error(w, http.StatusBadRequest, messages)
		return
	}

	w.WriteHeader(http.StatusCreated)
	responses.Success(w, http.StatusCreated,
		struct {
			ID int `json:"id"`
		}{
			ID: employee.ID,
		})
}
