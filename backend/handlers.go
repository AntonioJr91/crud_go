package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"github.com/go-chi/chi/v5"
)

func ListItemsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, email FROM items")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var c Item
		if err := rows.Scan(&c.ID, &c.Name, &c.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func GetItemHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var c Item
	err := db.QueryRow("SELECT id, name, email FROM items WHERE id = ?", id).
		Scan(&c.ID, &c.Name, &c.Email)

	if err == sql.ErrNoRows {
		http.Error(w, "Item não encontrado", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	var c Item
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if c.Name == "" || c.Email == ""{
		http.Error(w, "Campos obrigatórios: nome, email", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO items (name, email) VALUES (?, ?)",
		c.Name, c.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	insertedID, _ := result.LastInsertId()
	c.ID = insertedID

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var c Item
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`
		UPDATE items 
		SET name = ?, email = ?
		WHERE id = ?`,
		c.Name, c.Email, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Item não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Item atualizado com sucesso"})
}

func DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result, err := db.Exec("DELETE FROM items WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Item não encontrado", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
