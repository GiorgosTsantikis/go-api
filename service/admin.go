package service

import (
	"api/converter"
	"api/internal/database"
	"api/model"
	"context"
)

func NewAdminService(db *database.Queries) AdminService {
	return &adminService{
		DB: db,
	}
}

type AdminService interface {
	GetAllUsers() []*model.User
}

type adminService struct {
	DB *database.Queries
}

func (s *adminService) GetAllUsers() []*model.User {
	users, err := s.DB.GetAllUsers(context.Background())
	if err != nil || len(users) == 0 {
		return nil
	}
	result := make([]*model.User, len(users))
	for idx, user := range users {
		result[idx] = converter.UserEntityToUserModel(&user)
	}
	return result
}
