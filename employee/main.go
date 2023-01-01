package employee

import (
	"github/trananhtung/HR-system/employee/routes"
	"net/http"
)

func checkMethod(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		routes.Get(w, r)
	case "PUT":
		routes.Put(w, r)
	case "DELETE":
		routes.Delete(w, r)
	case "POST":
		routes.Post(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func EmployeeService() {
	http.HandleFunc("/", routes.HomePage)
	http.HandleFunc("/api/v1/employee/", checkMethod)
	http.ListenAndServe(":8080", nil)
}
