package employee_controller

import (
	"HR-system/employee_service/models"
	employee_storage "HR-system/employee_service/storage"
	responses "HR-system/employee_service/utils"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	Storage *employee_storage.Storage
	Router  *gin.Engine
}

type loginAccount struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var identityKey = "id"

type User struct {
	Email     string
	FirstName string
	LastName  string
	Position  string
}

func (s *Server) Run(
	host, port, user, dbname, password, sslmode, jwtKey string,
) {
	s.initializeDB(host, port, user, dbname, password, sslmode, jwtKey)
	s.startServe(":8080")
}

func (server *Server) startServe(addr string) {
	log.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

func (s *Server) initializeDB(host, port, user, dbname, password, sslmode, jwtKey string) {
	databaseUrl := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbname, password, sslmode)
	db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})

	v1Route := "/api/v1/employee"

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

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "authentication",
		Key:         []byte(jwtKey),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.Email,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				Email: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginAccountVal loginAccount
			if err := c.ShouldBind(&loginAccountVal); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			email := loginAccountVal.Email
			password := loginAccountVal.Password

			employees, err := s.Storage.GetByEmail(email)
			dbEmail := employees[0].Email
			dbHashPassword := employees[0].Password

			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			errComparePassword := bcrypt.CompareHashAndPassword([]byte(dbHashPassword), []byte(password))
			if errComparePassword == nil && dbEmail == email {
				return &User{
					Email:     email,
					LastName:  employees[0].LastName,
					FirstName: employees[0].FirstName,
					Position:  employees[0].Position,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if user, ok := data.(*User); ok {
				employees, err := s.Storage.GetByEmail(user.Email)

				if err != nil {
					return false
				}

				switch employees[0].Position {
				case "manager":
					return true
				// prevent developer get all data
				case "developer":
					if c.FullPath() == v1Route && c.Request.Method == "GET" {
						return false
					}
					return true
				default:
					return false
				}
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("Middleware init error:" + errInit.Error())
	}

	s.Router.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		responses.BadRequest(c, []string{"Page not found"})
	})

	v1 := s.Router.Group(v1Route)

	// No need login route
	{
		s.Router.LoadHTMLGlob("templates/*")
		s.Router.POST("/login", authMiddleware.LoginHandler)
		v1.GET("/refresh_token", authMiddleware.RefreshHandler)
		v1.POST("/", s.create)
	}

	v1.Use(authMiddleware.MiddlewareFunc())

	// Need login route
	{
		v1.GET("/", s.index)
		v1.GET("/all", s.getAll)
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

func (s *Server) getAll(c *gin.Context) {
	employees, err := s.Storage.GetAll()
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
			{Functionality: "Get all employees", Path: "GET /api/v1/employee/all", ReturnCodes: []string{"200", "404"}},
			{Functionality: "Create employee", Path: "POST /api/v1/employee/", ReturnCodes: []string{"201", "400"}},
			{Functionality: "Update employee", Path: "PUT /api/v1/employee/{id}", ReturnCodes: []string{"200", "400", "404"}},
			{Functionality: "Delete employee", Path: "DELETE /api/v1/employee/{id}", ReturnCodes: []string{"200", "404"}},
		},
	}

	c.HTML(http.StatusOK, "home.html", data)
}
