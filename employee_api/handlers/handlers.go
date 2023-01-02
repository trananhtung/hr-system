package handlers

import (
	"HR-system/employee_api/models"
	"fmt"
	"log"
	"net/http"

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
	DB     *gorm.DB
	Router *mux.Router
	Routes []Route
}

func (s *Server) StartServe(
	host, port, user, dbname, password, sslmode string,
) {
	s.Initialize(host, port, user, dbname, password, sslmode)
	s.Run(":8080")
}

func (s *Server) CollectRoutes() []Route {
	return []Route{
		{
			"Index",
			"GET",
			"/",
			s.Index,
		},
		{
			"Create employee",
			"POST",
			"/api/v1/employee/",
			s.Create,
		},
		{
			"Get employee by id",
			"GET",
			"/api/v1/employee/{id}",
			s.Get,
		},
		{
			"Update employee",
			"PUT",
			"/api/v1/employee/{id}",
			s.Update,
		},
		{
			"Delete employee",
			"DELETE",
			"/api/v1/employee/{id}",
			s.Delete,
		},
	}

}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

func (s *Server) Initialize(host, port, user, dbname, password, sslmode string) {
	databaseUrl := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbname, password, sslmode)
	db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})

	if err != nil {
		fmt.Printf("Cannot connect to %s database", dbname)
		log.Fatal("This is the error:", err)
		return
	}

	s.DB = db

	s.DB.AutoMigrate(&models.Employee{})

	s.Router = mux.NewRouter()
	s.Routes = s.CollectRoutes()
	for _, route := range s.Routes {
		s.Router.HandleFunc(route.Pattern, route.HandlerFunc).Methods(route.Method)
	}
}
