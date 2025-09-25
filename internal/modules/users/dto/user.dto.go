package dto

import (
	"time"

	model "github.com/hafizhproject45/Golang-Boilerplate.git/internal/modules/users/models"
)

// === DTO Structs ===

type UserListDTO struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserDetailDTO struct {
	UserListDTO
}

// === Mapper Functions ===

func ToUserListDTO(m model.User) UserListDTO {
	return UserListDTO{
		Id:   m.Id,
		Name: m.Name,
	}
}

func ToUserListDTOs(m []model.User) []UserListDTO {
	result := make([]UserListDTO, len(m))
	for i, r := range m {
		result[i] = ToUserListDTO(r)
	}
	return result
}
