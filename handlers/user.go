package handlers

import (
	"api/service"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserHandler interface {
	UserExistsByEmail(http.ResponseWriter, *http.Request)
	GetUserProfile(http.ResponseWriter, *http.Request)
}

type userHandler struct {
	UserService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{UserService: userService}
}

// user/{email}
func (u *userHandler) UserExistsByEmail(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// user/profile
func (u *userHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userIdVal := r.Context().Value("user-id")
	userId := userIdVal.(string)

	fmt.Printf("UserHandler.GetUserProfile id:%v \n", userId)

	user, err := u.UserService.GetUserByID(userId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User with " + userId + " was not found"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	encodingError := json.NewEncoder(w).Encode(user)
	if encodingError != nil {
		fmt.Println("Failed to encode user object")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
