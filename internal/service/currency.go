package service

import ( "errors"; "os"; "github.com/shopspring/decimal" )

type CurrencyService interface {
	List() []string
	Rate(code string) (decimal.Decimal, error)
	Convert(amount decimal.Decimal, from, to string) (decimal.Decimal, decimal.Decimal, error)
	DisplayCurrency() string
}

type StaticCurrencyService struct { rates map[string]decimal.Decimal; display string }
func NewCurrencyServiceFromEnv(display string) *StaticCurrencyService {
	return &StaticCurrencyService{ rates: map[string]decimal.Decimal{
		"USD": decimal.RequireFromString(getenv("RATE_USD", "1")),
		"EUR": decimal.RequireFromString(getenv("RATE_EUR", "0.9")),
		"JPY": decimal.RequireFromString(getenv("RATE_JPY", "155")),
	}, display: display }
}
func (s *StaticCurrencyService) List() []string { return []string{"USD","EUR","JPY"} }
func (s *StaticCurrencyService) Rate(code string) (decimal.Decimal, error) { r, ok := s.rates[code]; if !ok { return decimal.Zero, errors.New("unsupported currency") }; return r, nil }
func (s *StaticCurrencyService) Convert(amount decimal.Decimal, from, to string) (decimal.Decimal, decimal.Decimal, error) {
	if from==to { return amount, decimal.NewFromInt(1), nil }
	fromRate, ok := s.rates[from]; if !ok { return decimal.Zero, decimal.Zero, errors.New("unsupported currency") }
	toRate, ok := s.rates[to]; if !ok { return decimal.Zero, decimal.Zero, errors.New("unsupported currency") }
	usd := amount.Div(fromRate); converted := usd.Mul(toRate); rate := toRate.Div(fromRate); return converted, rate, nil
}
func (s *StaticCurrencyService) DisplayCurrency() string { return s.display }
func getenv(k, def string) string { if v := os.Getenv(k); v != "" { return v }; return def }
