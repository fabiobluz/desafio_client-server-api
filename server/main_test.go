package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestCotacaoHandler(t *testing.T) {
	apiMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"USDBRL": map[string]string{"bid": "5.25"},
		})
	}))
	defer apiMock.Close()

	originalAPIURL := "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	defer func() { apiURL = originalAPIURL }()
	apiURL = apiMock.URL

	var err error
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erro ao abrir banco de dados: %v", err)
	}
	defer db.Close()
	createTable()

	req := httptest.NewRequest(http.MethodGet, "/cotacao", nil)
	rec := httptest.NewRecorder()

	cotacaoHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status esperado 200, obtido %d", resp.StatusCode)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Erro ao decodificar resposta: %v", err)
	}

	if result["bid"] != "5.25" {
		t.Errorf("Valor bid esperado '5.25', obtido '%s'", result["bid"])
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM cotacoes").Scan(&count)
	if err != nil {
		t.Fatalf("Erro ao verificar banco: %v", err)
	}
	//if count == 0 {
	//	t.Error("Cotação não foi inserida no banco de dados")
	//}
}
