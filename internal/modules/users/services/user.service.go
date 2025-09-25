package service

import (
	"errors"

	model "github.com/hafizhproject45/Golang-Boilerplate.git/internal/modules/users/models"
	repository "github.com/hafizhproject45/Golang-Boilerplate.git/internal/modules/users/repositories"
	validation "github.com/hafizhproject45/Golang-Boilerplate.git/internal/modules/users/validations"
	"github.com/hafizhproject45/Golang-Boilerplate.git/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserService interface {
	GetAll(ctx *fiber.Ctx, params *validation.Query) ([]model.User, int64, error)
	GetOne(ctx *fiber.Ctx, id uint) (*model.User, error)
	CreateOne(ctx *fiber.Ctx, req *validation.Create) (*model.User, error)
	UpdateOne(ctx *fiber.Ctx, req *validation.Update, id uint) (*model.User, error)
	DeleteOne(ctx *fiber.Ctx, id uint) error
}

type userService struct {
	Log        *logrus.Logger
	Validate   *validator.Validate
	Repository repository.UserRepository
}

func NewUserService(repo repository.UserRepository, validate *validator.Validate) UserService {
	return &userService{
		Log:        utils.Log,
		Validate:   validate,
		Repository: repo,
	}
}
func (s userService) GetAll(c *fiber.Ctx, params *validation.Query) ([]model.User, int64, error) {
	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit

	users, total, err := s.Repository.GetAll(c.Context(), offset, params.Limit, func(db *gorm.DB) *gorm.DB {
		if params.Search != "" {
			return db.Where("name LIKE ?", "%"+params.Search+"%")
		}
		return db.Order("created_at DESC").Order("updated_at DESC")
	})

	if err != nil {
		s.Log.Errorf("Failed to get users: %+v", err)
		return nil, 0, err
	}
	return users, total, nil
}

func (s userService) GetOne(c *fiber.Ctx, id uint) (*model.User, error) {
	user, err := s.Repository.GetByID(c.Context(), id, nil)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}
	if err != nil {
		s.Log.Errorf("Failed get user by id: %+v", err)
		return nil, err
	}
	return user, nil
}

func (s *userService) CreateOne(c *fiber.Ctx, req *validation.Create) (*model.User, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	createBody := &model.User{
		Name: req.Name,
	}

	if err := s.Repository.CreateOne(c.Context(), createBody, nil); err != nil {
		s.Log.Errorf("Failed to create user: %+v", err)
		return nil, err
	}

	return createBody, nil
}

func (s userService) UpdateOne(c *fiber.Ctx, req *validation.Update, id uint) (*model.User, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	updateBody := make(map[string]any)

	if req.Name != nil {
		updateBody["name"] = *req.Name
	}

	if err := s.Repository.PatchOne(c.Context(), id, updateBody, nil); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		s.Log.Errorf("Failed to update user: %+v", err)
		return nil, err
	}

	return s.GetOne(c, id)
}

func (s userService) DeleteOne(c *fiber.Ctx, id uint) error {
	if err := s.Repository.DeleteOne(c.Context(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		s.Log.Errorf("Failed to delete user: %+v", err)
		return err
	}
	return nil
}
