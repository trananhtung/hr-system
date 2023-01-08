package employee_controller

import (
	"HR-system/employee_service/models"
	employee_storage "HR-system/employee_service/storage"
	responses "HR-system/employee_service/utils"
	"fmt"
	"log"
	"net/http"

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
	err = s.Storage.AutoMigrate()
	if err != nil {
		log.Printf("Cannot migrate %s database", dbname)
		return
	}

	s.Router = gin.Default()

	s.Router.LoadHTMLGlob("templates/*")

	v1 := s.Router.Group("/api/v1/employee")
	{
		v1.GET("/", s.index)
		v1.POST("/", s.create)
		v1.GET("/:id", s.get)
		v1.PUT("/:id", s.update)
		v1.DELETE("/:id", s.delete)
	}
}

type ID struct {
	ID uint `json:"id" uri:"id"`
}

func (s *Server) create(c *gin.Context) {
	var employee models.EmployeeDTO

	if err := c.ShouldBindJSON(&employee); err != nil {
		responses.BadRequest(c, []string{err.Error()})
		return
	}

	employee.SetRequired(true)
	messages := employee.Validate()
	if len(messages) > 0 {
		responses.BadRequest(c, messages)
		return
	}

	createEmployee := employee.MapForCreate()

	// validate data

	tx := s.Storage.Create(&createEmployee)
	if tx.Error != nil {
		messages = []string{tx.Error.Error()}
		responses.BadRequest(c, messages)
		return
	}

	responses.Success(c, ID{ID: createEmployee.ID})
}

func (s *Server) delete(c *gin.Context) {
	var id ID
	err := c.ShouldBindUri(&id)
	if err != nil {
		responses.BadRequest(c, []string{"Invalid id"})
		return
	}

	tx := s.Storage.DeleteById(id.ID)
	if tx.Error != nil {
		messages := []string{tx.Error.Error()}
		responses.NotFound(c, messages)
		return
	}

	responses.Success(c, id)
}

func (s *Server) get(c *gin.Context) {
	var id ID
	err := c.ShouldBindUri(&id)
	if err != nil {
		responses.BadRequest(c, []string{"Invalid id"})
		return
	}

	employees, err := s.Storage.GetById(id.ID)
	if err != nil {
		responses.NotFound(c, []string{err.Error()})
		return
	}
	responses.Success(c, employees)
}

func (s *Server) update(c *gin.Context) {
	var id ID
	err := c.ShouldBindUri(&id)
	if err != nil {
		responses.BadRequest(c, []string{"Invalid id"})
		return
	}

	// get employee
	var employee models.EmployeeDTO
	err = c.ShouldBindJSON(&employee)
	if err != nil {
		responses.BadRequest(c, []string{err.Error()})
		fmt.Println(err)
		return
	}

	employee.SetRequired(false)
	messages := employee.Validate()
	if len(messages) > 0 {
		responses.BadRequest(c, messages)
		return
	}

	updateEmployee := employee.MapForUpdate()
	_, err = s.Storage.UpdateById(id.ID, updateEmployee)
	if err != nil {
		messages := []string{err.Error()}
		responses.BadRequest(c, messages)
		return
	}

	responses.Success(c, id)
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
