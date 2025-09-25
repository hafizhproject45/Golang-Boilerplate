package controller

import (
	"math"
	"strconv"

	"github.com/hafizhproject45/Golang-Boilerplate.git/internal/modules/users/dto"
	service "github.com/hafizhproject45/Golang-Boilerplate.git/internal/modules/users/services"
	validation "github.com/hafizhproject45/Golang-Boilerplate.git/internal/modules/users/validations"
	"github.com/hafizhproject45/Golang-Boilerplate.git/internal/response"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

func (u *UserController) GetAll(c *fiber.Ctx) error {
	query := &validation.Query{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
		Search: c.Query("search", ""),
	}

	result, totalResults, err := u.UserService.GetAll(c, query)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.SuccessWithPaginate[dto.UserListDTO]{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Get all users successfully",
			Meta: response.Meta{
				Page:         query.Page,
				Limit:        query.Limit,
				TotalPages:   int64(math.Ceil(float64(totalResults) / float64(query.Limit))),
				TotalResults: totalResults,
			},
			Data: dto.ToUserListDTOs(result),
		})
}

func (u *UserController) GetOne(c *fiber.Ctx) error {
	param := c.Params("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Id")
	}

	result, err := u.UserService.GetOne(c, uint(id))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Success{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Get user successfully",
			Data:    dto.ToUserListDTO(*result),
		})
}

func (u *UserController) CreateOne(c *fiber.Ctx) error {
	req := new(validation.Create)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := u.UserService.CreateOne(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).
		JSON(response.Success{
			Code:    fiber.StatusCreated,
			Status:  "success",
			Message: "Create user successfully",
			Data:    dto.ToUserListDTO(*result),
		})
}

func (u *UserController) UpdateOne(c *fiber.Ctx) error {
	req := new(validation.Update)
	param := c.Params("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Id")
	}

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := u.UserService.UpdateOne(c, req, uint(id))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Success{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Update user successfully",
			Data:    dto.ToUserListDTO(*result),
		})
}

func (u *UserController) DeleteOne(c *fiber.Ctx) error {
	param := c.Params("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Id")
	}

	if err := u.UserService.DeleteOne(c, uint(id)); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Delete user successfully",
		})
}
