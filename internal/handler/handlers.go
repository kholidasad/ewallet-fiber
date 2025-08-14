package handler

import (
	"strconv"
	"kholid/ewallet/v2/internal/models"
	"kholid/ewallet/v2/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
)

type Handlers struct { auth *service.AuthService; wallet *service.WalletService; trx *service.TransactionService; cc service.CurrencyService }
func NewHandlers(a *service.AuthService, w *service.WalletService, t *service.TransactionService, cc service.CurrencyService) *Handlers { return &Handlers{auth:a, wallet:w, trx:t, cc:cc} }

func (h *Handlers) Register(c *fiber.Ctx) error {
	t := struct{Email, Password string}{}
	if err := c.BodyParser(&t); err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()}) }
	u, err := h.auth.Register(t.Email, t.Password); if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()}) }
	return c.Status(fiber.StatusCreated).JSON(u)
}
func (h *Handlers) Login(c *fiber.Ctx) error {
	t := struct{Email, Password string}{}
	if err := c.BodyParser(&t); err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()}) }
	s, u, err := h.auth.Login(t.Email, t.Password); if err != nil { return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(fiber.Map{"token": s, "user": u})
}
func (h *Handlers) Me(c *fiber.Ctx) error { userID := c.Locals("userID").(uint); w, _ := h.wallet.ListWallets(userID); return c.JSON(fiber.Map{"id": userID, "wallets": w}) }

func (h *Handlers) CreateWallet(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint); var body struct{ Name string `json:"name"` }
	if err := c.BodyParser(&body); err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()}) }
	w, err := h.wallet.CreateWallet(userID, body.Name); if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()}) }
	return c.Status(fiber.StatusCreated).JSON(w)
}
func (h *Handlers) ListWallets(c *fiber.Ctx) error { userID := c.Locals("userID").(uint); w, err := h.wallet.ListWallets(userID); if err != nil { return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()}) }; return c.JSON(w) }
func (h *Handlers) WalletBalances(c *fiber.Ctx) error { userID := c.Locals("userID").(uint); id, _ := strconv.ParseUint(c.Params("id"), 10, 64); b, err := h.wallet.WalletBalances(userID, uint(id)); if err != nil { return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()}) }; return c.JSON(b) }

func (h *Handlers) Deposit(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint); id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
	var body struct{ Currency string; Amount string }; if err := c.BodyParser(&body); err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()}) }
	amt, err := decimal.NewFromString(body.Amount); if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"invalid amount"}) }
	trx, err := h.wallet.Deposit(userID, uint(id), body.Currency, amt); if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(trx)
}
func (h *Handlers) Withdraw(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint); id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
	var body struct{ Currency string; Amount string }; if err := c.BodyParser(&body); err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()}) }
	amt, err := decimal.NewFromString(body.Amount); if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"invalid amount"}) }
	trx, err := h.wallet.Withdraw(userID, uint(id), body.Currency, amt); if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(trx)
}
func (h *Handlers) Transfer(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint); id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
	var body struct{ ToWalletID uint `json:"to_wallet_id"`; FromCurrency, ToCurrency, Amount string }
	if err := c.BodyParser(&body); err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()}) }
	amt, err := decimal.NewFromString(body.Amount); if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"invalid amount"}) }
	trx, err := h.wallet.Transfer(userID, uint(id), body.ToWalletID, body.FromCurrency, body.ToCurrency, amt); if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(trx)
}
func (h *Handlers) Payment(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint); id, _ := strconv.ParseUint(c.Params("id"), 10, 64)
	var body struct{ Currency, Amount, Reference, Metadata string }; if err := c.BodyParser(&body); err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()}) }
	amt, err := decimal.NewFromString(body.Amount); if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"invalid amount"}) }
	trx, err := h.wallet.Payment(userID, uint(id), body.Currency, amt, body.Reference, body.Metadata); if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()}) }
	return c.JSON(trx)
}
func (h *Handlers) ListTransactions(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint); page, _ := strconv.Atoi(c.Query("page", "1")); size, _ := strconv.Atoi(c.Query("page_size", "20"))
	var q service.TrxQuery; if t := c.Query("type"); t != "" { tt := models.TransactionType(t); q.Type = &tt }; if cur := c.Query("currency"); cur != "" { q.Currency = &cur }
	res, err := h.trx.List(userID, page, size, q); if err != nil { return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()}) }
	pages := (res.Total + int64(size) - 1) / int64(size); return c.JSON(fiber.Map{"data": res.Data, "meta": fiber.Map{"page": page, "page_size": size, "total": res.Total, "total_pages": pages}})
}
func (h *Handlers) Summary(c *fiber.Ctx) error { userID := c.Locals("userID").(uint); res, err := h.wallet.Summary(userID, h.cc.DisplayCurrency()); if err != nil { return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()}) }; return c.JSON(res) }
