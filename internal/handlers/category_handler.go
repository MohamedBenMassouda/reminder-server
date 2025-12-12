package handlers

import (
	"net/http"
	"reminder-server/internal/models"
	"reminder-server/internal/services"
	"reminder-server/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// Create godoc
// @Summary      Create a new category
// @Description  Create a new reminder category
// @Tags         categories
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        category  body      models.CategoryCreateRequest  true  "Category data"
// @Success      201       {object}  models.Category
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /categories/ [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	var req models.CategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.categoryService.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// List godoc
// @Summary      List all categories
// @Description  Get all categories for the authenticated user
// @Tags         categories
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {array}   models.Category
// @Failure      500  {object}  map[string]string
// @Router       /categories/ [get]
func (h *CategoryHandler) List(c *gin.Context) {
	categories, err := h.categoryService.List(1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// Get godoc
// @Summary      Get a category by ID
// @Description  Get category details by ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  models.Category
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /categories/{id} [get]
func (h *CategoryHandler) Get(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.categoryService.Get(int64(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrorSqlNoRows(err).Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// Update godoc
// @Summary      Update a category
// @Description  Update category details
// @Tags         categories
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id        path      int                           true  "Category ID"
// @Param        category  body      models.CategoryUpdateRequest  true  "Category update data"
// @Success      200       {object}  models.Category
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /categories/{id} [patch]
func (h *CategoryHandler) Update(c *gin.Context) {
	var req models.CategoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	categoryID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.categoryService.Update(int64(categoryID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// Delete godoc
// @Summary      Delete a category
// @Description  Delete a category by ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id   path      int  true  "Category ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.categoryService.Delete(int64(categoryID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
