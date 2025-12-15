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

// List godoc
// @Summary      List all reminders
// @Description  Get all reminders for the authenticated user
// @Tags         reminders
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {array}   models.Reminder
// @Failure      500  {object}  map[string]string
// @Router       /reminders/ [get]
func (h *ReminderHandler) List(c *gin.Context) {
	// TODO: Take the user id from the authenticated user
	userID := c.GetInt64("user_id")

	if userID == 0 {
		userID = 1 // Temporary hardcoded user ID for testing
	}

	reminders, err := h.service.List(int32(userID))

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, reminders)
}

// Get godoc
// @Summary      Get a reminder by ID
// @Description  Get reminder details by ID
// @Tags         reminders
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id   path      int  true  "Reminder ID"
// @Success      200  {object}  models.Reminder
// @Failure      500  {object}  map[string]string
// @Router       /reminders/{id} [get]
func (h *ReminderHandler) Get(c *gin.Context) {
	reminderID, err := strconv.Atoi(c.Param("id"))

	reminder, err := h.service.Get(int64(reminderID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reminder)
}

// Create godoc
// @Summary      Create a new reminder
// @Description  Create a new reminder
// @Tags         reminders
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        reminder  body      models.ReminderCreateRequest  true  "Reminder data"
// @Success      201       {object}  models.Reminder
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /reminders/ [post]
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

// Update godoc
// @Summary      Update a reminder
// @Description  Update reminder details
// @Tags         reminders
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id        path      int                           true  "Reminder ID"
// @Param        reminder  body      models.ReminderUpdateRequest  true  "Reminder update data"
// @Success      200       {object}  models.Reminder
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /reminders/{id} [patch]
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

// Delete godoc
// @Summary      Delete a reminder
// @Description  Delete a reminder by ID
// @Tags         reminders
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id   path      int  true  "Reminder ID"
// @Success      204  {object}  nil
// @Failure      500  {object}  map[string]string
// @Router       /reminders/{id} [delete]
func (h *ReminderHandler) Delete(c *gin.Context) {
	reminderID, err := strconv.Atoi(c.Param("id"))

	err = h.service.Delete(int64(reminderID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

// UpdateStatus godoc
// @Summary      Update reminder status
// @Description  Update the status of a reminder (pending, completed, cancelled)
// @Tags         reminders
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id      path      int                       true  "Reminder ID"
// @Param        status  body      object{status=string}     true  "Status data"
// @Success      200     {object}  models.Reminder
// @Failure      400     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /reminders/{id}/status [put]
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
