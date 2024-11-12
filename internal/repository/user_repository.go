package repository

import (
	"reminder-server/internal/models"

	"gorm.io/gorm"
)

type userRepository interface {
	FindAll() ([]models.User, error)
	FindByID(id int64) (models.User, error)
	FindByEmail(email string) (models.User, error)
	Create(user models.User) (models.User, error)
	Update(user models.User) (models.User, error)
	Delete(id int64) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{
		db: db,
	}
}

func (ur *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	result := ur.db.Find(&users)

	return users, result.Error
}

func (ur *UserRepository) FindByID(id int64) (models.User, error) {
	var user models.User
	result := ur.db.First(&user, id)

	return user, result.Error
}

func (ur *UserRepository) FindByEmail(email string) (models.User, error) {
	var user models.User
	result := ur.db.Where("email = ?", email).First(&user)

	return user, result.Error
}

func (ur *UserRepository) Create(user models.User) (models.User, error) {
	result := ur.db.Create(&user)

	return user, result.Error
}

func (ur *UserRepository) Update(user models.User) (models.User, error) {
	result := ur.db.Save(&user)

	return user, result.Error
}

func (ur *UserRepository) Delete(id int64) error {
	result := ur.db.Delete(&models.User{}, id)

	return result.Error
}
