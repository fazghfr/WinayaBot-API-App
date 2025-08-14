package repository

import (
	"Discord_API_DB_v1/internal/model"
	"errors"
	"gorm.io/gorm"
)

type TaskRepo struct {
	db *gorm.DB
}

func InitTaskRepository(db *gorm.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) CreateTask(task *model.Task) (*model.Task, error) {
	err := r.db.Create(&task).Error
	if err != nil {
		return &model.Task{}, err
	}

	return task, nil
}

func (r *TaskRepo) GetTasksByUserID(userID string, offset, limit int) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepo) CountTasksByUserID(userID string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Task{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// only gets called by an authorized entity
func (r *TaskRepo) IsTaskExistAndAuthorized(id string, userID string) (*model.Task, bool, error) {
	var task model.Task
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&task).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil // Not found, but no DB error
		}
		return nil, false, err // Some other DB error
	}

	return &task, true, nil
}

// ONLY GETS CALLED FOR AN EXISTING MODEL. MAKE SURE TO CHECK TASK EXISTANCE. USING IsTaskExist function
func (r *TaskRepo) EditTaskByID(existingTask *model.Task) (*model.Task, error) {
	if err := r.db.Save(existingTask).Error; err != nil {
		return nil, err
	}
	return existingTask, nil
}

func (r *TaskRepo) DeleteTaskByID(taskID string) error {
	if err := r.db.Delete(&model.Task{}, "id = ?", taskID).Error; err != nil {
		return err
	}
	return nil
}
