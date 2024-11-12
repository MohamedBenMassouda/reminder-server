package services

import (
	"errors"
	"log"
	"reminder-server/internal/models"
	"reminder-server/internal/repository"
	"reminder-server/internal/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		repo: repository.NewUserRepository(db),
	}
}

func (us *UserService) ListAll() ([]models.User, error) {
	users, err := us.repo.FindAll()

	if err != nil {
		return []models.User{}, err
	}

	return users, nil
}

func (us *UserService) Get(id int64) (models.User, error) {
	user, err := us.repo.FindByID(id)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (us *UserService) GetByEmail(email string) (models.User, error) {
	user, err := us.repo.FindByEmail(email)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (us *UserService) Create(request models.UserCreateRequest) (models.User, error) {
	newUser := models.User{
		Email:    request.Email,
		Password: request.Password,
	}

	existingUser, err := us.GetByEmail(newUser.Email)

	if err != nil && err != gorm.ErrRecordNotFound {
		log.Printf("Error getting user by email: %v", err)
		return models.User{}, errors.New(utils.ErrorUserAlreadyExists)
	}

	if existingUser.ID != 0 {
		return models.User{}, errors.New(utils.ErrorUserAlreadyExists)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return models.User{}, errors.New(utils.ErrorFailedToHash)
	}

	newUser.Password = string(hash)

	user, err := us.repo.Create(newUser)

	return user, err
}

func (us *UserService) Login(request models.UserLoginRequest) (models.User, error) {
	user, err := us.GetByEmail(request.Email)

	if err != nil {
		return models.User{}, errors.New(utils.ErrorUserNotFound)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))

	if err != nil {
		return models.User{}, errors.New(utils.ErrorInvalidPassword)
	}

	_, err = utils.GenerateToken(user)

	if err != nil {
		return models.User{}, err
	}

	log.Printf("User %v logged in", user.ID)

	return user, nil
}
