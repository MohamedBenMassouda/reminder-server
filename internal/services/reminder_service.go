package services

import (
	"errors"
	"reminder-server/internal/models"
	"reminder-server/internal/repository"
	"reminder-server/internal/utils"

	"gorm.io/gorm"
)

type ReminderService struct {
	repo repository.ReminderRepository
}

func NewReminderService(db *gorm.DB) *ReminderService {
	return &ReminderService{
		repo: repository.NewReminderRepository(db),
	}
}

func (rs *ReminderService) ListAll() ([]models.Reminder, error) {
	reminders, err := rs.repo.FindAll()

	if err != nil {
		return []models.Reminder{}, err
	}

	return reminders, nil
}

func (rs *ReminderService) List(userID int32) ([]models.Reminder, error) {
	reminders, err := rs.repo.FindByUserID(int64(userID))

	if err != nil {
		return []models.Reminder{}, err
	}

	return reminders, nil
}

func (rs *ReminderService) Get(id int64) (models.Reminder, error) {
	reminder, err := rs.repo.FindByID(id)

	if err != nil {
		return models.Reminder{}, err
	}

	return reminder, nil
}

func (rs *ReminderService) Create(request models.ReminderCreateRequest) (models.Reminder, error) {
	// TODO: Make it use a real user ID
	newReminder := models.Reminder{
		Title:            request.Title,
		Description:      request.Description,
		DueDate:          request.DueDate.String(),
		CategoryID:       request.CategoryID,
		IsRecurring:      request.IsRecurring,
		RecurringPattern: request.RecurringPattern,
		Priority:         request.Priority,
		Status:           models.StatusPending,
		UserID:           1,
	}

	// Check if the category exists
	category, err := NewCategoryService(rs.repo.GetDB()).Get(request.CategoryID)

	if err != nil {
		return models.Reminder{}, err
	}

	if category.ID == 0 {
		return models.Reminder{}, errors.New(utils.ErrorCategoryNotFound)
	}

	reminder, err := rs.repo.Create(newReminder)

	return reminder, err
}

func (rs *ReminderService) Update(id int64, request models.ReminderUpdateRequest) (models.Reminder, error) {
	// Get the reminder
	reminder, err := rs.Get(id)

	if err != nil {
		return models.Reminder{}, err
	}

	if request.Title != nil {
		reminder.Title = *request.Title
	}

	if request.Description != nil {
		reminder.Description = *request.Description
	}

	if request.CategoryID != nil {
		reminder.CategoryID = *request.CategoryID
	}

	if request.DueDate != nil {
		reminder.DueDate = request.DueDate.String()
	}

	if request.Priority != nil {
		reminder.Priority = *request.Priority
	}

	if request.Status != nil {
		reminder.Status = *request.Status
	}

	if request.IsRecurring != nil {
		reminder.IsRecurring = *request.IsRecurring
	}

	if request.RecurringPattern != nil {
		reminder.RecurringPattern = *request.RecurringPattern
	}

	updatedReminder, err := rs.repo.Update(reminder)

	return updatedReminder, err
}

func (rs *ReminderService) Delete(id int64) error {
	err := rs.repo.Delete(id)

	return err
}
