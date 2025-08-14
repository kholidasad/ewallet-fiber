# E‑Wallet API (Go Fiber + MySQL) — `kholid/ewallet/v2`

Multi‑currency **e‑wallet** API dengan **JWT auth**, **pagination**, **real‑time FX rates** (cache), dan **MySQL**. 
Dibangun pakai **Go Fiber**, **GORM**, dan **shopspring/decimal** untuk akurasi uang.

---

## Fitur
- Users & multi‑wallet, multi‑currency balances (USD, EUR, JPY default).
- Transaksi: **deposit**, **withdraw**, **transfer** (dengan konversi), **payment**.
- Catatan transaksi lengkap (amount, currency, rate, converted, status, reference, timestamp).
- **Pagination** di `/api/v1/transactions` dengan indeks yang sesuai.
- **Live FX** (`exchangerate.host` by default) + cache TTL, fallback **static rates** via env.
- **JWT auth** (register/login), user hanya bisa akses resource miliknya.
- **Swagger UI**: `http://localhost:8080/swagger/index.html`.

---

## Arsitektur (High Level)

```
[Client] ──HTTP──> [Fiber API Handlers]
                      │
                      ▼
               [Service Layer]
          (Auth, Wallet, Transaction, FX)
                      │
                      ▼
                [GORM Repository]
                      │
                      ▼
                   [MySQL DB]
                      │
                      └── (Indexes: uniq_wallet_currency, idx_trx_wallet_created, ...)

           [Live FX Provider API]  <─ cache TTL ─
```

**Stateless** di layer API → cocok untuk **horizontal scaling**. Konsistensi saldo: transaksi DB + **row‑level lock** (`SELECT ... FOR UPDATE` via GORM `clause.Locking`).

---

## Teknologi
- Go 1.22 (jalan juga di 1.24), Fiber v2, GORM
- MySQL 8
- JWT (golang‑jwt v5)
- Decimal math (shopspring/decimal)
- Swagger (gofiber/swagger + swag runtime loader)

---

## Struktur Project
```
.
├─ cmd/server/main.go
├─ internal/
│  ├─ config/            # load env (DB_HOST/PORT/USERNAME/PASSWORD/NAME, JWT, FX)
│  ├─ db/                # connect & migrate (GORM AutoMigrate)
│  ├─ middleware/        # JWT auth middleware
│  ├─ models/            # GORM models
│  ├─ handler/           # Fiber handlers (REST)
│  └─ service/           # business logic (auth, wallet, transaction, currency)
├─ migrations/mysql/001_init.sql   # schema + indexes + seed currencies
├─ docs/                 # minimal Swagger embed
├─ .github/workflows/    # ci.yml & cicd.yml
├─ Dockerfile
├─ docker-compose.yml    # dev MySQL
├─ Makefile
├─ .env.example
└─ go.mod (module: kholid/ewallet/v2)
```

---

## Database Schema & Indexes (ringkas)
Tabel inti:
- **users** (email unik, password hash)
- **wallets** (FK → users)
- **currencies** (seed USD/EUR/JPY)
- **wallet_balances** (FK → wallets & currencies)  
  - `UNIQUE (wallet_id, currency_code)` → **`uniq_wallet_currency`**
  - `INDEX (wallet_id, currency_code)` → **`idx_wb_wallet_cur`**
- **transactions** (FK → wallets)  
  Index penting:
  - `INDEX (wallet_id, created_at)` → **`idx_trx_wallet_created`** (pagination & riwayat)
  - `INDEX (type)`, `INDEX (currency_code)`, `INDEX (status)`, `INDEX (reference)`

> File SQL: `migrations/mysql/001_init.sql` (disertai seed currencies). Runtime juga menjalankan **GORM AutoMigrate**.

---

## Environment Variables
```
APP_ENV=development
PORT=8080

DB_HOST=127.0.0.1
DB_PORT=3306
DB_USERNAME=user
DB_PASSWORD=password
DB_NAME=ewallet

JWT_SECRET=supersecret
DISPLAY_CURRENCY=USD

# FX
FX_PROVIDER=live                          # live|static
FX_PROVIDER_URL=https://api.exchangerate.host/latest
FX_BASE=USD
FX_CURRENCIES=USD,EUR,JPY
FX_TTL_SECONDS=300

# Static rates (fallback jika FX_PROVIDER=static)
RATE_USD=1
RATE_EUR=0.9
RATE_JPY=155
```

---

## Instalasi & Menjalankan

### 1) Prasyarat
- Go 1.22+
- Docker & Docker Compose

### 2) Jalankan MySQL dev
```bash
docker compose up -d db
```

### 3) Konfigurasi Env
```bash
cp .env.example .env
# edit jika perlu (DB_HOST/PORT/USERNAME/PASSWORD/NAME, JWT_SECRET, dll)
```

