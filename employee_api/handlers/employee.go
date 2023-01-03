package employee_handler

import (
	"HR-system/employee_api/models"
	employee_storage "HR-system/employee_api/storage"
	responses "HR-system/employee_api/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Server struct {
	Storage *employee_storage.Storage
	Router  *mux.Router
}

func (s *Server) Run(
	host, port, user, dbname, password, sslmode string,
) {
	s.initializeDB(host, port, user, dbname, password, sslmode)
	s.startServe(":8080")
}

func (s *Server) collectRoutes() []Route {
	return []Route{
		{
			"Index",
			"GET",
			"/",
			s.index,
		},
		{
			"Create employee",
			"POST",
			"/api/v1/employee/",
			s.create,
		},
		{
			"Get employee by id",
			"GET",
			"/api/v1/employee/{id}",
			s.get,
		},
		{
			"Update employee",
			"PUT",
			"/api/v1/employee/{id}",
			s.update,
		},
		{
			"Delete employee",
			"DELETE",
			"/api/v1/employee/{id}",
			s.delete,
		},
	}

}

func (server *Server) startServe(addr string) {
	log.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

func (s *Server) initializeDB(host, port, user, dbname, password, sslmode string) {
	databaseUrl := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbname, password, sslmode)
	db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})

	if err != nil {
		log.Printf("Cannot connect to %s database", dbname)
		log.Fatal("This is the error:", err)
		return
	}

	s.Storage.SetDB(db)
	s.Storage.AutoMigrate()

	s.Router = mux.NewRouter()
	Routes := s.collectRoutes()
	for _, route := range Routes {
		s.Router.HandleFunc(route.Pattern, route.HandlerFunc).Methods(route.Method)
	}
}

func (s *Server) create(w http.ResponseWriter, r *http.Request) {
	// read data from request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, []string{"Invalid request"})
		return
	}

	// convert data to Employee struct
	employee := models.Employee{}
	err = json.Unmarshal(body, &employee)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, []string{"Invalid request"})
		return
	}

	// validate data
	messages := employee.Validate()
	if len(messages) > 0 {
		responses.Error(w, http.StatusBadRequest, messages)
		return
	}

	tx := s.Storage.Create(&employee)
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

func (s *Server) delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idQuery := vars["id"]
	id, err := strconv.Atoi(idQuery)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, []string{"Invalid id"})
		return
	}

	tx := s.Storage.DeleteById(id)
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

func (s *Server) get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	idQuery := vars["id"]
	// convert string to int
	id, err := strconv.Atoi(idQuery)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, []string{"Invalid id"})
	}

	employees, err := s.Storage.GetById(id)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, []string{err.Error()})
		return
	}
	responses.Success(w, http.StatusOK, employees)
}

func (s *Server) update(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, []string{"Invalid request"})
		return
	}

	vars := mux.Vars(r)
	idQuery := vars["id"]
	id, err := strconv.Atoi(idQuery)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, []string{"Invalid request"})
		return
	}

	var updateEmployee models.Employee
	err = json.Unmarshal(body, &updateEmployee)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, []string{"Invalid request"})
		return
	}

	// update employee
	err = s.Storage.UpdateById(id, updateEmployee)
	if err != nil {
		messages := []string{err.Error()}
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

type API struct {
	Functionality string
	Path          string
	ReturnCodes   []string
}

type Data struct {
	APIs []API
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {

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
