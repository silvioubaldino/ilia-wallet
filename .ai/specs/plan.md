# Wallet Microservice — SDD Implementation Plan (Clean Architecture)

## Context

Implementar o microsserviço de carteira digital (Parte 1 do desafio ília) em **Go** com **Clean Architecture**. O escopo é definido pelo arquivo `.claude/ms-transactions.yaml` — nenhum endpoint além dos listados ali deve ser implementado.

---

## Por que não precisamos de CRUD de usuários?

O microserviço wallet tem responsabilidade única: **armazenar e consultar transações financeiras**. O `user_id` vem no body da requisição e é validado contra o claim do JWT. Não há endpoints de `/auth/*` no spec. Gerenciamento de usuários é escopo da Parte 2 (Users Microservice). O balance é calculado via query (SUM CREDIT - SUM DEBIT) — não há tabela `wallets` separada.

---

## Tech Stack

| Concern | Choice |
|---|---|
| Language | Go 1.22+ |
| Style | Uber Go style guide |
| HTTP | `github.com/gin-gonic/gin` |
| ORM | `gorm.io/gorm` + `gorm.io/driver/postgres` |
| JWT | `github.com/golang-jwt/jwt/v5` |
| Decimals | `github.com/shopspring/decimal` |
| UUIDs | `github.com/google/uuid` |
| Testing | `github.com/stretchr/testify` + `testify/mock` |

---

## API Endpoints (source of truth: `.claude/ms-transactions.yaml`)

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/health` | No | Liveness check |
| POST | `/transactions` | JWT | Create transaction (CREDIT or DEBIT) |
| GET | `/transactions` | JWT | List transactions (filter `?type=CREDIT\|DEBIT`) |
| GET | `/balance` | JWT | Consolidated balance (SUM via query) |

### POST /transactions — Request
```json
{ "user_id": "string", "type": "CREDIT|DEBIT", "amount": 100 }
```

### POST /transactions — Response 200
```json
{ "id": "string", "user_id": "string", "type": "CREDIT", "amount": 100 }
```

### GET /transactions — Response 200
```json
[{ "id": "string", "user_id": "string", "type": "CREDIT", "amount": 100 }]
```

### GET /balance — Response 200
```json
{ "amount": 250 }
```
Balance = SUM(amount WHERE type=CREDIT) - SUM(amount WHERE type=DEBIT), filtrado pelo `user_id` do JWT claim.

---

## Database Schema

### `transactions` (única tabela necessária)

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PK, gen_random_uuid() |
| user_id | UUID | NOT NULL, INDEX |
| type | VARCHAR(10) | NOT NULL, CHECK IN ('CREDIT','DEBIT') |
| amount | INTEGER | NOT NULL, CHECK > 0 |
| created_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |

**Sem tabela `users`, sem tabela `wallets`.** O balance é calculado via SQL:
```sql
SELECT
  SUM(CASE WHEN type = 'CREDIT' THEN amount ELSE -amount END) AS amount
FROM transactions
WHERE user_id = $1
```

---

## Clean Architecture — Estrutura de Diretórios

```
ilia-wallet/
├── cmd/
│   └── server/
│       └── main.go                        # Wire all deps, start HTTP server
│
├── internal/
│   ├── domain/
│   │   └── transaction/
│   │       ├── entity.go                  # Transaction entity + Type enum (CREDIT/DEBIT)
│   │       └── repository.go              # Repository interface (port)
│   │
│   ├── usecase/
│   │   └── transaction/
│   │       ├── create.go                  # CreateTransaction use case
│   │       ├── create_test.go
│   │       ├── list.go                    # ListTransactions use case
│   │       ├── list_test.go
│   │       ├── balance.go                 # GetBalance use case
│   │       └── balance_test.go
│   │
│   ├── adapter/
│   │   ├── http/
│   │   │   ├── handler/
│   │   │   │   ├── transaction.go         # Gin handlers
│   │   │   │   └── transaction_test.go
│   │   │   └── middleware/
│   │   │       ├── auth.go                # JWT validation gin.HandlerFunc
│   │   │       └── auth_test.go
│   │   └── repository/
│   │       └── postgres/
│   │           ├── transaction.go         # GORM implementation of domain.Repository
│   │           └── transaction_test.go
│   │
│   └── infrastructure/
│       ├── config/
│       │   └── config.go                  # Config struct, Load() from env
│       └── database/
│           └── postgres.go                # GORM connection + AutoMigrate
│
├── pkg/
│   └── apperrors/
│       └── errors.go                      # Sentinel errors
│
├── docker/
│   └── Dockerfile                         # Multi-stage: golang:1.22-alpine → alpine
│
├── docker-compose.yml
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### Dependency Rule (Clean Arch)
```
Infrastructure → Adapter → UseCase → Domain
```
- **Domain**: entidades e interfaces de repositório (zero dependências externas)
- **UseCase**: lógica de negócio, depende apenas do Domain
- **Adapter**: HTTP handlers e implementação do repositório (GORM), depende de UseCase
- **Infrastructure**: configuração, DB connection, wiring — depende de tudo

