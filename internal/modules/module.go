package modules

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Module interface {
	RegisterRoutes(router fiber.Router, db *gorm.DB, validate *validator.Validate)
}
