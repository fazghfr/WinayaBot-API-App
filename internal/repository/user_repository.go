package repository

import (
	"Discord_API_DB_v1/internal/model"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func InitUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db}
}

// Check user existance by User id
func (r *UserRepo) FindUserByID(id string) (*model.User, bool, error) {
	var User model.User
	err := r.db.Where("id = ?", id).First(&User).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Record not found is not an application error in this case.
			// It simply means the user does not exist.
			return &model.User{}, false, nil
		}
		// An actual database error occurred.
		return &model.User{}, false, err
	}

	// User was found.
	return &User, true, nil
}

// Check user existance by Username
func (r *UserRepo) CheckUserByDiscordID(DiscordID string) (*model.User, bool, error) {
	var User model.User
	err := r.db.Where("discord_id = ?", DiscordID).First(&User).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Record not found is not an application error in this case.
			// It simply means the user does not exist.
			return &model.User{}, false, nil
		}
		// An actual database error occurred.
		return &model.User{}, false, err
	}

	// User was found.
	return &User, true, nil
}

// Insert a new user (happens when checking existance of a user returns false)
func (r *UserRepo) CreateNewUser(user *model.User) (*model.User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return &model.User{}, err
	}
	return user, err
}
