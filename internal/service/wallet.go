package service

import (
	"errors"; "fmt"; "time"
	"kholid/ewallet/v2/internal/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"; "gorm.io/gorm/clause"
)

type WalletService struct { db *gorm.DB; cc CurrencyService }
func NewWalletService(db *gorm.DB, cc CurrencyService) *WalletService { return &WalletService{db: db, cc: cc} }

func (s *WalletService) CreateWallet(userID uint, name string) (*models.Wallet, error) {
	w := &models.Wallet{UserID: userID, Name: name}; if err := s.db.Create(w).Error; err != nil { return nil, err }
	for _, code := range s.cc.List() { wb := models.WalletBalance{WalletID: w.ID, CurrencyCode: code, Balance: decimal.Zero}; s.db.Create(&wb) }
	return w, nil
}

func (s *WalletService) ListWallets(userID uint) ([]models.Wallet, error) {
	var ws []models.Wallet; if err := s.db.Where("user_id = ?", userID).Find(&ws).Error; err != nil { return nil, err }; return ws, nil
}

func (s *WalletService) WalletBalances(userID, walletID uint) ([]models.WalletBalance, error) {
	var w models.Wallet; if err := s.db.Where("id = ? AND user_id = ?", walletID, userID).First(&w).Error; err != nil { return nil, gorm.ErrRecordNotFound }
	var b []models.WalletBalance; if err := s.db.Where("wallet_id = ?", walletID).Find(&b).Error; err != nil { return nil, err }; return b, nil
}

func (s *WalletService) adjustBalance(tx *gorm.DB, walletID uint, currency string, delta decimal.Decimal) (*models.WalletBalance, error) {
	var wb models.WalletBalance
	if err := tx.Clauses(clause.Locking{Strength:"UPDATE"}).Where("wallet_id = ? AND currency_code = ?", walletID, currency).First(&wb).Error; err != nil { return nil, err }
	newBal := wb.Balance.Add(delta); if newBal.IsNegative() { return nil, errors.New("insufficient funds") }
	wb.Balance = newBal; wb.UpdatedAt = time.Now(); if err := tx.Save(&wb).Error; err != nil { return nil, err }
	return &wb, nil
}

func (s *WalletService) Deposit(userID, walletID uint, currency string, amount decimal.Decimal) (*models.Transaction, error) {
	if amount.LessThanOrEqual(decimal.Zero) { return nil, errors.New("amount must be > 0") }
	var w models.Wallet; if err := s.db.Where("id=? AND user_id=?", walletID, userID).First(&w).Error; err != nil { return nil, gorm.ErrRecordNotFound }
	var trx *models.Transaction; err := s.db.Transaction(func(tx *gorm.DB) error {
		if _, err := s.adjustBalance(tx, walletID, currency, amount); err != nil { return err }
		trx = &models.Transaction{ WalletID: w.ID, Type: models.TrxDeposit, Status: models.TrxSuccess, Amount: amount, CurrencyCode: currency, ExchangeRate: decimal.NewFromInt(1), ConvertedAmount: amount, Reference: fmt.Sprintf("DEP-%d", time.Now().UnixNano()) }
		return tx.Create(trx).Error
	}); return trx, err
}

func (s *WalletService) Withdraw(userID, walletID uint, currency string, amount decimal.Decimal) (*models.Transaction, error) {
	if amount.LessThanOrEqual(decimal.Zero) { return nil, errors.New("amount must be > 0") }
	var w models.Wallet; if err := s.db.Where("id=? AND user_id=?", walletID, userID).First(&w).Error; err != nil { return nil, gorm.ErrRecordNotFound }
	var trx *models.Transaction; err := s.db.Transaction(func(tx *gorm.DB) error {
		if _, err := s.adjustBalance(tx, walletID, currency, amount.Neg()); err != nil { return err }
		trx = &models.Transaction{ WalletID: w.ID, Type: models.TrxWithdraw, Status: models.TrxSuccess, Amount: amount, CurrencyCode: currency, ExchangeRate: decimal.NewFromInt(1), ConvertedAmount: amount, Reference: fmt.Sprintf("WDR-%d", time.Now().UnixNano()) }
		return tx.Create(trx).Error
	}); return trx, err
}

