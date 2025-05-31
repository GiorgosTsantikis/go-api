package handlers

import (
	"api/model"
	"api/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type UserHandler struct {
	UserService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

//url length/method ? eg POST /profile GET profile/abc GET friends
func (u *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	w.Header().Set("Content-Type", "application/json")
	fmt.Printf("path %v ", path)
	if r.Method == "GET" {
		switch path[len(path)-2] {
		case "profile":
			u.getUserProfile(w, r)
		case "/friends":

		default:
			http.NotFound(w, r)
		}
	} else if r.Method == "POST" {
		switch path[len(path)-1] {
		case "signup":
			u.createUser(w, r)
		case "login":

		case "/friend":

		case "/profile":

		default:
			http.NotFound(w, r)
		}
	}
}

func (u *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var user model.RegistrationDTO
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("user data: ", user)
	a, err := u.UserService.CreateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (u *UserHandler) getUserProfile(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	if len(path) < 4 || path[len(path)-1] == "" {
		http.Error(w, "Username not provided", http.StatusBadRequest)
		return
	}
	username := path[len(path)-1]
	model, err := u.UserService.GetUser(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(model)
}
