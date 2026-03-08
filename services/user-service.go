package services

import (
	"github.com/Skyvko6607/go-api-example/models"
	"github.com/Skyvko6607/go-api-example/repositories"
	"golang.org/x/crypto/bcrypt"

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

func (s *UserService) CreateUser(c *gin.Context, userName string, email string, password string) (models.UserDTO, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.UserDTO{}, err
	}
	user, err := s.Repo.CreateUser(c, userName, email, string(hash))
	if err != nil {
		return models.UserDTO{}, err
	}

	return user.AsDTO(), err
}