func (s *WalletService) Transfer(userID, fromWalletID uint, toWalletID uint, fromCur, toCur string, amount decimal.Decimal) (*models.Transaction, error) {
	if amount.LessThanOrEqual(decimal.Zero) { return nil, errors.New("amount must be > 0") }
	var fromW, toW models.Wallet
	if err := s.db.Where("id=? AND user_id=?", fromWalletID, userID).First(&fromW).Error; err != nil { return nil, gorm.ErrRecordNotFound }
	if err := s.db.Where("id=? AND user_id=?", toWalletID, userID).First(&toW).Error; err != nil { return nil, gorm.ErrRecordNotFound }
	var converted, rate decimal.Decimal; var trx *models.Transaction
	err := s.db.Transaction(func(tx *gorm.DB) error {
		c, r, err := s.cc.Convert(amount, fromCur, toCur); if err != nil { return err }
		converted, rate = c, r
		if _, err := s.adjustBalance(tx, fromWalletID, fromCur, amount.Neg()); err != nil { return err }
		if _, err := s.adjustBalance(tx, toWalletID, toCur, converted); err != nil { return err }
		trx = &models.Transaction{ WalletID: fromW.ID, ToWalletID: &toW.ID, Type: models.TrxTransfer, Status: models.TrxSuccess, Amount: amount, CurrencyCode: fromCur, ExchangeRate: rate, ConvertedAmount: converted, Reference: fmt.Sprintf("TRF-%d", time.Now().UnixNano()) }
		return tx.Create(trx).Error
	}); return trx, err
}

func (s *WalletService) Payment(userID, walletID uint, currency string, price decimal.Decimal, reference, metadata string) (*models.Transaction, error) {
	if price.LessThanOrEqual(decimal.Zero) { return nil, errors.New("amount must be > 0") }
	var w models.Wallet; if err := s.db.Where("id=? AND user_id=?", walletID, userID).First(&w).Error; err != nil { return nil, gorm.ErrRecordNotFound }
	var trx *models.Transaction; err := s.db.Transaction(func(tx *gorm.DB) error {
		if _, err := s.adjustBalance(tx, walletID, currency, price.Neg()); err != nil { return err }
		trx = &models.Transaction{ WalletID: w.ID, Type: models.TrxPayment, Status: models.TrxSuccess, Amount: price, CurrencyCode: currency, ExchangeRate: decimal.NewFromInt(1), ConvertedAmount: price, Reference: reference, Metadata: metadata }
		return tx.Create(trx).Error
	}); return trx, err
}

func (s *WalletService) Summary(userID uint, display string) (map[string]interface{}, error) {
	var wallets []models.Wallet; if err := s.db.Where("user_id=?", userID).Find(&wallets).Error; err != nil { return nil, err }
	var balances []models.WalletBalance; if len(wallets)==0 { return map[string]interface{}{"display_currency": display, "total": decimal.Zero, "by_currency": map[string]decimal.Decimal{}}, nil }
	var idsList []uint; for _, w := range wallets { idsList = append(idsList, w.ID) }
	if err := s.db.Where("wallet_id IN ?", idsList).Find(&balances).Error; err != nil { return nil, err }
	total := decimal.Zero; perCur := map[string]decimal.Decimal{}
	for _, b := range balances {
		if v, ok := perCur[b.CurrencyCode]; ok { perCur[b.CurrencyCode] = v.Add(b.Balance) } else { perCur[b.CurrencyCode] = b.Balance }
		conv, _, err := s.cc.Convert(b.Balance, b.CurrencyCode, display); if err != nil { return nil, err }
		total = total.Add(conv)
	}
	return map[string]interface{}{"display_currency": display, "total": total, "by_currency": perCur}, nil
}
