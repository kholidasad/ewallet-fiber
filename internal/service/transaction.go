package service

import ( "kholid/ewallet/v2/internal/models"; "gorm.io/gorm" )

type TransactionService struct { db *gorm.DB; cc CurrencyService }
func NewTransactionService(db *gorm.DB, cc CurrencyService) *TransactionService { return &TransactionService{db: db, cc: cc} }

type TrxQuery struct { Type *models.TransactionType; Currency *string }
type PagedTrx struct { Data []models.Transaction; Total int64 }

func (s *TransactionService) List(userID uint, page, pageSize int, q TrxQuery) (PagedTrx, error) {
	var ids []uint; s.db.Model(&models.Wallet{}).Where("user_id = ?", userID).Pluck("id", &ids)
	dbq := s.db.Model(&models.Transaction{}).Where("wallet_id IN ?", ids)
	if q.Type != nil { dbq = dbq.Where("type = ?", *q.Type) }
	if q.Currency != nil { dbq = dbq.Where("currency_code = ?", *q.Currency) }
	var total int64; if err := dbq.Count(&total).Error; err != nil { return PagedTrx{}, err }
	if page<1 { page=1 }; if pageSize<=0 || pageSize>100 { pageSize=20 }; offset := (page-1)*pageSize
	var trxs []models.Transaction
	if err := dbq.Order("created_at DESC, id DESC").Limit(pageSize).Offset(offset).Find(&trxs).Error; err != nil { return PagedTrx{}, err }
	return PagedTrx{Data: trxs, Total: total}, nil
}
