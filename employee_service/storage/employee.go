package employee_storage

import (
	"HR-system/employee_service/models"

	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func (s *Storage) SetDB(db *gorm.DB) {
	s.db = db
}

func (s *Storage) AutoMigrate() error {
	return s.db.AutoMigrate(&models.EmployeeModel{})
}

func (s *Storage) Create(value interface{}) *gorm.DB {
	return s.db.Create(value)
}

func (s *Storage) DeleteById(id uint) *gorm.DB {
	return s.db.Delete(&models.EmployeeModel{}, id)
}

type EmployeeData struct {
	ID        uint   `json:"id" gorm:"primarykey"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" gorm:"unique"`
	Phone     string `json:"phone"`
	Birthday  string `json:"birthday"`
	Position  string `json:"position"`
}

func MapFromEmployeeModel(employees []models.EmployeeModel) []EmployeeData {
	var employeeData []EmployeeData
	for _, e := range employees {
		employeeData = append(employeeData, EmployeeData{
			ID:        e.ID,
			FirstName: e.FirstName,
			LastName:  e.LastName,
			Email:     e.Email,
			Phone:     e.Phone,
			Birthday:  e.Birthday,
			Position:  e.Position,
		})
	}
	return employeeData
}

func (s *Storage) GetById(id uint) ([]EmployeeData, error) {
	var employees []models.EmployeeModel
	result := s.db.First(&employees, id)

	return MapFromEmployeeModel(employees), result.Error
}

func (s *Storage) GetAll() ([]EmployeeData, error) {
	var employees []models.EmployeeModel
	result := s.db.Find(&employees)

	return MapFromEmployeeModel(employees), result.Error
}

func (s *Storage) UpdateById(id uint, updateEmployee models.EmployeeModel) (int64, error) {
	tx := s.db.Model(&models.EmployeeModel{}).Where("id = ?", id).Updates(updateEmployee)
	return tx.RowsAffected, tx.Error
}

func (s *Storage) GetByEmail(email string) ([]models.EmployeeModel, error) {
	var employees []models.EmployeeModel
	result := s.db.First(&employees, "email = ?", email)
	return employees, result.Error
}
