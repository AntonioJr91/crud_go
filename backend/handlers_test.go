package main

import(
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"bytes"
	"context"

	"github.com/DATA-DOG/go-sqlmock"
    "github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
);

//GET
func TestListItemsHandler(t *testing.T) {
	var mockDB *sql.DB
	var mock sqlmock.Sqlmock
	var err error

	mockDB, mock, err = sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	// substitui a variável global db pelo mock
	db = mockDB

	rows := sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(1, "Doka", "doka@example.com").
		AddRow(2, "Coto", "@example.com")

	mock.ExpectQuery("SELECT id, name, email FROM items").
		WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	w := httptest.NewRecorder()

	ListItemsHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var items []Item
	err = json.NewDecoder(resp.Body).Decode(&items)
	require.NoError(t, err)

	require.Len(t, items, 2)
	require.Equal(t, "Doka", items[0].Name)
	require.Equal(t, "Coto", items[1].Name)

	require.NoError(t, mock.ExpectationsWereMet())
}

//GET:ID
func TestGetItemHandler(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()
	db = mockDB

	rows := sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(1, "Doka", "doka@example.com")

	mock.ExpectQuery("SELECT id, name, email FROM items WHERE id = ?").
		WithArgs("1").
		WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/items/1", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	GetItemHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var item Item
	err = json.NewDecoder(resp.Body).Decode(&item)
	require.NoError(t, err)
	require.Equal(t, int64(1), item.ID)
	require.Equal(t, "Doka", item.Name)
	require.Equal(t, "doka@example.com", item.Email)

	mock.ExpectQuery("SELECT id, name, email FROM items WHERE id = ?").
		WithArgs("2").
		WillReturnError(sql.ErrNoRows)

	req2 := httptest.NewRequest(http.MethodGet, "/items/2", nil)
	w2 := httptest.NewRecorder()

	rctx2 := chi.NewRouteContext()
	rctx2.URLParams.Add("id", "2")
	req2 = req2.WithContext(context.WithValue(req2.Context(), chi.RouteCtxKey, rctx2))

	GetItemHandler(w2, req2)

	resp2 := w2.Result()
	defer resp2.Body.Close()
	require.Equal(t, http.StatusNotFound, resp2.StatusCode)

	require.NoError(t, mock.ExpectationsWereMet())
}

//POST
func TestCreateItemHandler(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()
	db = mockDB

	mock.ExpectExec("INSERT INTO items").
		WithArgs("João", "joao@example.com").
		WillReturnResult(sqlmock.NewResult(10, 1))

	body := []byte(`{"name":"João","email":"joao@example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateItemHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)
	require.NoError(t, mock.ExpectationsWereMet())
}

//PUT
func TestUpdateItemHandler(t *testing.T) {
	var mockDB *sql.DB
	var mock sqlmock.Sqlmock
	var err error

	mockDB, mock, err = sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	db = mockDB

	item := Item{Name: "Novo Nome", Email: "novo@email.com"}
	body, _ := json.Marshal(item)

	req := httptest.NewRequest(http.MethodPut, "/items/1", bytes.NewReader(body))
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	mock.ExpectExec("UPDATE items").
		WithArgs(item.Name, item.Email, "1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	UpdateItemHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var res map[string]string
	err = json.NewDecoder(resp.Body).Decode(&res)
	require.NoError(t, err)
	require.Equal(t, "Item atualizado com sucesso", res["message"])

	req2 := httptest.NewRequest(http.MethodPut, "/items/2", bytes.NewReader(body))
	w2 := httptest.NewRecorder()

	rctx2 := chi.NewRouteContext()
	rctx2.URLParams.Add("id", "2")
	req2 = req2.WithContext(context.WithValue(req2.Context(), chi.RouteCtxKey, rctx2))

	mock.ExpectExec("UPDATE items").
		WithArgs(item.Name, item.Email, "2").
		WillReturnResult(sqlmock.NewResult(1, 0))

	UpdateItemHandler(w2, req2)
	resp2 := w2.Result()
	defer resp2.Body.Close()
	require.Equal(t, http.StatusNotFound, resp2.StatusCode)

	require.NoError(t, mock.ExpectationsWereMet())
}

//DEL
func TestDeleteItemHandler(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()
	db = mockDB

	req := httptest.NewRequest(http.MethodDelete, "/items/1", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	mock.ExpectExec("DELETE FROM items").
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	DeleteItemHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	req2 := httptest.NewRequest(http.MethodDelete, "/items/2", nil)
	w2 := httptest.NewRecorder()

	rctx2 := chi.NewRouteContext()
	rctx2.URLParams.Add("id", "2")
	req2 = req2.WithContext(context.WithValue(req2.Context(), chi.RouteCtxKey, rctx2))

	mock.ExpectExec("DELETE FROM items").
		WithArgs("2").
		WillReturnResult(sqlmock.NewResult(1, 0))

	DeleteItemHandler(w2, req2)

	resp2 := w2.Result()
	defer resp2.Body.Close()
	require.Equal(t, http.StatusNotFound, resp2.StatusCode)

	require.NoError(t, mock.ExpectationsWereMet())
}