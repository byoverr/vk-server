package handlers

import (
	"backend/internal/models"
	"backend/internal/storage"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type Handler struct {
	DB     *storage.PostgresqlDB
	logger *slog.Logger
}

func NewHandler(db *storage.PostgresqlDB, logger *slog.Logger) *Handler {
	return &Handler{DB: db, logger: logger}
}

func InitRoutes(router *gin.Engine, handler *Handler) {
	router.GET("/containers", handler.GetContainers)
	router.POST("/containers", handler.PostContainer)
}

func (h *Handler) GetContainers(c *gin.Context) {
	h.logger.Info("start get all statuses")

	statuses, err := h.DB.ReadAllStatuses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting all statuses"})
		return
	}

	h.logger.Info("successful to get all statuses")
	c.JSON(http.StatusOK, statuses)
}

func (h *Handler) PostContainer(c *gin.Context) {
	h.logger.Info("start adding container")

	var rawStatus models.RawContainerStatus
	if err := c.ShouldBindJSON(&rawStatus); err != nil {
		h.logger.Error("invalid request data")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}

	status, err := rawStatus.ToContainerStatus()
	if err != nil {
		h.logger.Error("invalid time format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid time format"})
		return
	}

	if err := h.DB.CreateStatus(status); err != nil {
		h.logger.Error("failed to insert container status")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert container status"})
		return
	}

	h.logger.Info("container status successfully added")
	c.JSON(http.StatusCreated, gin.H{"message": "container status added"})
}
