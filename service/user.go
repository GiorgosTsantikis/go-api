package service

import (
	"api/converter"
	"api/internal/database"
	"api/model"
	"context"
	"fmt"
)

type UserService interface {
	GetUserByEmail(email string) (*model.User, error)
	GetUserByCookie(sessionToken string) (*model.User, error)
}

type userService struct { /*dependencies*/
	DB *database.Queries
}

func (u *userService) GetUserByEmail(email string) (*model.User, error) {
	fmt.Printf("UserService.GetUserByEmail %v", email)
	val, err := u.DB.GetUserByEmail(context.Background(), email)
	if err != nil {
		return nil, err
	}
	return converter.UserEntityToUserModel(&val), nil
}

func (u *userService) GetUserByCookie(sessionToken string) (*model.User, error) {
	fmt.Printf("UserService.GetUserByCookie %v", sessionToken)
	user, err := u.DB.GetUserBySession(context.Background(), sessionToken)
	if err != nil {
		return nil, err
	}
	return converter.UserEntityToUserModel(&user), nil
}

func NewUserService(queries *database.Queries) UserService {
	return &userService{DB: queries}
}
