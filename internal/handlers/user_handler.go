package handlers

import (
	"net/http"
	"reminder-server/internal/models"
	"reminder-server/internal/services"
	"reminder-server/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// List godoc
// @Summary      List all users
// @Description  Get all users
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {array}   models.User
// @Failure      500  {object}  map[string]string
// @Router       /users/ [get]
func (h *UserHandler) List(c *gin.Context) {
	users, err := h.userService.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// Get godoc
// @Summary      Get a user by ID
// @Description  Get user details by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  models.User
// @Failure      500  {object}  map[string]string
// @Router       /users/{id} [get]
func (h *UserHandler) Get(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))

	user, err := h.userService.Get(int64(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// SignUp godoc
// @Summary      Register a new user
// @Description  Create a new user account
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      models.UserCreateRequest  true  "User registration data"
// @Success      201   {object}  models.UserResponse
// @Failure      400   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /users/signup [post]
func (h *UserHandler) SignUp(c *gin.Context) {
	var req models.UserCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Create(req)

	if err != nil {
		if err.Error() == utils.ErrorUserAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and return JWT token
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        credentials  body      models.UserLoginRequest  true  "Login credentials"
// @Success      200          {object}  models.UserResponse
// @Failure      400          {object}  map[string]string
// @Failure      401          {object}  map[string]string
// @Failure      404          {object}  map[string]string
// @Failure      500          {object}  map[string]string
// @Router       /users/login [post]
func (h *UserHandler) Login(c *gin.Context) {

	var req models.UserLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Login(req)

	if err != nil {
		if err.Error() == utils.ErrorUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if err.Error() == utils.ErrorInvalidPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
