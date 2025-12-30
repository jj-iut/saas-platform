package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	checks := make(map[string]string)

	// Check database connection
	if err := h.db.Ping(); err != nil {
		checks["database"] = "unhealthy"
		c.JSON(http.StatusServiceUnavailable, HealthResponse{
			Status:    "unhealthy",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Checks:    checks,
		})
		return
	}

	checks["database"] = "healthy"

	c.JSON(http.StatusOK, HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    checks,
	})
}
