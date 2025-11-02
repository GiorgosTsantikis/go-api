package handlers

import (
	"api/model"
	"api/service"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

type StoreHandler interface {
	IsValidURL(http.ResponseWriter, *http.Request)
	GenerateQRCode(w http.ResponseWriter, r *http.Request)
	GetQRCodePDF(w http.ResponseWriter, r *http.Request)
	DeleteMenuItem(w http.ResponseWriter, r *http.Request)
	CreateMenuItem(w http.ResponseWriter, r *http.Request)
}

type storeHandler struct {
	StoreService service.StoreService
}

func NewStoreHandler(storeService service.StoreService) StoreHandler {
	return &storeHandler{StoreService: storeService}
}

// user/profile
func (s *storeHandler) IsValidURL(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("true"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte("false"))
	w.WriteHeader(http.StatusOK)
}

// /store/generate-qrcode
func (s *storeHandler) GenerateQRCode(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hit generateQRCode")
	storeId := r.Context().Value("store-id")
	fmt.Printf("store-id %v", storeId.(int))
	qrCodes, err := s.StoreService.GenerateQRCodes(storeId.(int))
	if err != nil {
		fmt.Printf("failed to generate QR code %v", err)
		fmt.Println("fcked")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	base64Encoded := make([]string, len(qrCodes))
	for i, data := range qrCodes {
		base64Encoded[i] = base64.StdEncoding.EncodeToString(data)
	}

	w.Header().Set("Content-Type", "application/json")
	jsonErr := json.NewEncoder(w).Encode(base64Encoded)
	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// /store/generate-qr-pdf
func (s *storeHandler) GetQRCodePDF(w http.ResponseWriter, r *http.Request) {
	storeId := r.Context().Value("store-id")
	fmt.Printf("GenerateQRCodePDF store-id %v", storeId.(int))
	buffer, err := s.StoreService.GenerateQRCodePDF(storeId.(int))
	if err != nil {
		w.Write([]byte("Error generating pdf"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=\"qrcodes.pdf\"")
	_, errOut := w.Write(buffer.Bytes())
	if errOut != nil {
		w.Write([]byte("Error attaching pdf"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *storeHandler) DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("delete-menu-item", r.PathValue("id"))
	w.WriteHeader(http.StatusNoContent)
}

func (s *storeHandler) CreateMenuItem(w http.ResponseWriter, r *http.Request) {
	var menuItem model.MenuItem
	storeId := r.Context().Value("store-id")

	if err := json.NewDecoder(r.Body).Decode(&menuItem); err != nil {
		http.Error(w, "Invalid json: "+err.Error(), http.StatusBadRequest)
		return
	}

	id, err := s.StoreService.CreateMenuItem(int32(storeId.(int)), &menuItem)
	if err != nil {
		http.Error(w, "Error creating menu item "+err.Error(), http.StatusInternalServerError)
		return
	}
	type CreateMenuItemResponse struct {
		ID int64 `json:"id"`
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(CreateMenuItemResponse{
		ID: id,
	}); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}
