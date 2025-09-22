package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	valor, err := buscarCotacao("http://localhost:8080/cotacao")
	if err != nil {
		log.Fatal(err)
	}

	err = salvarCotacaoEmArquivo("cotacao.txt", valor)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Cotação salva com sucesso:", valor)
}

// buscarCotacao faz uma requisição com timeout e retorna o valor "bid"
func buscarCotacao(url string) (string, error) {

	requestID := uuid.New().String()
	ctx := context.WithValue(context.Background(), "request_id", requestID)

	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("erro ao criar requisição: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao obter resposta do servidor: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("erro ao decodificar JSON: %w", err)
	}

	return result["bid"], nil
}

// salvarCotacaoEmArquivo grava o valor no formato correto no arquivo
func salvarCotacaoEmArquivo(nomeArquivo, valor string) error {
	file, err := os.Create(nomeArquivo)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar: %s", valor))
	if err != nil {
		return fmt.Errorf("erro ao escrever no arquivo: %w", err)
	}

	return nil
}
