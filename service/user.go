package service

import (
	"api/converter"
	"api/internal/database"
	"api/model"
	"context"
	"fmt"
)

type UserService interface {
	GetUserByID(id string) (*model.User, error)
}

type userService struct { /*dependencies*/
	DB *database.Queries
}

func (u *userService) GetUserByID(id string) (*model.User, error) {
	fmt.Printf("UserService.GetUserByID %v", id)
	user, err := u.DB.GetUserByID(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return converter.UserEntityToUserModel(&user), nil
}

func NewUserService(queries *database.Queries) UserService {
	return &userService{DB: queries}
}
