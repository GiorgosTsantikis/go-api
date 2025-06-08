package handlers

import (
	"api/service"
	"encoding/json"
	"net/http"
	"strings"
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

// user/profile/{email}
func (u *userHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	if len(path) < 4 || path[len(path)-1] == "" {
		http.Error(w, "Email not provided", http.StatusBadRequest)
		return
	}
	email := path[len(path)-1]
	model, err := u.UserService.GetUserByEmail(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(model)
}
