package models

import (
	"net/mail"
	"strconv"
	"strings"
	"time"
)

type Employee struct {
	ID        int    `json:"id" binding:"required,number" gorm:"primary_key"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"required,number"`
	Birthday  string `json:"birthday" binding:"required,datetime=2006-01-02"`
	StartDay  string `json:"start_day" binding:"required,datetime=2006-01-02"`
	Position  string `json:"position" binding:"required,oneof=developer tester manager"`
}

func (e *Employee) validateFirstName() []string {
	var messages = make([]string, 0)
	if len(e.FirstName) == 0 {
		messages = append(messages, "First name is required")
	}

	// change first character to uppercase
	e.FirstName = strings.ToTitle(e.FirstName)
	return messages
}

func (e *Employee) validateLastName() []string {
	// make to create a new slice
	var messages = make([]string, 0)
	if len(e.LastName) == 0 {
		messages = append(messages, "Last name is required")
	}

	// change first character to uppercase
	e.LastName = strings.ToTitle(e.LastName)
	return messages
}

func (e *Employee) validateEmail() []string {
	var messages = make([]string, 0)
	email, err := mail.ParseAddress(e.Email)
	if err != nil {
		messages = append(messages, "Email is invalid")
	}
	e.Email = email.Address
	return messages
}

func (e *Employee) validatePhone() []string {
	var messages = make([]string, 0)
	// number only
	_, err := strconv.Atoi(e.Phone)
	if err != nil {
		messages = append(messages, "Phone is invalid")
	}
	return messages
}

func (e *Employee) validateBirthday() []string {
	var messages = make([]string, 0)
	// check validate date format yyyy-mm-dd
	t, err := time.Parse("2006-01-02", e.Birthday)
	if err != nil {
		messages = append(messages, "Birthday is invalid")
	}
	// check > 18 years old
	if err == nil && time.Now().Year()-t.Year() < 18 {
		messages = append(messages, "Birthday is invalid")
	}
	return messages
}

func (e *Employee) validateStartDay() []string {
	var messages = make([]string, 0)
	// check validate date format yyyy-mm-dd
	t, err := time.Parse("2006-01-02", e.StartDay)
	if err != nil {
		messages = append(messages, "Start day is invalid")
	}

	// check future date
	if err == nil && t.After(time.Now()) {
		messages = append(messages, "Start day is future date")
	}

	return messages
}

func (e *Employee) validatePosition() []string {
	var messages = make([]string, 0)
	var listPosition = [3]string{"Manager", "Developer", "Tester"}
	var isValid = false
	for _, position := range listPosition {
		if position == e.Position {
			isValid = true
			break
		}
	}

	if isValid {
		messages = append(messages, "Position is invalid")
	}
	return messages
}

func (e *Employee) Validate() []string {
	var messages []string

	messages = append(messages, e.validateFirstName()...)
	messages = append(messages, e.validateLastName()...)
	messages = append(messages, e.validateEmail()...)
	messages = append(messages, e.validatePhone()...)
	messages = append(messages, e.validateBirthday()...)
	messages = append(messages, e.validateStartDay()...)
	messages = append(messages, e.validatePosition()...)

	return messages
}
