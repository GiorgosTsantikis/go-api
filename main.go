package main

import (
	"api/handlers"
	"api/internal/auth"
	"api/internal/database"
	"api/service"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbConnection := os.Getenv("DB_URL")
	if dbConnection == "" {
		log.Fatal("DB_URL is not set")
	}
	fmt.Println("DB_URL:", dbConnection)
	conn, err := sql.Open("postgres", dbConnection)
	if err != nil {
		log.Fatal("Error opening database", err)
	}
	db := database.New(conn)

	mux := http.NewServeMux()
	userService := service.NewUserService(db)
	storeService := service.NewStoreService(db, conn)
	middleware := auth.NewMiddleware()

	userHandler := handlers.NewUserHandler(userService)
	storeHandler := handlers.NewStoreHandler(storeService)

	mux.HandleFunc("GET /user/profile", middleware.WithCORS(middleware.WithSecurityHeaders(middleware.AuthenticationMiddleware(userHandler.GetUserProfile))))

	mux.HandleFunc("POST /store/is-valid-url", middleware.CleanXSS(middleware.WithSecurityHeaders(middleware.AuthenticateStoreMiddleware(middleware.WithCORS(middleware.AuthenticationMiddleware(storeHandler.IsValidURL))))))
	mux.HandleFunc("/store/create-menu-item", middleware.WithSecurityHeaders(middleware.WithCORS(middleware.AuthenticationMiddleware(middleware.AuthenticateStoreMiddleware(storeHandler.CreateMenuItem)))))
	mux.HandleFunc("GET /store/generate-qrcode", middleware.WithCORS(middleware.WithSecurityHeaders(middleware.AuthenticateStoreMiddleware(storeHandler.GenerateQRCode))))
	mux.HandleFunc("GET /store/generate-qrcode-pdf", middleware.WithCORS(middleware.WithSecurityHeaders(middleware.AuthenticateStoreMiddleware(storeHandler.GetQRCodePDF))))
	mux.HandleFunc("DELETE /store/delete-menu-item/{id}", middleware.WithCORS(middleware.WithSecurityHeaders(middleware.AuthenticateStoreMiddleware(storeHandler.DeleteMenuItem))))

	adminHandler := handlers.NewAdminHandler(service.NewAdminService(db))
	//TODO protect
	mux.HandleFunc("/admin/users", adminHandler.GetUsers)
	fmt.Println("Listening on port ", os.Getenv("PORT"))

	err = http.ListenAndServe(":"+os.Getenv("PORT"), mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
