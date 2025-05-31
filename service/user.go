package service

import (
	"api/converter"
	"api/internal/database"
	"api/model"
	"context"
	"database/sql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetUser(username string) (*model.User, error)
	CreateUser(dto *model.RegistrationDTO) (*model.User, error)
}

type userService struct { /*dependencies*/
	DB *database.Queries
}

func (u *userService) GetUser(username string) (*model.User, error) {
	val, err := u.DB.GetUserByUserName(context.Background(), username)
	if err != nil {
		return nil, err
	}
	return converter.UserEntityToUserModel(&val), nil
}

func (u *userService) CreateUser(dto *model.RegistrationDTO) (*model.User, error) {
	entity, err := u.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:       uuid.New(),
		Username: dto.Username,
		Profilepic: sql.NullString{
			String: "",
			Valid:  true,
		},
	})
	if err != nil {
		return nil, err
	}
	var pass []byte
	pass, err = bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	var cred database.Credential
	cred, err = u.DB.CreateCredentials(context.Background(), database.CreateCredentialsParams{
		ID:       entity.ID,
		Username: dto.Username,
		Password: string(pass),
		Email:    dto.Email,
	})
	if err != nil {
		return nil, err
	}
	return &model.User{
		Username:   cred.Username,
		ProfilePic: "",
	}, nil
}

func NewUserService(queries *database.Queries) UserService {
	return &userService{DB: queries}
}
