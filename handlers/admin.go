package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AdminHandler struct{}

func (h *AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Admin Handler")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("lalala")
}
