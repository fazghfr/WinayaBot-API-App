package service

import (
	"Discord_API_DB_v1/internal/dto"
	"Discord_API_DB_v1/internal/model"
	"Discord_API_DB_v1/internal/repository"
	"errors"
	"github.com/google/uuid"
)

type TaskService struct {
	taskRepo *repository.TaskRepo
	userRepo *repository.UserRepo
}

var availableStatus = map[string]bool{
	"done":        true,
	"in progress": true,
	"backlog":     true,
}

func isValidStatus(s string) bool {
	return availableStatus[s]
}

func InitTaskService(taskRepo *repository.TaskRepo, userRepo *repository.UserRepo) *TaskService {
	return &TaskService{taskRepo: taskRepo, userRepo: userRepo}
}

func (t *TaskService) CreateTask(dto dto.TaskDTO) (*model.Task, error) {
	/*
	 1. checking if user exists in database (precautionary ops)
	 2. checking if title is not empty
	 3. checking if status is not empty -> if it is, set it to a default value
	 4. calling the repository function
	*/
	User, is_exist, err := t.userRepo.CheckUserByDiscordID(dto.DiscordID)
	if err != nil {
		return &model.Task{}, err
	}
	if !is_exist {
		return &model.Task{}, errors.New("User does not exist. Please Register First")
	}

	if dto.Title == "" {
		return &model.Task{}, errors.New("Title can not be empty!")
	}
	if dto.Status == "" {
		dto.Status = "backlog"
	}

	if !isValidStatus(dto.Status) {
		return &model.Task{}, errors.New("Invalid status. options : done, in progress, backlog")
	}

	// constructing the model
	var Task = &model.Task{
		ID:     uuid.New().String(),
		Title:  dto.Title,
		Status: dto.Status,
		UserID: User.ID,
		User:   *User,
	}
	Task, err = t.taskRepo.CreateTask(Task)
	if err != nil {
		return &model.Task{}, err
	}

	return Task, nil
}
