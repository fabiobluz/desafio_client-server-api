package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ExchangeRate struct {
	Bid string `json:"bid"`
}

type ExchangeResponse struct {
	USDBRL ExchangeRate `json:"USDBRL"`
}

var apiURL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./cotacoes.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTable()

	http.HandleFunc("/cotacao", cotacaoHandler)
	log.Println("Servidor iniciado na porta 8080")
	http.ListenAndServe(":8080", nil)
}

func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	ctxAPI, cancelAPI := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancelAPI()

	req, err := http.NewRequestWithContext(ctxAPI, http.MethodGet, apiURL, nil)
	//req, err := http.NewRequestWithContext(ctxAPI, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		http.Error(w, "Erro ao criar requisição", http.StatusInternalServerError)
		log.Println("Erro ao criar requisição:", err, fmt.Sprintf("request_id: %v", ctxAPI.Value("request_id")))
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Erro ao obter cotação", http.StatusInternalServerError)
		log.Println("Erro ao obter cotação:", err, fmt.Sprintf("request_id: %v", ctxAPI.Value("request_id")))
		return
	}
	defer resp.Body.Close()

	var exchange ExchangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&exchange); err != nil {
		http.Error(w, "Erro ao decodificar resposta", http.StatusInternalServerError)
		log.Println("Erro ao decodificar JSON:", err, fmt.Sprintf("request_id: %v", ctxAPI.Value("request_id")))
		return
	}

	go saveToDatabase(ctxAPI, exchange.USDBRL.Bid)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"bid": exchange.USDBRL.Bid})
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS cotacoes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bid TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(query); err != nil {
		log.Fatal("Erro ao criar tabela:", err)
	}
}

func saveToDatabase(ctx context.Context, bid string) {
	ctxDB, cancelDB := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancelDB()

	query := "INSERT INTO cotacoes (bid) VALUES (?)"
	stmt, err := db.PrepareContext(ctxDB, query)
	if err != nil {
		log.Println("Erro ao preparar statement:", err, fmt.Sprintf("request_id: %v", ctx.Value("request_id")))
		return
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctxDB, bid); err != nil {
		log.Println("Erro ao inserir no banco de dados:", err, fmt.Sprintf("request_id: %v", ctx.Value("request_id")))
	}
}
