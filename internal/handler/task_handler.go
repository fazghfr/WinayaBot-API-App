package handler

import (
	"Discord_API_DB_v1/internal/dto"
	"Discord_API_DB_v1/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
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
