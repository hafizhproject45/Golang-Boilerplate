package route

import (
	"github.com/hafizhproject45/Golang-Boilerplate.git/internal/modules"
	"github.com/hafizhproject45/Golang-Boilerplate.git/internal/validation"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	users "github.com/hafizhproject45/Golang-Boilerplate.git/internal/modules/users"
	// MODULE IMPORTS
)

func Routes(app *fiber.App, db *gorm.DB) {
	validate := validation.Validator()
	api := app.Group("/api")

	// masterRoute.Routes(api, db)

	// root modules di sini
	allModules := []modules.Module{
		users.UserModule{},
		// MODULE REGISTRY
	}

	// daftarkan root modules
	for _, m := range allModules {
		m.RegisterRoutes(api, db, validate)
	}

}
