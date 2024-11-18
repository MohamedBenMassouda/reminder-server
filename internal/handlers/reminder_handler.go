package handlers

import (
	"net/http"
	"reminder-server/internal/models"
	"reminder-server/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReminderHandler struct {
	service *services.ReminderService
}

func NewReminderHandler(service *services.ReminderService) *ReminderHandler {
	return &ReminderHandler{
		service: service,
	}
}

func (h *ReminderHandler) List(c *gin.Context) {
	// TODO: Take the user id from the authenticated user

	reminders, err := h.service.List(1)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, reminders)
}

func (h *ReminderHandler) Get(c *gin.Context) {
	reminderID, err := strconv.Atoi(c.Param("id"))

	reminder, err := h.service.Get(int64(reminderID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reminder)
}

func (h *ReminderHandler) Create(c *gin.Context) {
	var req models.ReminderCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reminder, err := h.service.Create(req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reminder)
}

func (h *ReminderHandler) Update(c *gin.Context) {
	reminderID, err := strconv.Atoi(c.Param("id"))

	var req models.ReminderUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reminder, err := h.service.Update(int64(reminderID), req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reminder)
}

func (h *ReminderHandler) Delete(c *gin.Context) {
	reminderID, err := strconv.Atoi(c.Param("id"))

	err = h.service.Delete(int64(reminderID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (h *ReminderHandler) UpdateStatus(c *gin.Context) {
	reminderID, err := strconv.Atoi(c.Param("id"))

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reminder, err := h.service.UpdateStatus(int64(reminderID), req.Status)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reminder)
}
