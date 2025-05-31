package converter

import (
	"api/internal/database"
	"api/model"
)

// perfomante maxim
func UserEntityToUserModel(user *database.User) *model.User {
	a := user.Profilepic
	pic := ""
	if a.Valid {
		pic = a.String
	}
	return &model.User{
		Username:   user.Username,
		ProfilePic: pic,
	}
}
