package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

type HealthResponse struct {
	Status string `json:"status"`
}

type ReadinessResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// Liveness godoc
// @Summary      Liveness probe
// @Description  Kubernetes liveness probe endpoint - checks if the application is running
// @Tags         health
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Router       /health/liveness [get]
func (h *HealthHandler) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status: "alive",
	})
}

// Readiness godoc
// @Summary      Readiness probe
// @Description  Kubernetes readiness probe endpoint - checks if the application is ready to serve traffic
// @Tags         health
// @Produce      json
// @Success      200  {object}  ReadinessResponse
// @Failure      503  {object}  ReadinessResponse
// @Router       /health/readiness [get]
func (h *HealthHandler) Readiness(c *gin.Context) {
	// Check database connection
	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, ReadinessResponse{
			Status:   "not ready",
			Database: "disconnected",
		})
		return
	}

	// Ping database
	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, ReadinessResponse{
			Status:   "not ready",
			Database: "unhealthy",
		})
		return
	}

	c.JSON(http.StatusOK, ReadinessResponse{
		Status:   "ready",
		Database: "healthy",
	})
}

// Startup godoc
// @Summary      Startup probe
// @Description  Kubernetes startup probe endpoint - checks if the application has started successfully
// @Tags         health
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Failure      503  {object}  HealthResponse
// @Router       /health/startup [get]
func (h *HealthHandler) Startup(c *gin.Context) {
	// Check database connection
	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, HealthResponse{
			Status: "starting",
		})
		return
	}

	// Ping database
	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, HealthResponse{
			Status: "starting",
		})
		return
	}

	c.JSON(http.StatusOK, HealthResponse{
		Status: "started",
	})
}
