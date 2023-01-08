package models

import (
	"fmt"
	"time"

	"github.com/go-playground/validator"
)

// create enum for required fields
const (
	REQUIRE  = true
	OPTIONAL = false
)

var validate *validator.Validate

type EmployeeDTO struct {
	ID        uint   `json:"id" gorm:"<-:create,primarykey"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" gorm:"unique"`
	Phone     string `json:"phone"`
	Birthday  string `json:"birthday"`
	StartDay  string `json:"start_day"`
	Position  string `json:"position"`

	Require bool `json:"-"`
}

type EmployeeModel struct {
	ID        uint   `json:"id" gorm:"primarykey"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" gorm:"unique"`
	Phone     string `json:"phone"`
	Birthday  string `json:"birthday"`
	StartDay  string `json:"start_day"`
	Position  string `json:"position"`
	CreateAt  int64  `json:"create_at"`
	UpdateAt  int64  `json:"update_at"`
	DeleteAt  int64  `json:"delete_at"`
}

func (e *EmployeeDTO) MapForUpdate() EmployeeModel {
	var employeeDB EmployeeModel
	employeeDB.FirstName = e.FirstName
	employeeDB.LastName = e.LastName
	employeeDB.Email = e.Email
	employeeDB.Phone = e.Phone
	employeeDB.Birthday = e.Birthday
	employeeDB.StartDay = e.StartDay
	employeeDB.Position = e.Position
	employeeDB.UpdateAt = time.Now().Unix()

	return employeeDB
}

func (e *EmployeeDTO) MapForCreate() EmployeeModel {
	employeeDB := e.MapForUpdate()
	employeeDB.ID = e.ID
	employeeDB.CreateAt = time.Now().Unix()

	return employeeDB
}

func (e *EmployeeDTO) SetRequired(required bool) {
	e.Require = required
}

func formatTag(tag string, required bool) string {
	if required {
		return fmt.Sprintf("required,%s", tag)
	}
	return tag
}

func (e *EmployeeDTO) Validate() []string {

	var err error
	var messages []string
	validate = validator.New()
	validate.RegisterValidation("yyyy-mm-dd", YYYYMMDDValidator)
	err = validate.Var(e.Email, formatTag("email", e.Require))
	if err != nil {
		messages = append(messages, err.Error())
	}

	err = validate.Var(e.FirstName, formatTag("alpha", e.Require))
	if err != nil {
		messages = append(messages, err.Error())
	}

	err = validate.Var(e.LastName, formatTag("alpha", e.Require))
	if err != nil {
		messages = append(messages, err.Error())
	}

	err = validate.Var(e.Phone, formatTag("numeric", e.Require))
	if err != nil {
		messages = append(messages, err.Error())
	}

	err = validate.Var(e.Birthday, "yyyy-mm-dd")
	if err != nil {
		messages = append(messages, err.Error())
	}

	err = validate.Var(e.StartDay, "yyyy-mm-dd")
	if err != nil {
		messages = append(messages, err.Error())
	}

	err = validate.Var(e.Position, formatTag("alpha", e.Require))
	if err != nil {
		messages = append(messages, err.Error())
	}

	return messages
}

func YYYYMMDDValidator(fl validator.FieldLevel) bool {
	timeStr := fl.Field().String()
	// pass validation if yyyy-mm-dd is empty
	if timeStr == "" {
		return true
	}
	// check yyyy-mm-dd format yyyy-mm-dd
	_, err := time.Parse("2006-01-02", timeStr)
	return err == nil
}
