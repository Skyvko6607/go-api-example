package services

import (
	"TestAPI/models"
	"TestAPI/repositories"
)

type UserService struct {
	Repo *repositories.UserRepository
}

func (s *UserService) GetUser(userNameOrEmail string) (models.UserDTO, error) {
	user, err := s.Repo.FindByUserNameOrEmail(userNameOrEmail, userNameOrEmail)
	if err != nil {
		return models.UserDTO{}, err
	}

	return user.AsDTO(), err
}

func (s *UserService) CreateUser(userName string, email string) (models.UserDTO, error) {
	user, err := s.Repo.CreateUser(userName, email)
	if err != nil {
		return models.UserDTO{}, err
	}

	return user.AsDTO(), err
}
