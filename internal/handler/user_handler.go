package handler

import (
	"Discord_API_DB_v1/internal/dto"
	"Discord_API_DB_v1/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	s *service.UserService
}

func InitUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{
		s: s,
	}
}

// handle for initial call (!ping)
func (h *UserHandler) InitRegistration(c *gin.Context) {
	var userDTO dto.UserDTO
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	User, err := h.s.RegisterUser(userDTO)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, *User)
}
