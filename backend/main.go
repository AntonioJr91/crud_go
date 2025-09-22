package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func main() {
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASS", "123")
	os.Setenv("DB_HOST", "127.0.0.1")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_NAME", "items")

	InitDB()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
    AllowedOrigins:   []string{"http://localhost:3000",},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
    AllowCredentials: true,
    MaxAge: 300,
}))

	r.Route("/items", func(r chi.Router) {
		r.Get("/", ListItemsHandler)     
		r.Post("/", CreateItemHandler)    
		r.Get("/{id}", GetItemHandler)   
		r.Put("/{id}", UpdateItemHandler) 
		r.Delete("/{id}", DeleteItemHandler) 
	})

	log.Println("Backend server running on http://localhost:8080/items")
	log.Println("Frontend server running on http://localhost:3000")
	http.ListenAndServe("0.0.0.0:8080", r)

}

func connectWithRetry() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	for i := 1; i <= 5; i++ {
		DB, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("Failed to open connection: %v", err)
		} else {
			err = DB.Ping()
			if err == nil {
				log.Println("Database conenection established!")
				return
			}
			log.Printf("Failed to connect (attempt %d): %v", i, err)
		}
		time.Sleep(2 * time.Second)
	}

	log.Fatal("Unable to connect to database after 5 attempts.")
}