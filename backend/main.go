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
	os.Setenv("DB_HOST", "mariadb")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_NAME", "items")

	// Inicializa a conexão com o banco
	InitDB()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// config cors
	r.Use(cors.Handler(cors.Options{
    AllowedOrigins:   []string{"http://localhost:3000",},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
    AllowCredentials: true,
    MaxAge: 300,
}))

	// Rotas
	r.Route("/items", func(r chi.Router) {
		r.Get("/", ListItemsHandler)      // Listar todos
		r.Post("/", CreateItemHandler)    // Criar
		r.Get("/{id}", GetItemHandler)    // Buscar por ID
		r.Put("/{id}", UpdateItemHandler) // Atualizar
		r.Delete("/{id}", DeleteItemHandler) // Deletar
	})

	log.Println("Servidor rodando em http://localhost:8080")
	http.ListenAndServe("0.0.0.0:8080", r)

}

// InitDB inicializa a conexão com o banco


// connectWithRetry tenta conectar ao banco até 5 vezes
func connectWithRetry() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	for i := 1; i <= 5; i++ { // tenta 5 vezes
		DB, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("Erro ao abrir conexão: %v", err)
		} else {
			err = DB.Ping()
			if err == nil {
				log.Println("Conexão com o banco de dados estabelecida!")
				return
			}
			log.Printf("Erro ao conectar (tentativa %d): %v", i, err)
		}
		time.Sleep(2 * time.Second) // espera 2 segundos antes da próxima tentativa
	}

	log.Fatal("Não foi possível conectar ao banco de dados após 5 tentativas")
}