package models

import (
	"time"
	"github.com/shopspring/decimal"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex;size:190" json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Wallets   []Wallet  `json:"wallets"`
}

type Wallet struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Name      string    `gorm:"size:120" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Balances  []WalletBalance `json:"balances"`
}

type Currency struct { Code string `gorm:"primaryKey;size:3" json:"code"`; Name string `gorm:"size:50" json:"name"` }

type WalletBalance struct {
	ID           uint            `gorm:"primaryKey" json:"id"`
	WalletID     uint            `gorm:"index:idx_wb_wallet_cur,priority:1;index" json:"wallet_id"`
	CurrencyCode string          `gorm:"size:3;index:idx_wb_wallet_cur,priority:2" json:"currency_code"`
	Balance      decimal.Decimal `gorm:"type:decimal(38,18)" json:"balance"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type TransactionType string
const (
	TrxDeposit TransactionType = "deposit"
	TrxWithdraw TransactionType = "withdrawal"
	TrxTransfer TransactionType = "transfer"
	TrxPayment TransactionType = "payment"
)

type TransactionStatus string
const (
	TrxSuccess TransactionStatus = "success"
	TrxFailed TransactionStatus = "failed"
)

type Transaction struct {
	ID              uint              `gorm:"primaryKey" json:"id"`
	WalletID        uint              `gorm:"index:idx_trx_wallet_created,priority:1;index" json:"wallet_id"`
	ToWalletID      *uint             `json:"to_wallet_id,omitempty"`
	Type            TransactionType   `gorm:"size:20;index" json:"type"`
	Status          TransactionStatus `gorm:"size:20;index" json:"status"`
	Amount          decimal.Decimal   `gorm:"type:decimal(38,18)" json:"amount"`
	CurrencyCode    string            `gorm:"size:3;index" json:"currency_code"`
	ExchangeRate    decimal.Decimal   `gorm:"type:decimal(38,18)" json:"exchange_rate"`
	ConvertedAmount decimal.Decimal   `gorm:"type:decimal(38,18)" json:"converted_amount"`
	Reference       string            `gorm:"size:120;index" json:"reference"`
	Metadata        string            `gorm:"type:text" json:"metadata"`
	CreatedAt       time.Time         `gorm:"index:idx_trx_wallet_created,priority:2" json:"created_at"`
}
