package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// InitDB cria a conexão com o MariaDB
func InitDB() {
	connectWithRetry()
	// Pegando variáveis de ambiente
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, name)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Erro ao abrir conexão com o banco: %v", err)
	}

	// Testa conexão
	if err = db.Ping(); err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}
	log.Println("Banco conectado com sucesso!")

	//cria a tabela
	query := `
	CREATE TABLE IF NOT EXISTS items (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) NOT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("Erro ao criar tabela: %v", err)
	}

	log.Println("Tabela 'items' criada com sucesso!")
}
