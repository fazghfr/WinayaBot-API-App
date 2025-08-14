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

// Implement auto-register on non-existant case
func (t *TaskService) initRegisterUser(dto dto.UserDTO) (*model.User, error) {
	User := &model.User{
		ID:         uuid.New().String(),
		Discord_id: dto.DiscordID,
	}
	var err error
	User, err = t.userRepo.CreateNewUser(User)
	if err != nil {
		return nil, err
	}
	return User, nil
}

// autoregister using a function within this file
func (t *TaskService) CreateTask(dto_task *dto.TaskDTO) (*model.Task, error) {
	/*
	 1. checking if user exists in database (precautionary ops)
	 2. checking if title is not empty
	 3. checking if status is not empty -> if it is, set it to a default value
	 4. calling the repository function
	*/
	User, is_exist, err := t.userRepo.CheckUserByDiscordID(dto_task.DiscordID)
	if err != nil {
		return &model.Task{}, err
	}
	if !is_exist {
		userDTO := dto.UserDTO{
			DiscordID: dto_task.DiscordID,
		}
		User, err = t.initRegisterUser(userDTO)
		if err != nil {
			return &model.Task{}, err
		}
	}

	if dto_task.Title == "" {
		return &model.Task{}, errors.New("Title can not be empty!")
	}
	if dto_task.Status == "" {
		dto_task.Status = "backlog"
	}

	if !isValidStatus(dto_task.Status) {
		return &model.Task{}, errors.New("Invalid status. options : done, in progress, backlog")
	}

	// constructing the model
	var Task = &model.Task{
		ID:     uuid.New().String(),
		Title:  dto_task.Title,
		Status: dto_task.Status,
		UserID: User.ID,
		User:   *User,
	}
	Task, err = t.taskRepo.CreateTask(Task)
	if err != nil {
		return &model.Task{}, err
	}

	return Task, nil
}

func (t *TaskService) GetTasksByUser(discordID string, page, limit int) (*dto.PaginatedTasksDTO, error) {
	// Check if user exists
	user, exists, err := t.userRepo.CheckUserByDiscordID(discordID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("User not found")
	}

	// Calculate offset
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Get tasks from repository
	tasks, err := t.taskRepo.GetTasksByUserID(user.ID, offset, limit)
	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := t.taskRepo.CountTasksByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	taskDTOs := make([]dto.TaskResponseDTO, len(tasks))
	for i, task := range tasks {
		taskDTOs[i] = dto.TaskResponseDTO{
			ID:     task.ID,
			Title:  task.Title,
			Status: task.Status,
		}
	}

	// Calculate total pages
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &dto.PaginatedTasksDTO{
		Tasks:      taskDTOs,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (t *TaskService) EditTaskByID(taskDTO *dto.TaskDTO, taskID string) (*model.Task, error) {
	// no need to validate dto, handled in task handler (if implemented)

	// check if given id is in user's possession
	User, exists, err := t.userRepo.CheckUserByDiscordID(taskDTO.DiscordID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("User not found")
	}

	var existingTask *model.Task
	existingTask, exists, err = t.taskRepo.IsTaskExistAndAuthorized(taskID, User.ID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("Task does not exist for given user")
	}

	// only allow these two attributes to be changed
	// if no title or status is set, then keep it as is
	if taskDTO.Title != "" {
		existingTask.Title = taskDTO.Title
	}
	if taskDTO.Status != "" {
		existingTask.Status = taskDTO.Status
	}

	updatedTask, err := t.taskRepo.EditTaskByID(existingTask)
	if err != nil {
		return nil, err
	}

	return updatedTask, nil
}

func (t *TaskService) DeleteTaskByID(taskID string, DiscordID string) (bool, error) {
	// check for auth
	User, exists, err := t.userRepo.CheckUserByDiscordID(DiscordID)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, errors.New("User not found")
	}

	_, exists, err = t.taskRepo.IsTaskExistAndAuthorized(taskID, User.ID)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, errors.New("Task does not exist for given user")
	}

	// attempted deletion
	err = t.taskRepo.DeleteTaskByID(taskID)
	if err != nil {
		return false, errors.New("Task Delete Failed")
	}
	return true, nil
}
