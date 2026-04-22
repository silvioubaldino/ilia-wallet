# ilia-wallet

Microserviço de carteira digital — Parte 1 do desafio ília.

Armazena e consulta transações financeiras (CREDIT/DEBIT) por usuário. O balance é calculado via query SQL; não há tabela de usuários ou carteiras.

## Stack

- **Go 1.22+** · Uber style guide
- **Gin** — HTTP
- **GORM + PostgreSQL** — persistência
- **golang-migrate** — migrations SQL embed
- **golang-jwt/jwt v5** — autenticação JWT

## Pré-requisitos

- Go 1.22+
- Docker e Docker Compose

## Variáveis de ambiente

| Variável | Padrão | Descrição |
|---|---|---|
| `SERVER_PORT` | `3001` | Porta HTTP |
| `DB_HOST` | `localhost` | Host do PostgreSQL |
| `DB_PORT` | `5432` | Porta do PostgreSQL |
| `DB_USER` | — | Usuário do banco |
| `DB_PASSWORD` | — | Senha do banco |
| `DB_NAME` | — | Nome do banco |
| `DB_SSLMODE` | `disable` | SSL mode |
| `ILIACHALLENGE` | — | Segredo de assinatura JWT |

Copie `.env.example` e ajuste os valores:

```bash
cp .env.example .env
```

## Como rodar

### Com Docker Compose (recomendado)

```bash
docker compose up --build
```

O serviço sobe em `http://localhost:3001`. As migrations são aplicadas automaticamente na inicialização.

### Localmente

Suba apenas o banco:

```bash
docker compose up postgres -d
```

Execute a aplicação:

```bash
go run ./cmd/server
```

## Testes

```bash
go test ./...
```

## API

Todos os endpoints protegidos exigem o header:

```
Authorization: Bearer <JWT>
```

O JWT deve ser assinado com o segredo definido em `ILIACHALLENGE`. O claim `user_id` presente no token é usado para isolar as transações de cada usuário.

---

### GET /health

Liveness check — sem autenticação.

**Response 200**
```json
{ "status": "ok" }
```

---

### POST /transactions

Cria uma transação de CRÉDITO ou DÉBITO.

**Request**
```json
{
  "user_id": "uuid",
  "type": "CREDIT | DEBIT",
  "amount": 100
}
```

O `user_id` do body deve coincidir com o claim do JWT; caso contrário, retorna 401.

**Response 200**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "type": "CREDIT",
  "amount": 100
}
```

---

### GET /transactions

Lista as transações do usuário autenticado. Suporta filtro opcional por tipo.

**Query params**

| Param | Valores |
|---|---|
| `type` | `CREDIT` ou `DEBIT` (opcional) |

**Response 200**
```json
[
  {
    "id": "uuid",
    "user_id": "uuid",
    "type": "CREDIT",
    "amount": 100
  }
]
```

---

### GET /balance

Retorna o saldo consolidado do usuário autenticado.

Calculado como `SUM(CREDIT) - SUM(DEBIT)` via query SQL.

**Response 200**
```json
{ "amount": 70 }
```
