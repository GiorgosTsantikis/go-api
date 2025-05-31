package main

import (
	"api/handlers"
	"api/internal/database"
	"api/service"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
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
	mux.Handle("/", &handlers.HomeHandler{})
	mux.Handle("/user/", handlers.NewUserHandler(service.NewUserService(db)))
	fmt.Println("Listening on port ", os.Getenv("PORT"))

	err = http.ListenAndServe(":"+os.Getenv("PORT"), mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
