package service

import (
	"encoding/json"; "errors"; "fmt"; "net/http"; "os"; "strconv"; "strings"; "sync"; "time"
	"github.com/shopspring/decimal"
)

type LiveCurrencyService struct {
	client *http.Client; base string; allowed []string; cache map[string]decimal.Decimal
	mu sync.RWMutex; expiry time.Time; ttl time.Duration; display string; providerURL string; apiKey string
}
func NewLiveCurrencyService(display string) *LiveCurrencyService {
	allowed := strings.Split(getenv("FX_CURRENCIES", "USD,EUR,JPY"), ",")
	return &LiveCurrencyService{ client: &http.Client{Timeout: 5*time.Second}, base: getenv("FX_BASE","USD"), allowed: allowed,
		cache: map[string]decimal.Decimal{}, ttl: time.Duration(mustAtoi(getenv("FX_TTL_SECONDS","300")))*time.Second,
		display: display, providerURL: getenv("FX_PROVIDER_URL","https://api.exchangerate.host/latest"), apiKey: os.Getenv("FX_API_KEY") }
}
func (l *LiveCurrencyService) List() []string { return l.allowed }
func (l *LiveCurrencyService) DisplayCurrency() string { return l.display }
func (l *LiveCurrencyService) ensureRates() error {
	l.mu.RLock(); if time.Now().Before(l.expiry) && len(l.cache)>0 { l.mu.RUnlock(); return nil }; l.mu.RUnlock()
	l.mu.Lock(); defer l.mu.Unlock(); if time.Now().Before(l.expiry) && len(l.cache)>0 { return nil }
	symbols := strings.Join(l.allowed, ","); url := fmt.Sprintf("%s?base=%s&symbols=%s", l.providerURL, l.base, symbols)
	req, _ := http.NewRequest("GET", url, nil); if l.apiKey!="" { req.Header.Set("apikey", l.apiKey) }
	resp, err := l.client.Do(req); if err != nil { return err }; defer resp.Body.Close()
	if resp.StatusCode != 200 { return fmt.Errorf("fx provider status %d", resp.StatusCode) }
	var out struct{ Rates map[string]float64 `json:"rates"` }
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil { return err }
	if len(out.Rates)==0 { return errors.New("empty rates") }
	l.cache = map[string]decimal.Decimal{}; for k,v := range out.Rates { l.cache[k] = decimal.NewFromFloat(v) }; l.cache[l.base] = decimal.NewFromInt(1); l.expiry = time.Now().Add(l.ttl); return nil
}
func (l *LiveCurrencyService) Rate(code string) (decimal.Decimal, error) { if err := l.ensureRates(); err != nil { return decimal.Zero, err }; l.mu.RLock(); defer l.mu.RUnlock(); r, ok := l.cache[code]; if !ok { return decimal.Zero, errors.New("unsupported currency") }; return r, nil }
func (l *LiveCurrencyService) Convert(amount decimal.Decimal, from, to string) (decimal.Decimal, decimal.Decimal, error) {
	if from==to { return amount, decimal.NewFromInt(1), nil }
	if err := l.ensureRates(); err != nil { return decimal.Zero, decimal.Zero, err }
	l.mu.RLock(); defer l.mu.RUnlock(); fromRate, ok := l.cache[from]; if !ok { return decimal.Zero, decimal.Zero, errors.New("unsupported currency") }
	toRate, ok := l.cache[to]; if !ok { return decimal.Zero, decimal.Zero, errors.New("unsupported currency") }
	usd := amount.Div(fromRate); converted := usd.Mul(toRate); rate := toRate.Div(fromRate); return converted, rate, nil
}
func mustAtoi(s string) int { n, _ := strconv.Atoi(s); if n<=0 { return 300 }; return n }
