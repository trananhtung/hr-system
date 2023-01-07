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
	return s.db.AutoMigrate(&models.Employee{})
}

func (s *Storage) Create(value interface{}) *gorm.DB {
	return s.db.Create(value)
}

func (s *Storage) DeleteById(id int) *gorm.DB {
	return s.db.Delete(&models.Employee{}, id)
}

func (s *Storage) GetById(id int) ([]models.Employee, error) {
	var employees []models.Employee
	result := s.db.First(&employees, id)
	return employees, result.Error
}

func (s *Storage) UpdateById(id int, updateEmployee models.Employee) error {
	tx := s.db.Model(&models.Employee{}).Where("id = ?", id).Updates(updateEmployee)
	return tx.Error
}
