package main

import (
	"log"

	_ "kholid/ewallet/v2/docs"
	"kholid/ewallet/v2/internal/config"
	"kholid/ewallet/v2/internal/db"
	"kholid/ewallet/v2/internal/handler"
	"kholid/ewallet/v2/internal/middleware"
	"kholid/ewallet/v2/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @title E-Wallet API
// @version 1.0
// @description Multi-currency e-wallet API (Go + Fiber + GORM + MySQL). JWT Auth, Decimal-safe, Live FX, Pagination.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.Load()
	gdb := db.MustConnect(cfg)

	var currency service.CurrencyService
	// if os.Getenv("FX_PROVIDER") == "static" {
		currency = service.NewCurrencyServiceFromEnv(cfg.DisplayCurrency)
	// } else {
	// 	currency = service.NewLiveCurrencyService(cfg.DisplayCurrency)
	// }

	authSvc := service.NewAuthService(gdb, cfg.JWTSecret)
	walletSvc := service.NewWalletService(gdb, currency)
	trxSvc := service.NewTransactionService(gdb, currency)

	app := fiber.New()
	app.Get("/swagger/*", swagger.HandlerDefault) // http://localhost:8080/swagger/index.html

	h := handler.NewHandlers(authSvc, walletSvc, trxSvc, currency)

	v1 := app.Group("/api/v1")
	v1.Post("/auth/register", h.Register)
	v1.Post("/auth/login", h.Login)

	auth := v1.Use(middleware.Auth(cfg.JWTSecret))
	auth.Get("/me", h.Me)
	auth.Post("/wallets", h.CreateWallet)
	auth.Get("/wallets", h.ListWallets)
	auth.Get("/wallets/:id/balances", h.WalletBalances)

	auth.Post("/wallets/:id/deposit", h.Deposit)
	auth.Post("/wallets/:id/withdraw", h.Withdraw)
	auth.Post("/wallets/:id/transfer", h.Transfer)
	auth.Post("/wallets/:id/payment", h.Payment)

	auth.Get("/transactions", h.ListTransactions)
	auth.Get("/summary", h.Summary)

	log.Printf("listening on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil { log.Fatal(err) }
}
