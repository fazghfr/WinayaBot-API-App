package service

import (
	"Discord_API_DB_v1/internal/dto"
	"Discord_API_DB_v1/internal/model"
	"Discord_API_DB_v1/internal/repository"
	"errors"
	"github.com/google/uuid"
)

type UserService struct {
	r *repository.UserRepo
}

func InitUserService(r *repository.UserRepo) *UserService {
	return &UserService{
		r: r,
	}
}

// User insert
func (s *UserService) RegisterUser(dto dto.UserDTO) (*model.User, error) {
	/*
		1. check for existance of a user
		2. if user does not exist, create new record
		3. if not, return error
	*/

	// checking existance block
	User, Is_exist, err := s.r.CheckUserByDiscordID(dto.DiscordID)
	if err != nil {
		return &model.User{}, err
	}
	if Is_exist {
		return User, errors.New("User already registered")
	}

	// create a new user block,  using dto
	User = &model.User{
		ID:         uuid.New().String(),
		Discord_id: dto.DiscordID,
	}

	User, err = s.r.CreateNewUser(User)
	if err != nil {
		return &model.User{}, err
	}

	return User, nil
}
