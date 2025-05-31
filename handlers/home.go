package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HomeHandler struct{}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Home Handler")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("lalala")
}
