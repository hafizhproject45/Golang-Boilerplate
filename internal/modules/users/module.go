package users

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	rUser "github.com/hafizhproject45/Golang-Boilerplate.git/internal/modules/users/repositories"
	sUser "github.com/hafizhproject45/Golang-Boilerplate.git/internal/modules/users/services"
)

type UserModule struct{}

func (UserModule) RegisterRoutes(router fiber.Router, db *gorm.DB, validate *validator.Validate) {
	userRepo := rUser.NewUserRepository(db)

	userService := sUser.NewUserService(userRepo, validate)

	UserRoutes(router, userService)
}
