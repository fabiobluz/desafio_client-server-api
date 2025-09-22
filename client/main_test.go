package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestBuscarCotacao(t *testing.T) {
	// Simula um servidor com resposta mockada
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"bid": "5.99"}`))
	}))
	defer mockServer.Close()

	valor, err := buscarCotacao(mockServer.URL)
	if err != nil {
		t.Fatalf("Erro ao buscar cotação: %v", err)
	}

	if valor != "5.99" {
		t.Errorf("Esperado '5.99', obtido '%s'", valor)
	}
}

func TestSalvarCotacaoEmArquivo(t *testing.T) {
	arquivo := "teste_cotacao.txt"
	defer os.Remove(arquivo)

	err := salvarCotacaoEmArquivo(arquivo, "5.99")
	if err != nil {
		t.Fatalf("Erro ao salvar cotação: %v", err)
	}

	conteudo, err := ioutil.ReadFile(arquivo)
	if err != nil {
		t.Fatalf("Erro ao ler arquivo: %v", err)
	}

	esperado := "Dólar: 5.99"
	if !strings.Contains(string(conteudo), esperado) {
		t.Errorf("Conteúdo esperado '%s', obtido '%s'", esperado, string(conteudo))
	}
}
