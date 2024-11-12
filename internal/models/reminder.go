package models

import (
	"encoding/json"
	"time"
)

type Reminder struct {
	ID               int64  `json:"id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	CategoryID       int64  `json:"category_id"`
	DueDate          string `json:"due_date"`
	Priority         string `json:"priority"`
	Status           string `json:"status"`
	IsRecurring      bool   `json:"is_recurring"`
	RecurringPattern string `json:"recurring_pattern,omitempty"`
	UserID           int64  `json:"user_id"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

type ReminderResponse struct {
	ID               int64     `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Category         Category  `json:"category"`
	DueDate          time.Time `json:"due_date"`
	Priority         string    `json:"priority"`
	Status           string    `json:"status"`
	IsRecurring      bool      `json:"is_recurring"`
	RecurringPattern string    `json:"recurring_pattern,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type ReminderCreateRequest struct {
	Title            string    `json:"title" binding:"required"`
	Description      string    `json:"description"`
	CategoryID       int64     `json:"category_id" binding:"required"`
	DueDate          time.Time `json:"due_date" binding:"required"`
	Priority         string    `json:"priority" binding:"required,oneof=low medium high"`
	IsRecurring      bool      `json:"is_recurring"`
	RecurringPattern string    `json:"recurring_pattern,omitempty"`
}

// ReminderUpdateRequest for updating existing reminders
type ReminderUpdateRequest struct {
	Title            *string    `json:"title,omitempty"`
	Description      *string    `json:"description,omitempty"`
	CategoryID       *int64     `json:"category_id,omitempty"`
	DueDate          *time.Time `json:"due_date,omitempty"`
	Priority         *string    `json:"priority,omitempty" binding:"omitempty,oneof=low medium high"`
	Status           *string    `json:"status,omitempty" binding:"omitempty,oneof=pending completed"`
	IsRecurring      *bool      `json:"is_recurring,omitempty"`
	RecurringPattern *string    `json:"recurring_pattern,omitempty"`
}

// Constants for reminder status and priority
const (
	StatusPending   = "pending"
	StatusCompleted = "completed"

	PriorityLow    = "low"
	PriorityMedium = "medium"
	PriorityHigh   = "high"
)

func (rc *ReminderCreateRequest) UnmarshalJSON(data []byte) error {
	type Alias ReminderCreateRequest
	aux := &struct {
		DueDate string `json:"due_date"`
		*Alias
	}{
		Alias: (*Alias)(rc),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	dueDate, err := time.Parse(time.DateTime, aux.DueDate)
	if err != nil {
		return err
	}

	rc.DueDate = dueDate

	return nil
}
