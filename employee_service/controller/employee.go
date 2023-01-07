package employee_controller

import (
	"HR-system/employee_service/models"
	employee_storage "HR-system/employee_service/storage"
	responses "HR-system/employee_service/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	Storage *employee_storage.Storage
	Router  *gin.Engine
}

func (s *Server) Run(
	host, port, user, dbname, password, sslmode string,
) {
	s.initializeDB(host, port, user, dbname, password, sslmode)
	s.startServe(":8080")
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

	s.Router = gin.Default()

	s.Router.LoadHTMLGlob("templates/*")

	s.Router.GET("/", s.index)
	s.Router.POST("/api/v1/employee/", s.create)
	s.Router.GET("/api/v1/employee/:id", s.get)
	s.Router.PUT("/api/v1/employee/:id", s.update)
	s.Router.DELETE("/api/v1/employee/:id", s.delete)
}

type successResponseID struct {
	ID int `json:"id"`
}

func (s *Server) create(c *gin.Context) {
	var employee models.Employee

	if err := c.ShouldBindJSON(&employee); err != nil {
		responses.BadRequest(c, []string{err.Error()})
		return
	}

	// validate data
	messages := employee.Validate()
	if len(messages) > 0 {
		responses.BadRequest(c, messages)
		return
	}

	tx := s.Storage.Create(&employee)
	if tx.Error != nil {
		messages = []string{tx.Error.Error()}
		responses.BadRequest(c, messages)
		return
	}

	responses.Success(c, successResponseID{ID: employee.ID})
}

func (s *Server) delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		responses.BadRequest(c, []string{"Invalid id"})
	}

	tx := s.Storage.DeleteById(id)
	if tx.Error != nil {
		messages := []string{tx.Error.Error()}
		responses.NotFound(c, messages)
	}

	responses.Success(c, successResponseID{ID: id})
}

func (s *Server) get(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		responses.BadRequest(c, []string{"Invalid id"})
	}

	employees, err := s.Storage.GetById(id)
	if err != nil {
		responses.NotFound(c, []string{err.Error()})
	}
	responses.Success(c, employees)
}

func (s *Server) update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		responses.BadRequest(c, []string{"Invalid request"})
		return
	}

	var updateEmployee models.Employee

	if err := c.ShouldBindJSON(&updateEmployee); err != nil {
		responses.NotFound(c, []string{"Invalid request"})
		return
	}

	// update employee
	err = s.Storage.UpdateById(id, updateEmployee)
	if err != nil {
		messages := []string{err.Error()}
		responses.BadRequest(c, messages)
		return
	}

	responses.Success(c, successResponseID{ID: id})
}

type API struct {
	Functionality string
	Path          string
	ReturnCodes   []string
}

type Data struct {
	APIs []API
}

func (s *Server) index(c *gin.Context) {
	data := Data{
		APIs: []API{
			{Functionality: "Get employee by id", Path: "GET /api/v1/employee/{id}", ReturnCodes: []string{"200", "404"}},
			{Functionality: "Create employee", Path: "POST /api/v1/employee/", ReturnCodes: []string{"201", "400"}},
			{Functionality: "Update employee", Path: "PUT /api/v1/employee/{id}", ReturnCodes: []string{"200", "400", "404"}},
			{Functionality: "Delete employee", Path: "DELETE /api/v1/employee/{id}", ReturnCodes: []string{"200", "404"}},
		},
	}

	c.HTML(http.StatusOK, "home.html", data)
}
