package services

import (
	"log"
	"reminder-server/internal/models"
	"reminder-server/internal/repository"

	"gorm.io/gorm"
)

type CategoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{
		repo: repository.NewCategoryRepository(db),
	}
}

func (cs *CategoryService) List(userId int32) ([]models.Category, error) {
	categories, err := cs.repo.FindByUserID(int64(userId))

	if err != nil {
		return []models.Category{}, err
	}

	return categories, nil
}

func (cs *CategoryService) Get(id int64) (models.Category, error) {
	category, err := cs.repo.FindByID(id)

	if err != nil {
		return models.Category{}, err
	}

	return category, nil
}

func (cs *CategoryService) Create(request models.CategoryCreateRequest) (models.Category, error) {
	log.Println("Creating category")

	newCategory := models.Category{
		Name:  request.Name,
		Color: request.Color,
		Icon:  request.Icon,
	}

	category, err := cs.repo.Create(newCategory)

	return category, err
}

func (cs *CategoryService) Update(id int64, request models.CategoryUpdateRequest) (models.Category, error) {
	log.Println("Updating category")

	currentCategory, err := cs.Get(id)

	if err != nil {
		return models.Category{}, err
	}

	if request.Name == nil && request.Color == nil && request.Icon == nil {
		return currentCategory, nil
	}

	if request.Name != nil {
		currentCategory.Name = *request.Name
	}

	if request.Color != nil {
		currentCategory.Color = *request.Color
	}

	if request.Icon != nil {
		currentCategory.Icon = *request.Icon
	}

	category, err := cs.repo.Update(currentCategory)

	if err != nil {
		return models.Category{}, err
	}

	return category, nil
}

func (cs *CategoryService) Delete(id int64) error {
	log.Println("Deleting category")

	err := cs.repo.Delete(id)

	return err
}

func (cs *CategoryService) CreateBulk(categories []models.Category) error {
	_, err := cs.repo.CreateBulk(categories)

	return err
}
