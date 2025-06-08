package handlers

import (
	"api/service"
	"encoding/json"
	"fmt"
	"net/http"
)

type adminHandler struct {
	adminService service.AdminService
}

type AdminHandler interface {
	GetUsers(http.ResponseWriter, *http.Request)
}

func NewAdminHandler(adminService service.AdminService) AdminHandler {
	return &adminHandler{
		adminService: adminService,
	}
}

func (h *adminHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Admin Handler")
	users := h.adminService.GetAllUsers()
	if users == nil {
		w.WriteHeader(http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