### 4) Jalankan API
```bash
export $(grep -v '^#' .env | xargs)    # macOS/Linux
go run ./cmd/server
# atau: make run
```

### 5) Swagger
Buka: **http://localhost:8080/swagger/index.html**

---

## Autentikasi
- **Register**: `POST /api/v1/auth/register`
- **Login**: `POST /api/v1/auth/login` → balikan `{ "token": "<JWT>" }`
- Sertakan header: `Authorization: Bearer <JWT>` untuk semua endpoint selain register/login.

---

## API Endpoints (ringkas + contoh)

### Register
```
POST /api/v1/auth/register
Content-Type: application/json

{ "Email": "user@demo.io", "Password": "secret" }
```

### Login
```
POST /api/v1/auth/login
{ "Email": "user@demo.io", "Password": "secret" }
⇒ { "token": "...", "user": {...} }
```

### Wallet
- **Create**: `POST /api/v1/wallets` `{ "name": "Main" }`
- **List**: `GET /api/v1/wallets`
- **Balances**: `GET /api/v1/wallets/:id/balances`

### Transaksi

**Deposit**
```
POST /api/v1/wallets/:id/deposit
{ "Currency":"USD", "Amount":"100.00" }
```

**Withdraw**
```
POST /api/v1/wallets/:id/withdraw
{ "Currency":"USD", "Amount":"40" }
```

**Transfer** (konversi otomatis)
```
POST /api/v1/wallets/:id/transfer
{
  "to_wallet_id": 2,
  "FromCurrency": "USD",
  "ToCurrency": "EUR",
  "Amount": "50.25"
}
```

**Payment**
```
POST /api/v1/wallets/:id/payment
{
  "Currency": "USD",
  "Amount": "12.99",
  "Reference": "ORDER-123",
  "Metadata": "book purchase"
}
```

### Riwayat Transaksi (+ Pagination)
```
GET /api/v1/transactions?page=1&page_size=20&type=deposit&currency=USD
⇒ {
  "data": [ ... ],
  "meta": { "page":1, "page_size":20, "total":123, "total_pages":7 }
}
```
Sorting default: `created_at DESC, id DESC`. Batas `page_size` 1–100 (default 20).

### Ringkasan Saldo (display currency)
```
GET /api/v1/summary
⇒ { "display_currency":"USD", "total":"123.45", "by_currency": {"USD":"100","EUR":"20","JPY":"5000"} }
```

---

## Contoh Quick Run (cURL)
```bash
# Register
curl -s -X POST http://localhost:8080/api/v1/auth/register   -H "Content-Type: application/json"   -d '{"Email":"a@b.c","Password":"secret"}'

# Login (ambil token)
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login   -H "Content-Type: application/json"   -d '{"Email":"a@b.c","Password":"secret"}' | jq -r .token)

# Create wallet
curl -s -X POST http://localhost:8080/api/v1/wallets   -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json"   -d '{"name":"Main"}'

# Deposit
curl -s -X POST http://localhost:8080/api/v1/wallets/1/deposit   -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json"   -d '{"Currency":"USD","Amount":"100.00"}'

# List transaksi (pagination)
curl -s -H "Authorization: Bearer $TOKEN"   "http://localhost:8080/api/v1/transactions?page=1&page_size=10"
```

---

## Real‑Time FX (Live) vs Static
- `FX_PROVIDER=live` → fetch rates dari `FX_PROVIDER_URL` (default: exchangerate.host). Cache TTL = `FX_TTL_SECONDS`.
- `FX_PROVIDER=static` → pakai env `RATE_USD/EUR/JPY`.
- Konversi: `converted = amount / fromRate * toRate`, semua presisi decimal.

---

## Concurrency & Konsistensi
- Mutasi saldo berjalan di **transaction DB**.
- Baris saldo di‑lock via **`FOR UPDATE`** (`gorm.io/gorm/clause.Locking`), mencegah race/double spend.
- Idempotency‑Key belum di‑implement (bisa ditambah di header untuk POST mutasi → simpan ke tabel `idempotency_keys`).

---

## Error Handling
- Format umum: `{"error":"message"}` dengan kode HTTP tepat (`400`, `401`, `404`, `500`).
- Validasi amount > 0, currency valid, kepemilikan wallet, saldo cukup.

---

## Testing
- **Unit tests** contoh: `internal/service/currency_test.go`, `wallet_service_test.go`.
  ```bash
  go test ./...
  ```
- **Integration test** (butuh MySQL):
  ```bash
  export TEST_DB_DSN="user:password@tcp(127.0.0.1:3306)/ewallet?charset=utf8mb4&parseTime=True&loc=Local"
  go test -tags=integration ./... -v
  ```

---

## Docker
```bash
# Build
docker build -t ewallet-fiber:dev .

# Run
docker run -p 8080:8080 --env-file .env ewallet-fiber:dev
```

---