package repository

import (
	"reminder-server/internal/models"

	"gorm.io/gorm"
)

type reminderRepository interface {
	FindAll() ([]models.Reminder, error)
	FindByID(id int64) (models.Reminder, error)
	FindByUserID(userID int64) ([]models.Reminder, error)
	Create(reminder models.Reminder) (models.Reminder, error)
	Update(reminder models.Reminder) (models.Reminder, error)
	Delete(id int64) error
}

type ReminderRepository struct {
	db *gorm.DB
}

func NewReminderRepository(db *gorm.DB) ReminderRepository {
	return ReminderRepository{
		db: db,
	}
}

func (rr *ReminderRepository) GetDB() *gorm.DB {
	return rr.db
}

func (rr *ReminderRepository) FindAll() ([]models.Reminder, error) {
	var reminders []models.Reminder
	result := rr.db.Find(&reminders)

	return reminders, result.Error
}

func (rr *ReminderRepository) FindByID(id int64) (models.Reminder, error) {
	var reminder models.Reminder
	result := rr.db.First(&reminder, id)

	return reminder, result.Error
}

func (rr *ReminderRepository) FindByUserID(userID int64) ([]models.Reminder, error) {
	var reminders []models.Reminder
	result := rr.db.Where("user_id = ?", userID).Order(
		"due_date ASC",
	).Find(&reminders)

	return reminders, result.Error
}

func (rr *ReminderRepository) Create(reminder models.Reminder) (models.Reminder, error) {
	result := rr.db.Create(&reminder)

	return reminder, result.Error
}

func (rr *ReminderRepository) Update(reminder models.Reminder) (models.Reminder, error) {
	result := rr.db.Save(&reminder)

	return reminder, result.Error
}

func (rr *ReminderRepository) Delete(id int64) error {
	result := rr.db.Delete(&models.Reminder{}, id)

	return result.Error
}
