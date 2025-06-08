package converter

import (
	"api/internal/database"
	"api/model"
)

// perfomante maxim
func UserEntityToUserModel(user *database.User) *model.User {
	a := user.Image
	pic := ""
	if a.Valid {
		pic = a.String
	}
	return &model.User{
		Username: user.Name,
		Email:    user.Email,
		Pic:      pic,
		UserId:   user.ID,
	}
}
