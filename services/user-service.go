package services

import (
	"github.com/Skyvko6607/go-api-learning/models"
	"github.com/Skyvko6607/go-api-learning/repositories"

	"github.com/gin-gonic/gin"
)

type UserService struct {
	Repo *repositories.UserRepository
}

func (s *UserService) GetUser(c *gin.Context, userNameOrEmail string) (models.UserDTO, error) {
	user, err := s.Repo.FindByUserNameOrEmail(c, userNameOrEmail, userNameOrEmail)
	if err != nil {
		return models.UserDTO{}, err
	}

	return user.AsDTO(), err
}

func (s *UserService) CreateUser(c *gin.Context, userName string, email string) (models.UserDTO, error) {
	user, err := s.Repo.CreateUser(c, userName, email)
	if err != nil {
		return models.UserDTO{}, err
	}

	return user.AsDTO(), err
}
