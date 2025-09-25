package users

import (
	controller "github.com/hafizhproject45/Golang-Boilerplate.git/internal/modules/users/controllers"
	user "github.com/hafizhproject45/Golang-Boilerplate.git/internal/modules/users/services"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(v1 fiber.Router, s user.UserService) {
	ctrl := controller.NewUserController(s)

	route := v1.Group("/users")

	route.Get("/", ctrl.GetAll)
	route.Post("/", ctrl.CreateOne)
	route.Get("/:id", ctrl.GetOne)
	route.Patch("/:id", ctrl.UpdateOne)
	route.Delete("/:id", ctrl.DeleteOne)
}
