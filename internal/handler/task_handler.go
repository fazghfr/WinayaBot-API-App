package handler

import (
	"Discord_API_DB_v1/internal/dto"
	"Discord_API_DB_v1/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	s *service.TaskService
}

func InitTaskHandler(s *service.TaskService) *TaskHandler {
	return &TaskHandler{s: s}
}

func (t *TaskHandler) CreateNewTask(c *gin.Context) {
	var taskDTO dto.TaskDTO
	if err := c.ShouldBindJSON(&taskDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	Task, err := t.s.CreateTask(taskDTO)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, *Task)
}

func (t *TaskHandler) GetTasksByUser(c *gin.Context) {
	// Get query parameters
	discordID := c.Query("discord_id")
	if discordID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "discord_id is required",
		})
		return
	}

	page := 1
	if pageParam := c.Query("page"); pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	limit := 10
	if limitParam := c.Query("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
			limit = l
		}
	}

	// Get paginated tasks
	paginatedTasks, err := t.s.GetTasksByUser(discordID, page, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, paginatedTasks)
}
