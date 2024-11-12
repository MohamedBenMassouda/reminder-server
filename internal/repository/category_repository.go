package repository

import (
	"reminder-server/internal/models"

	"gorm.io/gorm"
)

type categoryRepository interface {
	FindAll() ([]models.Category, error)
	FindByID(id int64) (models.Category, error)
	FindByUserID(userID int64) ([]models.Category, error)
	Create(category models.Category) (models.Category, error)
	CreateBulk(categories []models.Category) ([]models.Category, error)
	Update(category models.Category) (models.Category, error)
	Delete(id int64) error
}

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return CategoryRepository{
		db: db,
	}
}

func (cr *CategoryRepository) FindAll() ([]models.Category, error) {
	var categories []models.Category
	result := cr.db.Find(&categories)

	return categories, result.Error
}

func (cr *CategoryRepository) FindByID(id int64) (models.Category, error) {
	var category models.Category
	result := cr.db.First(&category, id)

	return category, result.Error
}

func (cr *CategoryRepository) FindByUserID(userID int64) ([]models.Category, error) {
	var categories []models.Category
	result := cr.db.Where("user_id = ?", userID).Find(&categories)

	return categories, result.Error
}

func (cr *CategoryRepository) Create(category models.Category) (models.Category, error) {
	result := cr.db.Create(&category)

	return category, result.Error
}

func (cr *CategoryRepository) Update(category models.Category) (models.Category, error) {
	result := cr.db.Save(&category)

	return category, result.Error
}

func (cr *CategoryRepository) Delete(id int64) error {
	result := cr.db.Delete(&models.Category{}, id)

	return result.Error
}

func (cr *CategoryRepository) CreateBulk(categories []models.Category) ([]models.Category, error) {
	result := cr.db.Create(&categories)

	return categories, result.Error
}
