package repository

import (
	"Discord_API_DB_v1/internal/model"
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
