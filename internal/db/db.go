package db

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"kholid/ewallet/v2/internal/config"
	"kholid/ewallet/v2/internal/models"
)

func MustConnect(cfg config.Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	gdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{ Logger: logger.Default.LogMode(logger.Warn) })
	if err != nil { log.Fatalf("db connect: %v", err) }
	if err := migrate(gdb); err != nil { log.Fatalf("migrate: %v", err) }
	return gdb
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{}, &models.Wallet{}, &models.Currency{},
		&models.WalletBalance{}, &models.Transaction{},
	)
}