---

## Development Stages

### Stage 0 — Criar `.claude/specs/plan.md` ✅
**Deliverable:** Arquivo de spec criado e commitado no repositório

- Criar `mkdir -p .claude/specs` e `cat` o conteúdo deste plan para `.claude/specs/plan.md`

---

### Stage 1 — Project Bootstrap
**Deliverable:** Go module + Gin em :3001 + `GET /health` + Docker Compose funcionando

- `go mod init github.com/silvioubaldino/ilia-wallet`
- `internal/infrastructure/config/config.go`
- `internal/infrastructure/database/postgres.go`
- `cmd/server/main.go` com health route
- `docker/Dockerfile` (multi-stage)
- `docker-compose.yml` (app + postgres)
- `.env.example`

**Verify:** `docker-compose up` → `curl localhost:3001/health` → 200

---

### Stage 2 — Domain + Database Migration
**Deliverable:** Entidade Transaction definida, tabela criada via AutoMigrate

- `internal/domain/transaction/entity.go`
- `internal/domain/transaction/repository.go` (interface)
- `internal/infrastructure/database/postgres.go` com AutoMigrate da tabela `transactions`

**Verify:** `psql` → `\d transactions` mostra schema correto

---

### Stage 3 — JWT Middleware
**Deliverable:** Todas as rotas protegidas retornam 401 sem token válido

- `internal/adapter/http/middleware/auth.go`: valida JWT com `ILIACHALLENGE`, injeta claims no context
- `internal/adapter/http/middleware/auth_test.go`: table-driven, casos: sem header / expirado / assinatura inválida / válido

**Verify:** `POST /transactions` sem token → 401; com JWT válido → passa

---

### Stage 4 — Create Transaction
**Deliverable:** `POST /transactions` funcionando end-to-end + unit tests

- `internal/adapter/repository/postgres/transaction.go`: implementa `repo.Create`
- `internal/usecase/transaction/create.go`: `CreateTransaction.Execute`
- `internal/usecase/transaction/create_test.go`: table-driven com mock do repo
- `internal/adapter/http/handler/transaction.go`: handler POST

**Verify:** `POST /transactions` com JWT válido cria registro no banco e retorna 200

---

### Stage 5 — List Transactions
**Deliverable:** `GET /transactions` funcionando com filtro por tipo + unit tests

- `internal/adapter/repository/postgres/transaction.go`: implementa `repo.List` com filtro opcional por type e user_id
- `internal/usecase/transaction/list.go`: `ListTransactions.Execute`
- `internal/usecase/transaction/list_test.go`
- Handler GET no `transaction.go`

**Verify:** `GET /transactions?type=CREDIT` retorna apenas CREDITs do user autenticado

---

### Stage 6 — Balance
**Deliverable:** `GET /balance` retorna consolidado via query no banco + unit tests

- `internal/adapter/repository/postgres/transaction.go`: implementa `repo.Balance` com query `SUM(CASE WHEN type='CREDIT' THEN amount ELSE -amount END)`
- `internal/usecase/transaction/balance.go`: `GetBalance.Execute`
- `internal/usecase/transaction/balance_test.go`
- Handler GET /balance

**Verify:** Após CREDIT 100 + DEBIT 30 → `GET /balance` retorna `{"amount": 70}`

---

### Stage 7 — Polish, Lint & Docs
**Deliverable:** Lint limpo, README completo

- `golangci-lint run ./...` e corrigir issues
- `Makefile`: `make run`, `make test`, `make lint`, `make docker-up`
- `README.md`: pré-requisitos, env vars, como rodar, referência da API

**Verify:** `golangci-lint` exits 0; smoke test end-to-end completo

---

## Environment Variables

| Var | Example | Description |
|---|---|---|
| `SERVER_PORT` | `3001` | HTTP server port |
| `DB_HOST` | `postgres` | Postgres host |
| `DB_PORT` | `5432` | Postgres port |
| `DB_USER` | `wallet` | Postgres user |
| `DB_PASSWORD` | `secret` | Postgres password |
| `DB_NAME` | `wallet_db` | Database name |
| `DB_SSLMODE` | `disable` | SSL mode |
| `ILIACHALLENGE` | `ILIACHALLENGE` | JWT signing secret |

---

## Testing Rules

- Pacote: `package <feature>_test`
- Formato: `map[string]struct{ input; mocks; expected }{}`
- Nome: `"expects/should ... when"`
- Padrão AAA com comentários `// Arrange`, `// Act`, `// Assert`
- Sem condicionais em testes
- `assert.ErrorIs` para erros, `assert.AnError` para mockar erros
- Mocks em `mock_test.go`, operator address-of, `AssertExpectations`
- Sem `mock.Anything`
