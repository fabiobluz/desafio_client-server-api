# Desafio Client-Server API

Este projeto implementa um sistema cliente-servidor em Go que consome uma API de cotações de moedas e persiste os dados em banco SQLite.

## 📋 Funcionalidades

### Server (server.go)
- **Endpoint HTTP**: `/cotacao` na porta 8080
- **Consumo de API**: Integração com `https://economia.awesomeapi.com.br/json/last/USD-BRL`
- **Timeout de API**: 200ms para chamadas à API externa
- **Persistência**: Salva cotações no banco SQLite com timeout de 10ms
- **Resposta JSON**: Retorna apenas o campo `bid` da cotação
- **Logs**: Registra erros de timeout e falhas com contexto

### Client (client.go)
- **Requisição HTTP**: Faz chamada para o servidor local
- **Timeout**: 300ms para receber resposta do servidor
- **Arquivo de saída**: Salva cotação em `cotacao.txt` no formato "Dólar: {valor}"
- **Tratamento de erros**: Logs detalhados para falhas

## 🏗️ Arquitetura

```
┌─────────────┐    HTTP Request     ┌─────────────┐    API Call     ┌─────────────────────┐
│   Client    │ ──────────────────► │   Server    │ ──────────────► │  economia.awesomeapi │
│             │                     │             │                 │      .com.br        │
│             │ ◄────────────────── │             │ ◄────────────── │                     │
└─────────────┘    JSON Response    └─────────────┘    USD-BRL      └─────────────────────┘
       │                                    │
       │                                    ▼
       │                            ┌─────────────┐
       │                            │   SQLite    │
       │                            │   Database  │
       │                            └─────────────┘
       │
       ▼
┌─────────────┐
│ cotacao.txt │
└─────────────┘
```

## 🚀 Como Executar

### Pré-requisitos
- Go 1.24.2 ou superior
- Git

### 1. Clone o repositório
```bash
git clone <url-do-repositorio>
cd desafio_client-server-api
```

### 2. Execute o Servidor
```bash
cd server
go mod tidy
go run main.go
```

O servidor estará disponível em: `http://localhost:8080`

### 3. Execute o Cliente (em outro terminal)
```bash
cd client
go mod tidy
go run main.go
```

### 4. Verificar Resultado
Após a execução do cliente, verifique o arquivo `client/cotacao.txt`:
```
Dólar: 5.6457
```

## 🧪 Executar Testes

### Testes do Servidor
```bash
cd server
go test -v
```

### Testes do Cliente
```bash
cd client
go test -v
```

## 📊 Estrutura do Projeto

```
desafio_client-server-api/
├── README.md
├── server/
│   ├── main.go          # Servidor HTTP com endpoint /cotacao
│   ├── main_test.go     # Testes unitários do servidor
│   ├── go.mod           # Dependências do servidor
│   ├── go.sum           # Checksums das dependências
│   └── cotacoes.db      # Banco SQLite (criado automaticamente)
└── client/
    ├── main.go          # Cliente que consome o servidor
    ├── main_test.go     # Testes unitários do cliente
    ├── go.mod           # Dependências do cliente
    ├── go.sum           # Checksums das dependências
    └── cotacao.txt      # Arquivo de saída (criado automaticamente)
```

## 🔧 Configurações de Timeout

| Componente | Timeout | Contexto |
|------------|---------|----------|
| API Externa | 200ms | Chamada para economia.awesomeapi.com.br |
| Banco SQLite | 10ms | Inserção de dados |
| Cliente | 300ms | Requisição para o servidor |

## 📝 Estrutura do Banco de Dados

```sql
CREATE TABLE IF NOT EXISTS cotacoes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    bid TEXT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## 🔍 Logs e Monitoramento

O sistema registra logs detalhados para:
- Erros de timeout nas chamadas de API
- Falhas na persistência no banco de dados
- Erros de conexão e decodificação JSON
- Identificação única de requisições (request_id)

## 📦 Dependências

### Server
- `github.com/mattn/go-sqlite3` - Driver SQLite para Go

### Client
- `github.com/google/uuid` - Geração de UUIDs para identificação de requisições

## 🎯 Requisitos Atendidos

- ✅ Client.go realiza requisição HTTP no server.go
- ✅ Server.go consome API de cotações USD-BRL
- ✅ Retorno em formato JSON com campo "bid"
- ✅ Context com timeout de 200ms para API externa
- ✅ Context com timeout de 10ms para banco de dados
- ✅ Context com timeout de 300ms no cliente
- ✅ Persistência no SQLite
- ✅ Salvamento em arquivo cotacao.txt
- ✅ Endpoint /cotacao na porta 8080
- ✅ Logs de erro para timeouts

## 🚨 Tratamento de Erros

O sistema implementa tratamento robusto de erros:
- **Timeouts**: Logs específicos quando operações excedem o tempo limite
- **Falhas de API**: Retorna erro HTTP 500 com mensagem descritiva
- **Falhas de banco**: Logs de erro sem interromper o fluxo principal
- **Falhas de arquivo**: Tratamento de erros na criação/escrita de arquivos

## 📈 Performance

- **Execução assíncrona**: Salvamento no banco em goroutine separada
- **Timeouts otimizados**: Configurações específicas para cada operação
- **Reutilização de conexões**: Uso do http.DefaultClient
- **Prepared statements**: Otimização de queries SQL

---

**Desenvolvido como parte do desafio de Client-Server API em Go**