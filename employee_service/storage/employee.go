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

func (s *Storage) GetById(id uint) ([]models.EmployeeModel, error) {
	var employees []models.EmployeeModel
	result := s.db.First(&employees, id)
	return employees, result.Error
}

func (s *Storage) GetAll() ([]models.EmployeeModel, error) {
	var employees []models.EmployeeModel
	result := s.db.Find(&employees)
	return employees, result.Error
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
