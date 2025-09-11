package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASS", "123")
	os.Setenv("DB_HOST", "127.0.0.1")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_NAME", "items")

	InitDB()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// config cors
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://127.0.0.1:5500", // live Server
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Rotas
	r.Route("/items", func(r chi.Router) {
		r.Get("/", ListItemsHandler) // Listar todos
		r.Post("/", CreateItemHandler) // Criar
		r.Get("/{id}", GetItemHandler) // Buscar por ID
		r.Put("/{id}", UpdateItemHandler) // Atualizar
		r.Delete("/{id}", DeleteItemHandler)// Deletar
	})

	log.Println("Servidor rodando em http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
