package handlers

// example of a route handler

import (
	"net/http"
	"text/template"
)

type API struct {
	Functionality string
	Path          string
	ReturnCodes   []string
}

type Data struct {
	APIs []API
}

func (s *Server) Index(w http.ResponseWriter, r *http.Request) {
	// set header html content type
	w.Header().Set("Content-Type", "text/html")

	tmp := template.Must(template.ParseFiles("templates/home.html"))

	data := Data{
		APIs: []API{
			{Functionality: "Get employee by id", Path: "GET /api/v1/employee/{id}", ReturnCodes: []string{"200", "404"}},
			{Functionality: "Create employee", Path: "POST /api/v1/employee/", ReturnCodes: []string{"201", "400"}},
			{Functionality: "Update employee", Path: "PUT /api/v1/employee/{id}", ReturnCodes: []string{"200", "400", "404"}},
			{Functionality: "Delete employee", Path: "DELETE /api/v1/employee/{id}", ReturnCodes: []string{"200", "404"}},
		},
	}
	tmp.Execute(w, data)
}
