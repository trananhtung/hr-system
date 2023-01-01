package main

import (
	"github/trananhtung/HR-system/db"
	"github/trananhtung/HR-system/employee"
)

func main() {
	db.ConnectDB()
	employee.EmployeeService()
}
