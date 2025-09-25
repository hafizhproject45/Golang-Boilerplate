package repository

import (
	model "github.com/hafizhproject45/Golang-Boilerplate.git/internal/modules/users/models"
	"github.com/hafizhproject45/Golang-Boilerplate.git/internal/repository"

	"gorm.io/gorm"
)

type UserRepository interface {
	repository.BaseRepository[model.User]
}

type UserRepositoryImpl struct {
	*repository.BaseRepositoryImpl[model.User]
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		BaseRepositoryImpl: repository.NewBaseRepository[model.User](db),
	}
}
